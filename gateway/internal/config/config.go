package config

import "os"

type Config struct {
	GatewayPort          string
	AuthHost             string
	AuthPort             string
	MLHost               string
	MLPort               string
	MinioAccessKeyID     string
	MinioSecretAccessKey string
	MinioHost            string
	MinioPort            string
	MinioBucketName      string
	JwtSecretKey         string
	JwtExpired           string
}

func Load() *Config {
	return &Config{
		GatewayPort:          os.Getenv("GATEWAY_SERVICE_PORT"),
		AuthHost:             os.Getenv("AUTH_SERVICE_HOST"),
		AuthPort:             os.Getenv("AUTH_SERVICE_PORT"),
		MinioAccessKeyID:     os.Getenv("MINIO_ROOT_USER"),
		MinioSecretAccessKey: os.Getenv("MINIO_ROOT_PASSWORD"),
		MLHost:               os.Getenv("ML_SERVICE_HOST"),
		MLPort:               os.Getenv("ML_SERVICE_PORT"),
		MinioHost:            os.Getenv("MINIO_API_HOST"),
		MinioPort:            os.Getenv("MINIO_API_PORT"),
		MinioBucketName:      os.Getenv("MINIO_BUCKET_NAME"),
		JwtSecretKey:         os.Getenv("JWT_SECRET_KEY"),
		JwtExpired:           os.Getenv("JWT_EXPIRED"),
	}
}
