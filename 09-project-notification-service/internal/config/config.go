package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port              string
	Env               string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPass            string
	DBName            string
	DBSSL             string
	RedisHost         string
	RedisPort         string
	RedisPassword     string
	RedisDB           int
	WorkerConcurrency int
	ProviderFailureRate float64
	JWTSecret         string
	JWTExpiryHours    int
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
		dbName = "notification_db"
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

	concurrency := 5
	if val := os.Getenv("WORKER_CONCURRENCY"); val != "" {
		if c, err := strconv.Atoi(val); err == nil && c > 0 {
			concurrency = c
		}
	}

	failureRate := 0.3
	if val := os.Getenv("PROVIDER_FAILURE_RATE"); val != "" {
		if rate, err := strconv.ParseFloat(val, 64); err == nil && rate >= 0 && rate <= 1 {
			failureRate = rate
		}
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "super_secret_notification_signing_key_change_me"
	}

	jwtExpiryHours := 24
	if val := os.Getenv("JWT_EXPIRY_HOURS"); val != "" {
		if hours, err := strconv.Atoi(val); err == nil && hours > 0 {
			jwtExpiryHours = hours
		}
	}

	return &Config{
		Port:              port,
		Env:               env,
		DBHost:            dbHost,
		DBPort:            dbPort,
		DBUser:            dbUser,
		DBPass:            dbPass,
		DBName:            dbName,
		DBSSL:             dbSSL,
		RedisHost:         redisHost,
		RedisPort:         redisPort,
		RedisPassword:     redisPassword,
		RedisDB:           redisDB,
		WorkerConcurrency: concurrency,
		ProviderFailureRate: failureRate,
		JWTSecret:         jwtSecret,
		JWTExpiryHours:    jwtExpiryHours,
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
