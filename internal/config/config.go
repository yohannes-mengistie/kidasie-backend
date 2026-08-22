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
	address := strings.TrimSpace(os.Getenv("KIDASIE_HTTP_ADDRESS"))
	databaseURL := os.Getenv("DATABASE_URL")
	if address == "" {
		address = addressFromPort(os.Getenv("PORT"))
	}

	return Config{
		HTTPAddress: address,
		DatabaseURL: databaseURL,
	}
}

func addressFromPort(port string) string {
	port = strings.TrimSpace(port)
	if port == "" {
		return defaultHTTPAddress
	}
	if strings.HasPrefix(port, ":") {
		return port
	}

	return ":" + port
}

func (c *Config) Validate() error {
	if strings.TrimSpace(c.DatabaseURL) == "" {
		return fmt.Errorf("database URL is required")
	}
	return nil
}
