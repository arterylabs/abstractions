package storage

import (
	"abstractions/config"
	"abstractions/logger"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// ErrFileTooLarge is returned when an uploaded file exceeds maxUploadBytes.
var ErrFileTooLarge = errors.New("file exceeds maximum allowed size of 5 MB")

// ErrMIMENotAllowed is returned when an uploaded file's MIME type is not in allowedMIMETypes.
var ErrMIMENotAllowed = errors.New("file type is not allowed")

// maxUploadBytes is the maximum size accepted for a public upload (5 MB).
const maxUploadBytes = 5 << 20

// allowedMIMETypes is the set of MIME types accepted for public uploads.
// SVG is intentionally excluded due to XSS risk.
var allowedMIMETypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

type MinioStorageClient struct {
	client *minio.Client
	ctx    context.Context
}

func NewMinioStorageClient() (*MinioStorageClient, error) {
	client, err := minio.New(config.Environment.AWSStorageEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.Environment.AWSAccessKey, config.Environment.AWSSecretKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}

	return &MinioStorageClient{
		client: client,
		ctx:    context.Background(),
	}, nil
}

func (m *MinioStorageClient) Upload(bucket, objectName string, reader io.Reader, objectSize int64) (string, error) {
	inf, err := m.client.PutObject(m.ctx, bucket, objectName, reader, objectSize, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}
	return inf.Location, nil
}

func (m *MinioStorageClient) Get(bucket, objectName string) (io.Reader, error) {
	object, err := m.client.GetObject(m.ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (m *MinioStorageClient) Delete(bucket, objectName string) error {
	return m.client.RemoveObject(m.ctx, bucket, objectName, minio.RemoveObjectOptions{})
}

func (m *MinioStorageClient) UploadPublic(ctx context.Context, objectPath string, file multipart.File) (string, error) {
	// Cap the read to maxUploadBytes+1 so we can detect oversized files without
	// reading the entire payload into memory.
	limited := io.LimitReader(file, maxUploadBytes+1)
	fileContent, err := io.ReadAll(limited)
	if err != nil {
		logger.Error(ctx, err)
		return "", err
	}
	if int64(len(fileContent)) > maxUploadBytes {
		return "", ErrFileTooLarge
	}

	mimeType, err := mimetype.DetectReader(bytes.NewReader(fileContent))
	if err != nil {
		logger.Error(ctx, err)
		return "", err
	}
	if !allowedMIMETypes[mimeType.String()] {
		return "", fmt.Errorf("%w: %q", ErrMIMENotAllowed, mimeType.String())
	}

	objectName := fmt.Sprintf("%s/%s%s", objectPath, uuid.NewString(), mimeType.Extension())

	_, err = m.client.PutObject(m.ctx, config.Environment.AWSStorageBucket, objectName, bytes.NewReader(fileContent), int64(len(fileContent)), minio.PutObjectOptions{
		ContentType: mimeType.String(),
		UserMetadata: map[string]string{
			"x-amz-acl": "public-read",
		},
		CacheControl: "max-age=31536000",
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%s/%s/%s", config.Environment.AWSStorageEndpoint, config.Environment.AWSStorageBucket, objectName), nil
}
