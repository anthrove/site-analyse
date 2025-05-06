package util

import (
	"context"
	"errors"
	"github.com/anthrove/site-analyse/pkg/e621"
	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
	"os"
)

func DownloadE6File(ctx context.Context, minioClient *minio.Client, bucketName string, filename string) error {
	name := filename[:len(filename)-3]

	if stat, err := os.Stat(name); err == nil {
		// path/to/whatever exists
		if stat.IsDir() {
			return errors.New("filename is a directory")
		} else {
			return nil
		}
	}

	err := minioClient.FGetObject(ctx, bucketName, name, name, minio.GetObjectOptions{})

	if err == nil {
		return nil
	}
	// Problem downloading file from s3 storage

	log.WithError(err).WithField("filename", filename).Warn("Failed to download file from s3 storage")

	err = e621.DownloadData(ctx, filename, name)
	if err != nil {
		return err
	}

	_, err = minioClient.FPutObject(ctx, bucketName, name, name, minio.PutObjectOptions{})

	if err != nil {
		return err
	}

	return nil
}
