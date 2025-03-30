package models

import (
	pkg "git.solsynth.dev/hypernet/paperclip/pkg/internal"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	DestinationTypeLocal = "local"
	DestinationTypeS3    = "s3"
)

type BaseDestination struct {
	ID      int    `json:"id,omitempty"` // Auto filled with index, only for user
	Type    string `json:"type"`
	Label   string `json:"label"`
	Region  string `json:"region"`
	IsBoost bool   `json:"is_boost"`
}

type LocalDestination struct {
	BaseDestination

	Path          string `json:"path"`
	AccessBaseURL string `json:"access_baseurl"`
}

type S3Destination struct {
	BaseDestination

	Path          string `json:"path"`
	Bucket        string `json:"bucket"`
	Endpoint      string `json:"endpoint"`
	SecretID      string `json:"secret_id"`
	SecretKey     string `json:"secret_key"`
	AccessBaseURL string `json:"access_baseurl"`
	ImageProxyURL string `json:"image_proxy_baseurl"`
	EnableSSL     bool   `json:"enable_ssl"`
	EnableSigned  bool   `json:"enable_signed"`
	BucketLookup  int    `json:"bucket_lookup"`
}

func (v S3Destination) GetClient() (*minio.Client, error) {
	client, err := minio.New(v.Endpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(v.SecretID, v.SecretKey, ""),
		Secure:       v.EnableSSL,
		BucketLookup: minio.BucketLookupType(v.BucketLookup),
	})
	if err == nil {
		client.SetAppInfo("HyperNet.Paperclip", pkg.AppVersion)
	}
	return client, err
}
