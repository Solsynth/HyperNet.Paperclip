package models

import (
	"time"

	"git.solsynth.dev/hypernet/nexus/pkg/nex/cachekit"
	"git.solsynth.dev/hypernet/nexus/pkg/nex/cruda"
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"

	"git.solsynth.dev/hypernet/paperclip/pkg/internal/gap"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	AttachmentDstExternal  = -1 // The destination marked the file did not store inside our service
	AttachmentDstTemporary = 0  // Destination 0 is a reserved config for pre-upload processing
)

const (
	AttachmentTypeNormal = iota
	AttachmentTypeThumbnail
	AttachmentTypeCompressed
)

type Attachment struct {
	cruda.BaseModel

	// Random ID is for accessing (appear in URL)
	Rid string `json:"rid" gorm:"uniqueIndex"`
	// Unique ID is for storing (appear in local file name or object name)
	Uuid string `json:"uuid"`

	Size        int64  `json:"size"`
	Name        string `json:"name"`
	Alternative string `json:"alt"`
	MimeType    string `json:"mimetype"`
	HashCode    string `json:"hash"`
	Destination int    `json:"destination"`
	RefCount    int    `json:"ref_count"`
	Type        uint   `json:"type"`

	CleanedAt *time.Time `json:"cleaned_at"`

	Metadata datatypes.JSONMap `json:"metadata"` // This field is analyzer auto generated metadata
	Usermeta datatypes.JSONMap `json:"usermeta"` // This field is user set metadata

	ContentRating int `json:"content_rating"` // This field use to filter mature content or not
	QualityRating int `json:"quality_rating"` // This field use to filter good content or not

	IsAnalyzed  bool `json:"is_analyzed"`
	IsSelfRef   bool `json:"is_self_ref"`
	IsIndexable bool `json:"is_indexable"` // Show this attachment in the public directory api or not

	// Count the usage of this attachment across all services
	// If this number remain 0, it will be deleted during the maintenance
	UsedCount int `json:"used_count"`

	Thumbnail    *Attachment `json:"thumbnail"`
	ThumbnailID  *uint       `json:"thumbnail_id"`
	Compressed   *Attachment `json:"compressed"`
	CompressedID *uint       `json:"compressed_id"`

	Ref    *Attachment `json:"ref"`
	RefID  *uint       `json:"ref_id"`
	RefURL *string     `json:"ref_url"` // External URL for the attachment

	Pool   *AttachmentPool `json:"pool"`
	PoolID *uint           `json:"pool_id"`

	Boosts []AttachmentBoost `json:"boosts"`

	AccountID uint           `json:"account_id"`
	Account   models.Account `gorm:"-" json:"account"`
}

func (v *Attachment) AfterUpdate(tx *gorm.DB) error {
	_ = cachekit.Delete(
		gap.Ca,
		cachekit.FKey(cachekit.DAAttachment, v.Rid),
	)

	return nil
}

// Data model for in progress multipart attachments
type AttachmentFragment struct {
	cruda.BaseModel

	// Random ID is for accessing (appear in URL)
	Rid string `json:"rid" gorm:"uniqueIndex"`
	// Unique ID is for storing (appear in local file name or object name)
	Uuid string `json:"uuid"`

	Size        int64   `json:"size"`
	Name        string  `json:"name"`
	Alternative string  `json:"alt"`
	MimeType    string  `json:"mimetype"`
	HashCode    string  `json:"hash"`
	Fingerprint *string `json:"fingerprint"` // Client side generated hash, for continue uploading

	FileChunks datatypes.JSONMap `json:"file_chunks"`

	Metadata datatypes.JSONMap `json:"metadata"` // This field is analyzer auto generated metadata
	Usermeta datatypes.JSONMap `json:"usermeta"` // This field is user set metadata

	Pool   *AttachmentPool `json:"pool"`
	PoolID *uint           `json:"pool_id"`

	AccountID uint `json:"account_id"`

	FileChunksMissing []string `json:"file_chunks_missing" gorm:"-"` // This field use to prompt client which chunks is pending upload, do not store it
}

func (v *AttachmentFragment) AfterUpdate(tx *gorm.DB) error {
	_ = cachekit.Delete(
		gap.Ca,
		cachekit.FKey("attachment-fragment", v.Rid),
	)

	return nil
}

func (v *AttachmentFragment) ToAttachment() Attachment {
	return Attachment{
		Rid:         v.Rid,
		Uuid:        v.Uuid,
		Size:        v.Size,
		Name:        v.Name,
		Alternative: v.Alternative,
		MimeType:    v.MimeType,
		HashCode:    v.HashCode,
		Metadata:    v.Metadata,
		Usermeta:    v.Usermeta,
		Destination: AttachmentDstTemporary,
		Type:        AttachmentTypeNormal,
		Pool:        v.Pool,
		PoolID:      v.PoolID,
		AccountID:   v.AccountID,
	}
}
