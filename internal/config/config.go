package config

import "os"

const defaultHTTPAddress = ":8090"

type Config struct {
	HTTPAddress string
}

func Load() Config {
	address := os.Getenv("KIDASIE_HTTP_ADDRESS")
	if address == "" {
		address = defaultHTTPAddress
	}

	return Config{
		HTTPAddress: address,
	}
}
