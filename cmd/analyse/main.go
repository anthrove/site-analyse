package main

import (
	"context"
	"github.com/anthrove/site-analyse/internal/analyze"
	"github.com/anthrove/site-analyse/pkg/config"
	"github.com/anthrove/site-analyse/pkg/util"
	"github.com/caarlos0/env/v11"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/prometheus/client_golang/prometheus/push"
	log "github.com/sirupsen/logrus"

	_ "github.com/joho/godotenv/autoload"
)

func main() {

	var s3Config config.S3StorageConfig
	if err := env.Parse(&s3Config); err != nil {
		log.Fatal(err)
	}

	var promConfig config.PrometheusConfig
	if err := env.Parse(&promConfig); err != nil {
		log.Fatal(err)
	}

	minioClient, err := minio.New(s3Config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3Config.AccessKeyID, s3Config.SecretAccessKey, ""),
		Secure: s3Config.SSL,
	})

	if err != nil {
		log.Fatalln(err)
	}

	name, err := util.DownloadE6File(context.Background(), minioClient, s3Config.BucketName, "tags-2025-05-06.csv.gz")

	if err != nil {
		log.Fatalln(err)
	}

	promPusher := push.New(promConfig.URL, "site_analytics").BasicAuth(promConfig.Username, promConfig.Password)

	analyze.Tags(context.Background(), promPusher, name)

}
