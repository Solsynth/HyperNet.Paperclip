package api

import (
	"strings"

	"git.solsynth.dev/hydrogen/paperclip/pkg/internal/database"
	"git.solsynth.dev/hydrogen/paperclip/pkg/internal/models"
	"github.com/gofiber/fiber/v2"
)

func listAttachment(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)

	if take > 100 {
		take = 100
	}

	tx := database.C.Order("created_at DESC")

	if len(c.Query("id")) > 0 {
		tx = tx.Where("id IN ?", strings.Split(c.Query("id"), ","))
	}

	if len(c.Query("author")) > 0 {
		var author models.Account
		if err := database.C.Where("name = ?", c.Query("author")).First(&author).Error; err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		} else {
			tx = tx.Where("account_id = ?", author.ID)
		}
	}

	if usage := c.Query("usage"); len(usage) > 0 {
		tx = tx.Where("usage IN ?", strings.Split(usage, " "))
	}

	var count int64
	countTx := tx
	if err := countTx.Model(&models.Attachment{}).Count(&count).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	var attachments []models.Attachment
	if err := tx.Offset(offset).Limit(take).Preload("Account").Find(&attachments).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  attachments,
	})
}
