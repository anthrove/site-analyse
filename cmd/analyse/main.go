package main

import (
	"context"
	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"github.com/anthrove/site-analyse/internal/analyze"
	"github.com/anthrove/site-analyse/pkg/config"
	"github.com/anthrove/site-analyse/pkg/util"
	"github.com/caarlos0/env/v11"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
	"time"
	_ "time/tzdata"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	var s3Config config.S3StorageConfig
	if err := env.Parse(&s3Config); err != nil {
		log.Fatal(err)
	}

	var influxDBConfig config.InfluxDB
	if err := env.Parse(&influxDBConfig); err != nil {
		log.Fatal(err)
	}

	minioClient, err := minio.New(s3Config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3Config.AccessKeyID, s3Config.SecretAccessKey, ""),
		Secure: s3Config.SSL,
	})

	if err != nil {
		log.Fatalln(err)
	}

	client, err := influxdb3.New(influxdb3.ClientConfig{
		Host:     influxDBConfig.Host,
		Database: influxDBConfig.Database,
		Token:    influxDBConfig.Token,
	})

	defer func(client *influxdb3.Client) {
		err := client.Close()
		if err != nil {
			panic(err)
		}
	}(client)

	if err != nil {
		log.Fatalln(err)
	}

	analzyeModules := map[string]func(ctx context.Context, influxClient *influxdb3.Client, fileName string) error{
		"tags":  analyze.Tags,
		"posts": analyze.Posts,
		"pools": analyze.Pools,
	}

	dateStr := time.Now().Format("2006-01-02")
	for module, function := range analzyeModules {
		log.WithField("module", module).Info("Starting analyzing")
		fileName := module + "-" + dateStr + ".csv.gz"
		name, err := util.DownloadE6File(context.Background(), minioClient, s3Config.BucketName, fileName)

		if err != nil {
			log.Fatalln(err)
		}
		log.WithField("module", module).WithField("file_name", fileName).Info("File downloaded")

		err = function(context.Background(), client, name)
		if err != nil {
			log.Fatalln(err)
		}
		log.WithField("module", module).Info("Finished analyzing")
	}
}
