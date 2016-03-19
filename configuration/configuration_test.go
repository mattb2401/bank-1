package configuration

import (
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig()
	if err != nil {
		t.Errorf("LoadConfig does not pass. Looking for %v, got %v", nil, err)
	}

	if reflect.TypeOf(config).String() != "configuration.Configuration" {
		t.Errorf("LoadConfig does not pass. Looking for %v, got %v", "configuration.Configuration", reflect.TypeOf(config).String())
	}
}
