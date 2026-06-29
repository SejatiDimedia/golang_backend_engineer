package utils

import (
	"context"
	"io"
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageClient interface {
	UploadFile(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, contentType string) error
	GetPresignedURL(ctx context.Context, bucketName, objectName string, expiry time.Duration) (*url.URL, error)
	GetFileStream(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error)
	DeleteFile(ctx context.Context, bucketName, objectName string) error
	Ping(ctx context.Context) error
}

type minioStorageClient struct {
	client *minio.Client
}

func NewMinioStorageClient(endpoint, accessKey, secretKey string, useSSL bool) (StorageClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &minioStorageClient{client: client}, nil
}

func (c *minioStorageClient) UploadFile(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, contentType string) error {
	_, err := c.client.PutObject(ctx, bucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (c *minioStorageClient) GetPresignedURL(ctx context.Context, bucketName, objectName string, expiry time.Duration) (*url.URL, error) {
	reqParams := make(url.Values)
	return c.client.PresignedGetObject(ctx, bucketName, objectName, expiry, reqParams)
}

func (c *minioStorageClient) GetFileStream(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error) {
	return c.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
}

func (c *minioStorageClient) DeleteFile(ctx context.Context, bucketName, objectName string) error {
	return c.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}

func (c *minioStorageClient) Ping(ctx context.Context) error {
	_, err := c.client.ListBuckets(ctx)
	return err
}

// EnsureBucketExists checks if the MinIO bucket exists, creating it if it doesn't.
func EnsureBucketExists(ctx context.Context, endpoint, accessKey, secretKey string, useSSL bool, bucketName string) error {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return err
	}

	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if !exists {
		log.Printf("MinIO bucket '%s' not found. Creating it...", bucketName)
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
		log.Printf("MinIO bucket '%s' created successfully", bucketName)
	} else {
		log.Printf("MinIO bucket '%s' already exists", bucketName)
	}

	return nil
}
