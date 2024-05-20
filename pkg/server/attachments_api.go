package server

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"

	"git.solsynth.dev/hydrogen/paperclip/pkg/grpc"
	"git.solsynth.dev/hydrogen/passport/pkg/grpc/proto"

	"git.solsynth.dev/hydrogen/paperclip/pkg/database"
	"git.solsynth.dev/hydrogen/paperclip/pkg/models"
	"git.solsynth.dev/hydrogen/paperclip/pkg/services"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

func openAttachment(c *fiber.Ctx) error {
	id := c.Params("id")

	metadata, err := services.GetAttachmentByUUID(id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound)
	}

	destMap := viper.GetStringMap("destinations")
	dest, destOk := destMap[metadata.Destination]
	if !destOk {
		return fiber.NewError(fiber.StatusInternalServerError, "invalid destination: destination configuration was not found")
	}

	var destParsed models.BaseDestination
	rawDest, _ := jsoniter.Marshal(dest)
	_ = jsoniter.Unmarshal(rawDest, &destParsed)

	switch destParsed.Type {
	case models.DestinationTypeLocal:
		var destConfigured models.LocalDestination
		_ = jsoniter.Unmarshal(rawDest, &destConfigured)
		return c.SendFile(filepath.Join(destConfigured.Path, metadata.Uuid))
	case models.DestinationTypeS3:
		var destConfigured models.S3Destination
		_ = jsoniter.Unmarshal(rawDest, &destConfigured)
		protocol := lo.Ternary(destConfigured.EnableSSL, "https", "http")
		return c.Redirect(fmt.Sprintf(
			"%s://%s.%s/%s",
			protocol,
			destConfigured.Bucket,
			destConfigured.Endpoint,
			url.QueryEscape(filepath.Join(destConfigured.Path, metadata.Uuid)),
		))
	default:
		return fmt.Errorf("invalid destination: unsupported protocol %s", destParsed.Type)
	}
}

func getAttachmentMeta(c *fiber.Ctx) error {
	id := c.Params("id")

	metadata, err := services.GetAttachmentByUUID(id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound)
	}

	return c.JSON(metadata)
}

func createAttachment(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	destName := c.Query("destination", viper.GetString("preferred_destination"))

	hash := c.FormValue("hash")
	if len(hash) != 64 {
		return fiber.NewError(fiber.StatusBadRequest, "please provide a sha-256 hash code, length should be 64 characters")
	}
	usage := c.FormValue("usage")
	if !lo.Contains(viper.GetStringSlice("accepts_usage"), usage) {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("disallowed usage: %s", usage))
	}

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	requiredPerm, _ := jsoniter.Marshal(file.Size)
	if result, err := grpc.Auth.CheckPerm(context.Background(), &proto.CheckPermRequest{
		Token: c.Locals("token").(string),
		Key:   "CreatePaperclipAttachments",
		Value: requiredPerm,
	}); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to check permission: %v", err))
	} else if !result.GetIsValid() {
		return fiber.NewError(
			fiber.StatusForbidden,
			fmt.Sprintf("requires permission CreatePaperclipAttachments equals or greater than %d", file.Size),
		)
	}

	usermeta := make(map[string]any)
	_ = jsoniter.UnmarshalFromString(c.FormValue("metadata"), &usermeta)

	tx := database.C.Begin()
	metadata, linked, err := services.NewAttachmentMetadata(tx, user, file, models.Attachment{
		Usage:       usage,
		HashCode:    hash,
		Alternative: c.FormValue("alt"),
		MimeType:    c.FormValue("mimetype"),
		Metadata:    usermeta,
		IsMature:    len(c.FormValue("mature")) > 0,
		Destination: destName,
	})
	if err != nil {
		tx.Rollback()
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !linked {
		if err := services.UploadFile(destName, c, file, metadata); err != nil {
			tx.Rollback()
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	tx.Commit()

	return c.JSON(metadata)
}

func updateAttachmentMeta(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	var data struct {
		Alternative string         `json:"alt"`
		Usage       string         `json:"usage"`
		Metadata    map[string]any `json:"metadata"`
		IsMature    bool           `json:"is_mature"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	}

	var attachment models.Attachment
	if err := database.C.Where(models.Attachment{
		Uuid:      c.Params("id"),
		AccountID: user.ID,
	}).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	attachment.Alternative = data.Alternative
	attachment.Usage = data.Usage
	attachment.Metadata = data.Metadata
	attachment.IsMature = data.IsMature

	if err := database.C.Save(&attachment).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(attachment)
}

func deleteAttachment(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id", 0)
	user := c.Locals("principal").(models.Account)

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
