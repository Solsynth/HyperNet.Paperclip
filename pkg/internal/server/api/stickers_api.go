package api

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	"git.solsynth.dev/hypernet/paperclip/pkg/filekit/models"
	"git.solsynth.dev/hypernet/paperclip/pkg/internal/database"
	"git.solsynth.dev/hypernet/paperclip/pkg/internal/server/exts"
	"git.solsynth.dev/hypernet/paperclip/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
)

func lookupStickerBatch(c *fiber.Ctx) error {
	probe := c.Query("probe")
	if stickers, err := services.GetStickerLikeAlias(probe); err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	} else {
		return c.JSON(stickers)
	}
}

func getStickerByAlias(c *fiber.Ctx) error {
	alias := c.Params("alias")
	if sticker, err := services.GetStickerWithAlias(alias); err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	} else {
		return c.JSON(sticker)
	}
}

func openStickerByAlias(c *fiber.Ctx) error {
	alias := c.Params("alias")
	region := c.Query("region")

	sticker, err := services.GetStickerWithAlias(alias)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	var url, mimetype string
	if len(region) > 0 {
		url, mimetype, err = services.OpenAttachmentByRID(sticker.Attachment.Rid, true, region)
	} else {
		url, mimetype, err = services.OpenAttachmentByRID(sticker.Attachment.Rid, true)
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

func listStickers(c *fiber.Ctx) error {
	if err := sec.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("nex_user").(*sec.UserInfo)

	var ownerships []models.StickerPackOwnership
	if err := database.C.Where("account_id = ?", user.ID).Find(&ownerships).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	tx := database.C.Where("pack_id IN ?", lo.Map(ownerships, func(o models.StickerPackOwnership, _ int) uint {
		return o.PackID
	}))

	var stickers []models.Sticker
	if err := tx.
		Preload("Attachment").Preload("Pack").
		Find(&stickers).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(stickers)
}

func getSticker(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("stickerId", 0)
	sticker, err := services.GetSticker(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	return c.JSON(sticker)
}

func createSticker(c *fiber.Ctx) error {
	user := c.Locals("nex_user").(*sec.UserInfo)

	var data struct {
		Alias        string `json:"alias" validate:"required,alphanum,min=2,max=12"`
		Name         string `json:"name" validate:"required"`
		AttachmentID string `json:"attachment_id"`
		PackID       uint   `json:"pack_id"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	var attachment models.Attachment
	if err := database.C.Where("rid = ?", data.AttachmentID).First(&attachment).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to find attachment: %v", err))
	} else if !attachment.IsAnalyzed {
		return fiber.NewError(fiber.StatusBadRequest, "sticker attachment must be analyzed")
	}

	if strings.SplitN(attachment.MimeType, "/", 2)[0] != "image" {
		return fiber.NewError(fiber.StatusBadRequest, "sticker attachment must be an image")
	}

	var pack models.StickerPack
	if err := database.C.Where("id = ? AND account_id = ?", data.PackID, user.ID).First(&pack).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to find pack: %v", err))
	}

	sticker, err := services.NewSticker(models.Sticker{
		Alias:        data.Alias,
		Name:         data.Name,
		Attachment:   attachment,
		AccountID:    user.ID,
		PackID:       pack.ID,
		AttachmentID: attachment.ID,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(sticker)
}

func updateSticker(c *fiber.Ctx) error {
	user := c.Locals("nex_user").(*sec.UserInfo)

	var data struct {
		Alias        string `json:"alias" validate:"required,alphanum,min=2,max=12"`
		Name         string `json:"name" validate:"required"`
		AttachmentID string `json:"attachment_id"`
		PackID       uint   `json:"pack_id"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	id, _ := c.ParamsInt("stickerId", 0)
	sticker, err := services.GetStickerWithUser(uint(id), user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	var attachment models.Attachment
	if err := database.C.Where("rid = ?", data.AttachmentID).First(&attachment).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to find attachment: %v", err))
	} else if !attachment.IsAnalyzed {
		return fiber.NewError(fiber.StatusBadRequest, "sticker attachment must be analyzed")
	}

	if strings.SplitN(attachment.MimeType, "/", 2)[0] != "image" {
		return fiber.NewError(fiber.StatusBadRequest, "sticker attachment must be an image")
	}

	var pack models.StickerPack
	if err := database.C.Where("id = ? AND account_id = ?", data.PackID, user.ID).First(&pack).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to find pack: %v", err))
	}

	sticker.Alias = data.Alias
	sticker.Name = data.Name
	sticker.PackID = data.PackID
	sticker.AttachmentID = attachment.ID

	if sticker, err = services.UpdateSticker(sticker); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(sticker)
}

func deleteSticker(c *fiber.Ctx) error {
	user := c.Locals("nex_user").(*sec.UserInfo)

	id, _ := c.ParamsInt("stickerId", 0)
	sticker, err := services.GetStickerWithUser(uint(id), user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if sticker, err = services.DeleteSticker(sticker); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(sticker)
}
