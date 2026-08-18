package config

import (
	"fmt"
	"os"
	"strings"
)

const defaultHTTPAddress = ":8090"

type Config struct {
	HTTPAddress string
	DatabaseURL string
}

func Load() Config {
	address := os.Getenv("KIDASIE_HTTP_ADDRESS")
	databaseUrl := os.Getenv("DATABASE_URL")
	if address == "" {
		address = defaultHTTPAddress
	}

	return Config{
		HTTPAddress: address,
		DatabaseURL: databaseUrl,
	}
}

func (c *Config) Validate() error {
	if strings.TrimSpace(c.DatabaseURL) == "" {
		return fmt.Errorf("database URL is required")
	}
	return nil
}