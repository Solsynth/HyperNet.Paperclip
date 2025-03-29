package services

import (
	"time"

	"git.solsynth.dev/hypernet/paperclip/pkg/internal/database"
	"git.solsynth.dev/hypernet/paperclip/pkg/filekit/models"

	"github.com/rs/zerolog/log"
)

func DoUnusedAttachmentCleanup() {
	deadline := time.Now().Add(-60 * time.Minute)

	var result []models.Attachment
	if err := database.C.Where("created_at < ? AND used_count = 0", deadline).
		Find(&result).Error; err != nil {
		log.Error().Err(err).Msg("An error occurred when getting unused attachments...")
		return
	}

	if err := DeleteAttachmentInBatch(result); err != nil {
		log.Error().Err(err).Msg("An error occurred when deleting unused attachments...")
		return
	}

	log.Info().Int("count", len(result)).Msg("Deleted unused attachments...")
}
