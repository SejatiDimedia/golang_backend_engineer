package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port           string
	Env            string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPass         string
	DBName         string
	DBSSL          string
	JWTSecret      string
	JWTExpiryHours int
}

// LoadConfig memuat konfigurasi dari environment variables
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
		dbName = "booking_db"
	}

	dbSSL := os.Getenv("DB_SSLMODE")
	if dbSSL == "" {
		dbSSL = "disable"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "super_secret_signing_key_change_me_in_production"
	}

	jwtExpiryHours := 24
	if val := os.Getenv("JWT_EXPIRY_HOURS"); val != "" {
		if hours, err := strconv.Atoi(val); err == nil && hours > 0 {
			jwtExpiryHours = hours
		}
	}

	return &Config{
		Port:           port,
		Env:            env,
		DBHost:         dbHost,
		DBPort:         dbPort,
		DBUser:         dbUser,
		DBPass:         dbPass,
		DBName:         dbName,
		DBSSL:          dbSSL,
		JWTSecret:      jwtSecret,
		JWTExpiryHours: jwtExpiryHours,
	}, nil
}

// GetDSN menghasilkan Data Source Name untuk PostgreSQL
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		c.DBHost, c.DBUser, c.DBPass, c.DBName, c.DBPort, c.DBSSL,
	)
}
