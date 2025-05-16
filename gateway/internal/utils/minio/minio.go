package minio

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
)

type IMinio interface {
	UploadPDF(ctx context.Context,
		bucketName, objectName string,
		reader io.Reader, fileSize int64) (minio.UploadInfo, error)
	MakeBucket(ctx context.Context, bucketName string) error
}

type Minio struct {
	client *minio.Client
}

func NewMinio(client *minio.Client) IMinio {
	return &Minio{client: client}
}

func (m *Minio) MakeBucket(ctx context.Context, bucketName string) error {
	if err := m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
		exists, err := m.client.BucketExists(ctx, bucketName)
		if exists && err == nil {
			return nil
		} else {
			return fmt.Errorf("error creating new bucket: %v", err)
		}
	}

	log.Info().Msgf("successfully create %s bucket", bucketName)
	return nil
}

func (m *Minio) UploadPDF(ctx context.Context,
	bucketName, objectName string,
	reader io.Reader, fileSize int64) (minio.UploadInfo, error) {
	info, err := m.client.PutObject(ctx, bucketName, objectName, reader, fileSize, minio.PutObjectOptions{
		ContentType: "application/pdf",
	})
	if err != nil {
		return minio.UploadInfo{}, errors.New("failed upload pdf to minio object storage")
	}

	return info, nil
}
