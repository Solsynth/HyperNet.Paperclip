package api

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/viper"

	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	"git.solsynth.dev/hypernet/paperclip/pkg/internal/database"
	"git.solsynth.dev/hypernet/paperclip/pkg/internal/server/exts"

	"git.solsynth.dev/hypernet/paperclip/pkg/filekit/models"
	"git.solsynth.dev/hypernet/paperclip/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
)

func getBillingStatus(c *fiber.Ctx) error {
	if err := sec.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("nex_user").(*sec.UserInfo)

	currentBytes, err := services.GetLastDayUploadedBytes(user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	discountFileSize := viper.GetInt64("payment.discount")

	return c.JSON(fiber.Map{
		"current_bytes":      currentBytes,
		"discount_file_size": discountFileSize,
		"included_ratio":     float64(currentBytes) / float64(discountFileSize),
	})
}

func openAttachment(c *fiber.Ctx) error {
	id := c.Params("id")
	region := c.Query("region")

	var err error
	var url, mimetype string
	var filesize int64
	size := lo.Ternary(c.QueryBool("preview", true), 1024, -1)
	if len(region) > 0 {
		url, filesize, mimetype, err = services.OpenAttachmentByRID(id, size, region)
	} else {
		url, filesize, mimetype, err = services.OpenAttachmentByRID(id, size)
	}

	authenticated := false
	if err := sec.EnsureAuthenticated(c); err == nil {
		authenticated = true
	}

	if filesize > viper.GetInt64("traffic.maximum_size") && !authenticated {
		return fiber.NewError(fiber.StatusForbidden, "file is too large, you need authorized to access")
	}

	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	c.Set(fiber.HeaderContentType, mimetype)

	if strings.HasPrefix(url, "file://") {
		fp := strings.Replace(url, "file://", "", 1)
		return c.SendFile(fp)
	}

	return c.Redirect(url, fiber.StatusFound)
}

func getAttachmentMeta(c *fiber.Ctx) error {
	id := c.Params("id")

	metadata, err := services.GetAttachmentByRID(id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound)
	}

	return c.JSON(metadata)
}

func updateAttachmentMeta(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id", 0)
	user := c.Locals("nex_user").(*sec.UserInfo)

	var data struct {
		Thumbnail   *uint          `json:"thumbnail"`
		Compressed  *uint          `json:"compressed"`
		Alternative string         `json:"alt"`
		Metadata    map[string]any `json:"metadata"`
		IsIndexable bool           `json:"is_indexable"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	var attachment models.Attachment
	if err := database.C.
		Where("id = ? AND account_id = ?", id, user.ID).
		Preload("Thumbnail").
		Preload("Compressed").
		First(&attachment).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if data.Thumbnail != nil && attachment.ThumbnailID != data.Thumbnail {
		var thumbnail models.Attachment
		if err := database.C.
			Where("id = ? AND account_id = ?", data.Thumbnail, user.ID).
			First(&thumbnail).Error; err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable find thumbnail: %v", err))
		}
		if attachment.Thumbnail != nil {
			services.UnsetAttachmentAsThumbnail(*attachment.Thumbnail)
		}
		thumbnail, err := services.SetAttachmentAsThumbnail(thumbnail)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable set thumbnail: %v", err))
		}
		attachment.Thumbnail = &thumbnail
		attachment.ThumbnailID = &thumbnail.ID
	}
	if data.Compressed != nil && attachment.CompressedID != data.Compressed {
		var compressed models.Attachment
		if err := database.C.
			Where("id = ? AND account_id = ?", data.Compressed, user.ID).
			First(&compressed).Error; err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable find compressed: %v", err))
		}
		if attachment.Compressed != nil {
			services.UnsetAttachmentAsCompressed(*attachment.Compressed)
		}
		compressed, err := services.SetAttachmentAsCompressed(compressed)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable set compressed: %v", err))
		}
		attachment.Compressed = &compressed
		attachment.CompressedID = &compressed.ID
	}

	attachment.Alternative = data.Alternative
	attachment.Usermeta = data.Metadata
	attachment.IsIndexable = data.IsIndexable

	if attachment, err := services.UpdateAttachment(attachment); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(attachment)
	}
}

func updateAttachmentRating(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id", 0)
	user := c.Locals("nex_user").(*sec.UserInfo)

	var data struct {
		ContentRating int `json:"content_rating" validate:"required,min=3,max=21"`
		QualityRating int `json:"quality_rating" validate:"min=0,max=5"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	attachment, err := services.GetAttachmentByID(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	} else if attachment.AccountID != user.ID {
		if err = sec.EnsureGrantedPerm(c, "OverrideAttachmentRating", true); err != nil {
			return err
		}
	}

	attachment.ContentRating = data.ContentRating
	attachment.QualityRating = data.QualityRating

	if attachment, err = services.UpdateAttachment(attachment); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(attachment)
	}
}

func deleteAttachment(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id", 0)
	user := c.Locals("nex_user").(*sec.UserInfo)

	attachment, err := services.GetAttachmentByID(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	} else if attachment.AccountID != user.ID {
		return fiber.NewError(fiber.StatusNotFound, "record not created by you")
	}

	if err := services.DeleteAttachment(attachment); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		return c.SendStatus(fiber.StatusOK)
	}
}
