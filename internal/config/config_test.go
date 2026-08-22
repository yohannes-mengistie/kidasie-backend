package config

import (
	"testing"
)

func TestLoadUsesDefaultHTTPAddress(t *testing.T) {
	t.Setenv("KIDASIE_HTTP_ADDRESS", "")
	t.Setenv("PORT", "")
	config := Load()

	if config.HTTPAddress != defaultHTTPAddress {
		t.Errorf("Expected default HTTP address %s, but got %s", defaultHTTPAddress, config.HTTPAddress)
	}
}

func TestLoadUsesRenderPort(t *testing.T) {
	t.Setenv("KIDASIE_HTTP_ADDRESS", "")
	t.Setenv("PORT", "10000")

	config := Load()

	if config.HTTPAddress != ":10000" {
		t.Errorf("expected HTTP address :10000, but got %s", config.HTTPAddress)
	}
}

func TestLoadPrefersKidasieHTTPAddress(t *testing.T) {
	t.Setenv("KIDASIE_HTTP_ADDRESS", ":9000")
	t.Setenv("PORT", "10000")

	config := Load()

	if config.HTTPAddress != ":9000" {
		t.Errorf("expected HTTP address :9000, but got %s", config.HTTPAddress)
	}
}
func TestLoadUsesHTTPAddressFromEnvironment(t *testing.T) {
	const expectedAddress = ":8090"
	t.Setenv("KIDASIE_HTTP_ADDRESS", expectedAddress)
	config := Load()

	if config.HTTPAddress != expectedAddress {
		t.Errorf("expected HTTP address %s, but got %s", expectedAddress, config.HTTPAddress)
	}
}
