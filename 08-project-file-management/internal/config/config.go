package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port                         string
	Env                          string
	DBHost                       string
	DBPort                       string
	DBUser                       string
	DBPass                       string
	DBName                       string
	DBSSL                        string
	MinioEndpoint                string
	MinioAccessKey               string
	MinioSecretKey               string
	MinioUseSSL                  bool
	MinioBucketName              string
	MinioPresignedExpiryMinutes  int
	JWTSecret                    string
	JWTExpiryHours               int
}

func LoadConfig() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}

	dbPass := os.Getenv("DB_PASSWORD")
	if dbPass == "" {
		dbPass = "postgres"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "file_db"
	}

	dbSSL := os.Getenv("DB_SSLMODE")
	if dbSSL == "" {
		dbSSL = "disable"
	}

	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	if minioEndpoint == "" {
		minioEndpoint = "localhost:9000"
	}

	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	if minioAccessKey == "" {
		minioAccessKey = "minioadmin"
	}

	minioSecretKey := os.Getenv("MINIO_SECRET_KEY")
	if minioSecretKey == "" {
		minioSecretKey = "minioadmin"
	}

	minioUseSSL := false
	if val := os.Getenv("MINIO_USE_SSL"); val == "true" {
		minioUseSSL = true
	}

	minioBucketName := os.Getenv("MINIO_BUCKET_NAME")
	if minioBucketName == "" {
		minioBucketName = "user-files"
	}

	minioExpiryMinutes := 15
	if val := os.Getenv("MINIO_PRESIGNED_EXPIRY_MINUTES"); val != "" {
		if mins, err := strconv.Atoi(val); err == nil && mins > 0 {
			minioExpiryMinutes = mins
		}
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "super_secret_file_management_signing_key_change_me"
	}

	jwtExpiryHours := 24
	if val := os.Getenv("JWT_EXPIRY_HOURS"); val != "" {
		if hours, err := strconv.Atoi(val); err == nil && hours > 0 {
			jwtExpiryHours = hours
		}
	}

	return &Config{
		Port:                         port,
		Env:                          env,
		DBHost:                       dbHost,
		DBPort:                       dbPort,
		DBUser:                       dbUser,
		DBPass:                       dbPass,
		DBName:                       dbName,
		DBSSL:                        dbSSL,
		MinioEndpoint:                minioEndpoint,
		MinioAccessKey:               minioAccessKey,
		MinioSecretKey:               minioSecretKey,
		MinioUseSSL:                  minioUseSSL,
		MinioBucketName:              minioBucketName,
		MinioPresignedExpiryMinutes:  minioExpiryMinutes,
		JWTSecret:                    jwtSecret,
		JWTExpiryHours:               jwtExpiryHours,
	}, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		c.DBHost, c.DBUser, c.DBPass, c.DBName, c.DBPort, c.DBSSL,
	)
}
