package config

import "os"

const defaultHTTPAddress = ":8090"

type Config struct {
	HTTPAddress string
	DatabaseUrl string
}

func Load() Config {
	address := os.Getenv("KIDASIE_HTTP_ADDRESS")
	databaseUrl := os.Getenv("DATABASE_URL")
	if address == "" {
		address = defaultHTTPAddress
	}

	return Config{
		HTTPAddress: address,
		DatabaseUrl: databaseUrl,
	}
}
