package services

import (
	"errors"
	"fmt"
	"gorm.io/gorm"

	"git.solsynth.dev/hypernet/paperclip/pkg/internal/database"
	"git.solsynth.dev/hypernet/paperclip/pkg/internal/models"
	"github.com/spf13/viper"
)

func GetStickerLikeAlias(alias string) ([]models.Sticker, error) {
	var stickers []models.Sticker
	prefix := viper.GetString("database.prefix")
	if err := database.C.
		Joins(fmt.Sprintf("LEFT JOIN %ssticker_packs pk ON pack_id = pk.id", prefix)).
		Where("UPPER(CONCAT(pk.prefix, alias)) LIKE UPPER(?)", "%"+alias+"%").
		Preload("Attachment").Preload("Pack").
		Limit(10).
		Find(&stickers).Error; err != nil {
		return stickers, err
	}
	return stickers, nil
}

func GetStickerWithAlias(alias string) (models.Sticker, error) {
	var sticker models.Sticker
	prefix := viper.GetString("database.prefix")
	if err := database.C.
		Joins(fmt.Sprintf("LEFT JOIN %ssticker_packs pk ON pack_id = pk.id", prefix)).
		Where("UPPER(CONCAT(pk.prefix, alias)) = UPPER(?)", alias).
		Preload("Attachment").Preload("Pack").
		First(&sticker).Error; err != nil {
		return sticker, err
	}
	return sticker, nil
}

func GetSticker(id uint) (models.Sticker, error) {
	var sticker models.Sticker
	if err := database.C.Where("id = ?", id).Preload("Attachment").First(&sticker).Error; err != nil {
		return sticker, err
	}
	return sticker, nil
}

func GetStickerWithUser(id, userId uint) (models.Sticker, error) {
	var sticker models.Sticker
	if err := database.C.Where("id = ? AND account_id = ?", id, userId).First(&sticker).Error; err != nil {
		return sticker, err
	}
	return sticker, nil
}

func NewSticker(sticker models.Sticker) (models.Sticker, error) {
	if err := database.C.Save(&sticker).Error; err != nil {
		return sticker, err
	}
	return sticker, nil
}

func UpdateSticker(sticker models.Sticker) (models.Sticker, error) {
	if err := database.C.Save(&sticker).Error; err != nil {
		return sticker, err
	}
	return sticker, nil
}

func DeleteSticker(sticker models.Sticker) (models.Sticker, error) {
	if err := database.C.Delete(&sticker).Error; err != nil {
		return sticker, err
	}
	return sticker, nil
}

func AddStickerPack(user uint, pack models.StickerPack) (models.StickerPackOwnership, error) {
	var ownership models.StickerPackOwnership
	if err := database.C.
		Where("account_id = ?", user).
		First(&ownership).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return ownership, fmt.Errorf("unable to get current ownership: %v", err)
	} else if err == nil {
		return ownership, fmt.Errorf("you already own this pack")
	}

	ownership = models.StickerPackOwnership{
		AccountID: user,
		PackID:    pack.ID,
	}

	err := database.C.Save(&ownership).Error

	return ownership, err
}

func RemoveStickerPack(user uint, pack models.StickerPack) (models.StickerPackOwnership, error) {
	var ownership models.StickerPackOwnership
	if err := database.C.
		Where("account_id = ? AND pack_id = ?", user, pack.ID).
		First(&ownership).Error; err != nil {
		return ownership, fmt.Errorf("unable to get current ownership: %v", err)
	}

	err := database.C.Delete(&ownership).Error

	return ownership, err
}
