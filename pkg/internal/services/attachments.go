package services

import (
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"git.solsynth.dev/hypernet/nexus/pkg/nex/cachekit"
	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"

	"git.solsynth.dev/hypernet/paperclip/pkg/internal/database"
	"git.solsynth.dev/hypernet/paperclip/pkg/internal/fs"
	"git.solsynth.dev/hypernet/paperclip/pkg/internal/gap"

	"git.solsynth.dev/hypernet/paperclip/pkg/filekit/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func KgAttachmentCache(rid string) string {
	return cachekit.FKey(cachekit.DAAttachment, rid)
}

func GetAttachmentByID(id uint) (models.Attachment, error) {
	var attachment models.Attachment
	if err := database.C.
		Where("id = ?", id).
		Preload("Pool").
		Preload("Thumbnail").
		Preload("Compressed").
		Preload("Boosts").
		First(&attachment).Error; err != nil {
		return attachment, err
	} else {
		CacheAttachment(attachment)
	}

	return attachment, nil
}

func GetAttachmentByRID(rid string) (models.Attachment, error) {
	if val, err := cachekit.Get[models.Attachment](
		gap.Ca,
		KgAttachmentCache(rid),
	); err == nil {
		return val, nil
	}

	var attachment models.Attachment
	if err := database.C.Where(models.Attachment{
		Rid: rid,
	}).
		Preload("Pool").
		Preload("Thumbnail").
		Preload("Compressed").
		Preload("Boosts").
		First(&attachment).Error; err != nil {
		return attachment, err
	} else {
		CacheAttachment(attachment)
	}

	return attachment, nil
}

func GetAttachmentByHash(hash string) (models.Attachment, error) {
	var attachment models.Attachment
	if err := database.C.Where(models.Attachment{
		HashCode: hash,
	}).Preload("Pool").First(&attachment).Error; err != nil {
		return attachment, err
	}
	return attachment, nil
}

func GetAttachmentCache(rid string) (models.Attachment, bool) {
	if val, err := cachekit.Get[models.Attachment](
		gap.Ca,
		KgAttachmentCache(rid),
	); err == nil {
		return val, true
	}
	return models.Attachment{}, false
}

func CacheAttachment(item models.Attachment) {
	cachekit.Set[models.Attachment](
		gap.Ca,
		KgAttachmentCache(item.Rid),
		item,
		60*time.Minute,
	)
}

func NewAttachmentMetadata(tx *gorm.DB, user *sec.UserInfo, file *multipart.FileHeader, attachment models.Attachment) (models.Attachment, error) {
	attachment.Uuid = uuid.NewString()
	attachment.Rid = RandString(16)
	attachment.Size = file.Size
	attachment.Name = file.Filename
	attachment.AccountID = user.ID

	// If the user didn't provide file mimetype manually, we have to detect it
	if len(attachment.MimeType) == 0 {
		if ext := filepath.Ext(attachment.Name); len(ext) > 0 {
			// Detect mimetype by file extensions
			attachment.MimeType = mime.TypeByExtension(ext)
		} else {
			// Detect mimetype by file header
			// This method as a fallback method, because this isn't pretty accurate
			header, err := file.Open()
			if err != nil {
				return attachment, fmt.Errorf("failed to read file header: %v", err)
			}
			defer header.Close()

			fileHeader := make([]byte, 512)
			_, err = header.Read(fileHeader)
			if err != nil {
				return attachment, err
			}
			attachment.MimeType = http.DetectContentType(fileHeader)
		}
	}

	if err := tx.Save(&attachment).Error; err != nil {
		return attachment, fmt.Errorf("failed to save attachment record: %v", err)
	}

	return attachment, nil
}

func TryLinkAttachment(tx *gorm.DB, og models.Attachment, hash string) (bool, error) {
	prev, err := GetAttachmentByHash(hash)
	if err != nil {
		return false, err
	}

	if prev.PoolID != nil && og.PoolID != nil && prev.PoolID != og.PoolID && prev.Pool != nil && og.Pool != nil {
		if !prev.Pool.Config.Data().AllowCrossPoolEgress || !og.Pool.Config.Data().AllowCrossPoolIngress {
			// Pool config doesn't allow reference
			return false, nil
		}
	}

	if err := tx.Model(&og).Updates(&models.Attachment{
		RefID:       &prev.ID,
		Uuid:        prev.Uuid,
		Destination: prev.Destination,
		IsSelfRef:   og.AccountID == prev.AccountID,
	}).Error; err != nil {
		tx.Rollback()
		return true, err
	} else if err = tx.Model(&prev).Update("ref_count", prev.RefCount+1).Error; err != nil {
		tx.Rollback()
		return true, err
	}

	return true, nil
}

func UpdateAttachment(item models.Attachment) (models.Attachment, error) {
	if err := database.C.Save(&item).Error; err != nil {
		return item, err
	}

	return item, nil
}

func DeleteAttachment(item models.Attachment, txs ...*gorm.DB) error {
	dat := item

	var tx *gorm.DB
	if len(txs) == 0 {
		tx = database.C.Begin()
	} else {
		tx = txs[0]
	}

	if item.RefID != nil {
		var refTarget models.Attachment
		if err := database.C.Where("id = ?", *item.RefID).First(&refTarget).Error; err == nil {
			refTarget.RefCount--
			if err := tx.Save(&refTarget).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("unable to update ref count: %v", err)
			}
		}
	}
	if item.Thumbnail != nil {
		if err := DeleteAttachment(*item.Thumbnail, tx); err != nil {
			return err
		}
	}
	if item.Compressed != nil {
		if err := DeleteAttachment(*item.Compressed, tx); err != nil {
			return err
		}
	}
	if err := database.C.Delete(&item).Error; err != nil {
		tx.Rollback()
		return err
	} else {
		cachekit.Delete(gap.Ca, KgAttachmentCache(item.Rid))
	}

	tx.Commit()

	if dat.RefCount == 0 {
		go fs.DeleteFile(dat)
	}

	return nil
}

func DeleteAttachmentInBatch(items []models.Attachment, txs ...*gorm.DB) error {
	if len(items) == 0 {
		return nil
	}

	var tx *gorm.DB
	if len(txs) == 0 {
		tx = database.C.Begin()
	} else {
		tx = txs[0]
	}

	refIDs := []uint{}
	for _, item := range items {
		if item.RefID != nil {
			refIDs = append(refIDs, *item.RefID)
		}
	}

	if len(refIDs) > 0 {
		var refTargets []models.Attachment
		if err := tx.Where("id IN ?", refIDs).Find(&refTargets).Error; err == nil {
			for i := range refTargets {
				refTargets[i].RefCount--
			}
			if err := tx.Save(&refTargets).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("unable to update ref count: %v", err)
			}
		}
	}

	var subAttachments []models.Attachment
	for _, item := range items {
		if item.Thumbnail != nil {
			subAttachments = append(subAttachments, *item.Thumbnail)
		}
		if item.Compressed != nil {
			subAttachments = append(subAttachments, *item.Compressed)
		}
	}

	if len(subAttachments) > 0 {
		if err := DeleteAttachmentInBatch(subAttachments, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	rids := make([]string, len(items))
	for i, item := range items {
		rids[i] = item.Rid
	}

	if err := tx.Where("rid IN ?", rids).Delete(&models.Attachment{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, rid := range rids {
		cachekit.Delete(gap.Ca, KgAttachmentCache(rid))
	}

	tx.Commit()

	go func() {
		for _, item := range items {
			if item.RefCount == 0 {
				fs.DeleteFile(item)
			}
		}
	}()

	return nil
}

func CountAttachmentUsage(tx *gorm.DB, delta int) (int64, error) {
	if tx := tx.Model(&models.Attachment{}).
		Update("used_count", gorm.Expr("used_count + ?", delta)); tx.Error != nil {
		return tx.RowsAffected, tx.Error
	} else {
		return tx.RowsAffected, nil
	}
}
