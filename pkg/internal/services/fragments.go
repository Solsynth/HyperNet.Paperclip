package services

import (
	"errors"
	"fmt"
	"math"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"git.solsynth.dev/hypernet/nexus/pkg/nex/cachekit"
	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	"git.solsynth.dev/hypernet/paperclip/pkg/filekit/models"
	"git.solsynth.dev/hypernet/paperclip/pkg/internal/database"
	"git.solsynth.dev/hypernet/paperclip/pkg/internal/gap"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func KgAttachmentFragmentCache(rid string) string {
	return cachekit.FKey("attachment-fragment", rid)
}

func NewAttachmentFragment(tx *gorm.DB, user *sec.UserInfo, fragment models.AttachmentFragment) (models.AttachmentFragment, error) {
	if fragment.Fingerprint != nil {
		var existsFragment models.AttachmentFragment
		if err := database.C.Where(models.AttachmentFragment{
			Fingerprint: fragment.Fingerprint,
			AccountID:   user.ID,
		}).First(&existsFragment).Error; err == nil {
			return existsFragment, nil
		}
	}

	fragment.Uuid = uuid.NewString()
	fragment.Rid = RandString(16)
	fragment.FileChunks = datatypes.JSONMap{}
	fragment.AccountID = user.ID

	chunkSize := viper.GetInt64("performance.file_chunk_size")
	chunkCount := math.Ceil(float64(fragment.Size) / float64(chunkSize))
	for idx := 0; idx < int(chunkCount); idx++ {
		cid := RandString(8)
		fragment.FileChunks[cid] = idx
	}

	// If the user didn't provide file mimetype manually, we have to detect it
	if len(fragment.MimeType) == 0 {
		if ext := filepath.Ext(fragment.Name); len(ext) > 0 {
			// Detect mimetype by file extensions
			fragment.MimeType = mime.TypeByExtension(ext)
		}
	}

	if err := tx.Save(&fragment).Error; err != nil {
		return fragment, fmt.Errorf("failed to save attachment record: %v", err)
	}

	return fragment, nil
}

func GetFragmentByRID(rid string) (models.AttachmentFragment, error) {
	if val, err := cachekit.Get[models.AttachmentFragment](
		gap.Ca,
		KgAttachmentFragmentCache(rid),
	); err == nil {
		return val, nil
	}

	var attachment models.AttachmentFragment
	if err := database.C.Where(models.AttachmentFragment{
		Rid: rid,
	}).Preload("Pool").First(&attachment).Error; err != nil {
		return attachment, err
	} else {
		CacheAttachmentFragment(attachment)
	}

	return attachment, nil
}

func CacheAttachmentFragment(item models.AttachmentFragment) {
	cachekit.Set[models.AttachmentFragment](
		gap.Ca,
		KgAttachmentFragmentCache(item.Rid),
		item,
		60*time.Minute,
	)
}

func UploadFragmentChunk(ctx *fiber.Ctx, cid string, file *multipart.FileHeader, meta models.AttachmentFragment) error {
	destMap := viper.GetStringMap("destinations.0")

	var dest models.LocalDestination
	rawDest, _ := jsoniter.Marshal(destMap)
	_ = jsoniter.Unmarshal(rawDest, &dest)

	tempPath := filepath.Join(dest.Path, fmt.Sprintf("%s.part%s.partial", meta.Uuid, cid))
	destPath := filepath.Join(dest.Path, fmt.Sprintf("%s.part%s", meta.Uuid, cid))
	if err := ctx.SaveFile(file, tempPath); err != nil {
		return err
	}
	return os.Rename(tempPath, destPath)
}

func UploadFragmentChunkBytes(ctx *fiber.Ctx, cid string, raw []byte, meta models.AttachmentFragment) error {
	destMap := viper.GetStringMap("destinations.0")

	var dest models.LocalDestination
	rawDest, _ := jsoniter.Marshal(destMap)
	_ = jsoniter.Unmarshal(rawDest, &dest)

	tempPath := filepath.Join(dest.Path, fmt.Sprintf("%s.part%s.partial", meta.Uuid, cid))
	destPath := filepath.Join(dest.Path, fmt.Sprintf("%s.part%s", meta.Uuid, cid))
	if err := os.WriteFile(tempPath, raw, 0644); err != nil {
		return err
	}
	return os.Rename(tempPath, destPath)
}

func CheckFragmentChunkExists(meta models.AttachmentFragment, cid string) bool {
	destMap := viper.GetStringMap("destinations.0")

	var dest models.LocalDestination
	rawDest, _ := jsoniter.Marshal(destMap)
	_ = jsoniter.Unmarshal(rawDest, &dest)

	path := filepath.Join(dest.Path, fmt.Sprintf("%s.part%s", meta.Uuid, cid))
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}

func FindFragmentMissingChunks(meta models.AttachmentFragment) []string {
	var missing []string
	for cid := range meta.FileChunks {
		if !CheckFragmentChunkExists(meta, cid) {
			missing = append(missing, cid)
		}
	}
	return missing
}
