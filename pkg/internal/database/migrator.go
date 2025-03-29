package database

import (
	"git.solsynth.dev/hypernet/paperclip/pkg/filekit/models"
	"gorm.io/gorm"
)

var AutoMaintainRange = []any{
	&models.AttachmentPool{},
	&models.Attachment{},
	&models.AttachmentFragment{},
	&models.AttachmentBoost{},
	&models.StickerPack{},
	&models.Sticker{},
	&models.StickerPackOwnership{},
}

func RunMigration(source *gorm.DB) error {
	if err := source.AutoMigrate(
		AutoMaintainRange...,
	); err != nil {
		return err
	}

	return nil
}
