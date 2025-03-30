package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	nurl "net/url"
	"path/filepath"
	"strings"
	"time"

	"git.solsynth.dev/hypernet/nexus/pkg/nex/cachekit"
	"git.solsynth.dev/hypernet/paperclip/pkg/filekit/models"
	"git.solsynth.dev/hypernet/paperclip/pkg/internal/database"
	"git.solsynth.dev/hypernet/paperclip/pkg/internal/gap"
	jsoniter "github.com/json-iterator/go"
	"github.com/minio/minio-go/v7"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

type openAttachmentResult struct {
	Attachment models.Attachment        `json:"attachment"`
	Boosts     []models.AttachmentBoost `json:"boost"`
}

func KgAttachmentOpenCache(rid string) string {
	return fmt.Sprintf("attachment-open#%s", rid)
}

func OpenAttachmentByRID(rid string, preview bool, region ...string) (url string, mimetype string, err error) {
	var result *openAttachmentResult
	if val, err := cachekit.Get[openAttachmentResult](
		gap.Ca,
		KgAttachmentOpenCache(rid),
	); err == nil {
		result = &val
	}

	if result == nil {
		var attachment models.Attachment
		if err = database.C.Where(models.Attachment{
			Rid: rid,
		}).
			Preload("Pool").
			Preload("Thumbnail").
			Preload("Compressed").
			First(&attachment).Error; err != nil {
			return
		}

		var boosts []models.AttachmentBoost
		boosts, err = ListBoostByAttachmentWithStatus(attachment.ID, models.BoostStatusActive)
		if err != nil {
			return
		}

		result = &openAttachmentResult{
			Attachment: attachment,
			Boosts:     boosts,
		}
	}

	if len(result.Attachment.MimeType) > 0 {
		mimetype = result.Attachment.MimeType
	}

	var dest models.BaseDestination
	var rawDest []byte

	if len(region) > 0 {
		if des, ok := DestinationsByRegion[region[0]]; ok {
			for _, boost := range result.Boosts {
				if boost.Destination == des.Index {
					rawDest = des.Raw
					json.Unmarshal(rawDest, &dest)
				}
			}
		}
	}
	if rawDest == nil {
		if len(result.Boosts) > 0 {
			randomIdx := rand.IntN(len(result.Boosts))
			boost := result.Boosts[randomIdx]
			if des, ok := DestinationsByIndex[boost.Destination]; ok {
				rawDest = des.Raw
				json.Unmarshal(rawDest, &dest)
			}
		} else {
			if des, ok := DestinationsByIndex[result.Attachment.Destination]; ok {
				rawDest = des.Raw
				json.Unmarshal(rawDest, &dest)
			}
		}
	}

	if rawDest == nil {
		err = fmt.Errorf("no destination found")
		return
	}

	switch dest.Type {
	case models.DestinationTypeLocal:
		var destConfigured models.LocalDestination
		_ = jsoniter.Unmarshal(rawDest, &destConfigured)
		url = "file://" + filepath.Join(destConfigured.Path, result.Attachment.Uuid)
		return
	case models.DestinationTypeS3:
		var destConfigured models.S3Destination
		_ = jsoniter.Unmarshal(rawDest, &destConfigured)
		if destConfigured.EnableSigned {
			var client *minio.Client
			client, err = destConfigured.GetClient()
			if err != nil {
				return
			}

			var uri *nurl.URL
			uri, err = client.PresignedGetObject(context.Background(), destConfigured.Bucket, result.Attachment.Uuid, 60*time.Minute, nil)
			if err != nil {
				return
			}

			url = uri.String()
			return
		}
		if len(destConfigured.AccessBaseURL) > 0 {
			url = fmt.Sprintf(
				"%s/%s",
				destConfigured.AccessBaseURL,
				nurl.QueryEscape(filepath.Join(destConfigured.Path, result.Attachment.Uuid)),
			)
		} else {
			protocol := lo.Ternary(destConfigured.EnableSSL, "https", "http")
			url = fmt.Sprintf(
				"%s://%s.%s/%s",
				protocol,
				destConfigured.Bucket,
				destConfigured.Endpoint,
				nurl.QueryEscape(filepath.Join(destConfigured.Path, result.Attachment.Uuid)),
			)
		}
		if len(destConfigured.ImageProxyURL) > 0 && preview {
			size := viper.GetInt("imageproxy.size")
			url = fmt.Sprintf(
				"%s/%dx%d,fit/%s",
				destConfigured.ImageProxyURL,
				size,
				size,
				strings.Replace(url, destConfigured.AccessBaseURL, "", 1),
			)
		}
		return
	default:
		err = fmt.Errorf("invalid destination: unsupported protocol %s", dest.Type)
		return
	}
}

func CacheOpenAttachment(item *openAttachmentResult) {
	if item == nil {
		return
	}

	cachekit.Set[openAttachmentResult](
		gap.Ca,
		KgAttachmentCache(item.Attachment.Rid),
		*item,
		60*time.Minute,
	)
}
