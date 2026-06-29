package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port                     string
	Env                      string
	DBHost                   string
	DBPort                   string
	DBUser                   string
	DBPass                   string
	DBName                   string
	DBSSL                    string
	RedisHost                string
	RedisPort                string
	RedisPassword            string
	RedisDB                  int
	AccessTokenExpiryMinutes int
	RefreshTokenExpiryDays   int
	RSAPrivateKeyPath        string
	RSAPublicKeyPath         string
}

func LoadConfig() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
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
		dbName = "auth_db"
	}

	dbSSL := os.Getenv("DB_SSLMODE")
	if dbSSL == "" {
		dbSSL = "disable"
	}

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisDB := 0
	if val := os.Getenv("REDIS_DB"); val != "" {
		if db, err := strconv.Atoi(val); err == nil {
			redisDB = db
		}
	}

	accessTokenExpiry := 15
	if val := os.Getenv("ACCESS_TOKEN_EXPIRY_MINUTES"); val != "" {
		if min, err := strconv.Atoi(val); err == nil && min > 0 {
			accessTokenExpiry = min
		}
	}

	refreshTokenExpiry := 7
	if val := os.Getenv("REFRESH_TOKEN_EXPIRY_DAYS"); val != "" {
		if days, err := strconv.Atoi(val); err == nil && days > 0 {
			refreshTokenExpiry = days
		}
	}

	privateKeyPath := os.Getenv("RSA_PRIVATE_KEY_PATH")
	if privateKeyPath == "" {
		privateKeyPath = "certs/private.key"
	}

	publicKeyPath := os.Getenv("RSA_PUBLIC_KEY_PATH")
	if publicKeyPath == "" {
		publicKeyPath = "certs/public.key"
	}

	return &Config{
		Port:                     port,
		Env:                      env,
		DBHost:                   dbHost,
		DBPort:                   dbPort,
		DBUser:                   dbUser,
		DBPass:                   dbPass,
		DBName:                   dbName,
		DBSSL:                    dbSSL,
		RedisHost:                redisHost,
		RedisPort:                redisPort,
		RedisPassword:            redisPassword,
		RedisDB:                  redisDB,
		AccessTokenExpiryMinutes: accessTokenExpiry,
		RefreshTokenExpiryDays:   refreshTokenExpiry,
		RSAPrivateKeyPath:        privateKeyPath,
		RSAPublicKeyPath:         publicKeyPath,
	}, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		c.DBHost, c.DBUser, c.DBPass, c.DBName, c.DBPort, c.DBSSL,
	)
}

func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}
