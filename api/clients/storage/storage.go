package storage

import (
	"abstractions/config"
	"abstractions/logger"
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageClient interface {
	Upload(bucket, objectName string, reader io.Reader, objectSize int64) error
	Get(bucket, objectName string) (io.Reader, error)
	Delete(bucket, objectName string) error
	UploadPublic(bucket, objectName string, reader io.Reader, objectSize int64) error
}

type MinioConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
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

	fileContent, err := io.ReadAll(file)
	if err != nil {
		logger.Error(ctx, err)
		return "", err
	}

	fileReader := bytes.NewReader(fileContent)
	fileReaderCopy := bytes.NewReader(fileContent)

	mime, err := mimetype.DetectReader(fileReaderCopy)
	if err != nil {
		logger.Error(ctx, err)
		return "", err
	}

	objectName := fmt.Sprintf("%s/%s", objectPath, uuid.NewString()+mime.Extension())

	_, err = m.client.PutObject(m.ctx, config.Environment.AWSStorageBucket, objectName, fileReader, int64(fileReader.Len()), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
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
