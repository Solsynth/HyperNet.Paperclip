package api

import (
	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	"git.solsynth.dev/hypernet/paperclip/pkg/filekit/models"
	"git.solsynth.dev/hypernet/paperclip/pkg/internal/database"
	"git.solsynth.dev/hypernet/paperclip/pkg/internal/gap"
	"git.solsynth.dev/hypernet/paperclip/pkg/internal/server/exts"
	"git.solsynth.dev/hypernet/paperclip/pkg/internal/services"
	"git.solsynth.dev/hypernet/passport/pkg/authkit"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func listStickerPacks(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)

	if take > 100 {
		take = 100
	}

	tx := database.C

	if len(c.Query("author")) > 0 {
		author, err := authkit.GetUserByName(gap.Nx, c.Query("author"))
		if err == nil {
			tx = tx.Where("account_id = ?", author.ID)
		}
	}

	var count int64
	if err := database.C.Model(&models.StickerPack{}).Count(&count).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var packs []models.StickerPack
	if err := tx.Limit(take).Offset(offset).Find(&packs).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  packs,
	})
}

func listOwnedStickerPacks(c *fiber.Ctx) error {
	if err := sec.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("nex_user").(*sec.UserInfo)

	var ownerships []models.StickerPackOwnership
	if err := database.C.Where("account_id = ?", user.ID).Find(&ownerships).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	idSet := lo.Map(ownerships, func(o models.StickerPackOwnership, _ int) uint {
		return o.PackID
	})

	var packs []models.StickerPack
	if err := database.C.Where("id IN ?", idSet).Find(&packs).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(packs)
}

func getStickerPack(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("packId", 0)
	pack, err := services.GetStickerPack(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	var stickers []models.Sticker
	if err := database.C.Where("pack_id = ?", pack.ID).
		Preload("Attachment").
		Find(&stickers).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		pack.Stickers = stickers
	}

	return c.JSON(pack)
}

func createStickerPack(c *fiber.Ctx) error {
	user := c.Locals("nex_user").(*sec.UserInfo)

	var data struct {
		Prefix      string `json:"prefix" validate:"required,alphanum,min=2,max=12"`
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	pack, err := services.NewStickerPack(user, data.Prefix, data.Name, data.Description)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(pack)
}

func updateStickerPack(c *fiber.Ctx) error {
	user := c.Locals("nex_user").(*sec.UserInfo)

	var data struct {
		Prefix      string `json:"prefix" validate:"required,alphanum,min=2,max=12"`
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	id, _ := c.ParamsInt("packId", 0)
	pack, err := services.GetStickerPackWithUser(uint(id), user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	pack.Prefix = data.Prefix
	pack.Name = data.Name
	pack.Description = data.Description

	if pack, err = services.UpdateStickerPack(pack); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(pack)
}

func deleteStickerPack(c *fiber.Ctx) error {
	user := c.Locals("nex_user").(*sec.UserInfo)

	id, _ := c.ParamsInt("packId", 0)
	pack, err := services.GetStickerPackWithUser(uint(id), user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if pack, err = services.DeleteStickerPack(pack); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(pack)
}

func addStickerPack(c *fiber.Ctx) error {
	if err := sec.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("nex_user").(*sec.UserInfo)

	packId, _ := c.ParamsInt("packId", 0)
	var pack models.StickerPack
	if err := database.C.Where("id = ?", packId).First(&pack).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	ownership, err := services.AddStickerPack(user.ID, pack)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(ownership)
}

func removeStickerPack(c *fiber.Ctx) error {
	if err := sec.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("nex_user").(*sec.UserInfo)

	packId, _ := c.ParamsInt("packId", 0)
	var pack models.StickerPack
	if err := database.C.Where("id = ?", packId).First(&pack).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	ownership, err := services.RemoveStickerPack(user.ID, pack)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(ownership)
}
