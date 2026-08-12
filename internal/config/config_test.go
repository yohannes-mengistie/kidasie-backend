package config

import (
	"testing"
)
func TestLoadUsesDefaultHTTPAddress(t *testing.T){
	t.Setenv("KIDASIE_HTTP_ADDRESS","")
	config := Load()

	if config.HTTPAddress != defaultHTTPAddress {
		t.Errorf("Expected default HTTP address %s, but got %s", defaultHTTPAddress, config.HTTPAddress)
	}
}

func TestLoadUsesHTTPAddressFromEnvironment(t *testing.T){
	const expectedAddress = ":8090"
	t.Setenv("KIDASIE_HTTP_ADDRESS",expectedAddress)
	config := Load()

	if config.HTTPAddress != expectedAddress{
		t.Errorf("expected HTTP address %s, but got %s",expectedAddress,config.HTTPAddress)
	}
}