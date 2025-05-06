package config

type S3StorageConfig struct {
	Endpoint        string `env:"S3_ENDPOINT,required"`
	AccessKeyID     string `env:"S3_ACCESS_KEY,required"`
	SecretAccessKey string `env:"S3_SECRET_KEY,required"`
	SSL             bool   `env:"S3_SSL,required"`
	BucketName      string `env:"S3_BUCKET_NAME,required"`
}
