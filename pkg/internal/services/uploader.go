package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"git.solsynth.dev/hypernet/paperclip/pkg/filekit/models"
	"git.solsynth.dev/hypernet/paperclip/pkg/internal/database"
	"git.solsynth.dev/hypernet/paperclip/pkg/internal/fs"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
)

func UploadFileToTemporary(ctx *fiber.Ctx, file *multipart.FileHeader, meta models.Attachment) error {
	destMap := viper.GetStringMap(fmt.Sprintf("destinations.%d", meta.Destination))

	var dest models.BaseDestination
	rawDest, _ := jsoniter.Marshal(destMap)
	_ = jsoniter.Unmarshal(rawDest, &dest)

	switch dest.Type {
	case models.DestinationTypeLocal:
		var destConfigured models.LocalDestination
		_ = jsoniter.Unmarshal(rawDest, &destConfigured)
		return ctx.SaveFile(file, filepath.Join(destConfigured.Path, meta.Uuid))
	default:
		return fmt.Errorf("invalid destination: unsupported protocol %s", dest.Type)
	}
}

func ReUploadFile(meta models.Attachment, dst int, doNotUpdate ...bool) error {
	if meta.Destination == dst {
		return fmt.Errorf("destnation cannot be reversed temporary or the same as the original")
	}

	prevDst := meta.Destination
	inDst, err := fs.DownloadFileToLocal(meta, prevDst)
	if err != nil {
		return fmt.Errorf("unable to retrieve file content: %v", err)
	}

	cleanupDst := func() {
		if len(doNotUpdate) == 0 || !doNotUpdate[0] {
			database.C.Model(&meta).Update("destination", dst)
		}
		if prevDst == models.AttachmentDstTemporary {
			return
		}
		os.Remove(inDst)
	}

	meta.Destination = dst
	destMap := viper.GetStringMap(fmt.Sprintf("destinations.%d", dst))

	var dest models.BaseDestination
	rawDest, _ := jsoniter.Marshal(destMap)
	_ = jsoniter.Unmarshal(rawDest, &dest)

	switch dest.Type {
	case models.DestinationTypeLocal:
		var destConfigured models.LocalDestination
		_ = jsoniter.Unmarshal(rawDest, &destConfigured)

		in, err := os.Open(inDst)
		if err != nil {
			return fmt.Errorf("unable to open file in temporary storage: %v", err)
		}
		defer in.Close()

		out, err := os.Create(filepath.Join(destConfigured.Path, meta.Uuid))
		if err != nil {
			return fmt.Errorf("unable to open dest file: %v", err)
		}
		defer out.Close()

		_, err = io.Copy(out, in)
		if err != nil {
			return fmt.Errorf("unable to copy data to dest file: %v", err)
		}

		cleanupDst()
		return nil
	case models.DestinationTypeS3:
		var destConfigured models.S3Destination
		_ = jsoniter.Unmarshal(rawDest, &destConfigured)

		client, err := destConfigured.GetClient()
		if err != nil {
			return fmt.Errorf("unable to configure s3 client: %v", err)
		}

		_, err = client.FPutObject(context.Background(), destConfigured.Bucket, filepath.Join(destConfigured.Path, meta.Uuid), inDst, minio.PutObjectOptions{
			ContentType:           meta.MimeType,
			SendContentMd5:        false,
			DisableContentSha256:  true,
			PartSize:              10 * 1024 * 1024,
			ConcurrentStreamParts: true,
			NumThreads:            4,
		})
		if err != nil {
			return fmt.Errorf("unable to upload file to s3: %v", err)
		}

		cleanupDst()
		return nil
	default:
		return fmt.Errorf("invalid destination: unsupported protocol %s", dest.Type)
	}
}
