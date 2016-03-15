package accounts

import (
	"fmt"
	"testing"

	"github.com/ksred/bank/configuration"
)

func TestLoadConfiguration(t *testing.T) {
	// Load app config
	Config, err := configuration.LoadConfig()
	if err != nil {
		t.Errorf("loadDatabase does not pass. Configuration does not load, looking for %v, got %v", nil, err)
	}
	fmt.Println(Config)
}

func TestLoadDatabase(t *testing.T) {
	/*
		var Config configuration.Configuration
		// Load app config
		Config, err = configuration.LoadConfig()
		if err != nil {
			t.Errorf("loadDatabase does not pass. Configuration does not load, looking for %v, got %v", nil, err)
		}
		fmt.Println(&Config)
			// Set config in packages
			SetConfig(&Config)

			_, err := loadDatabase()

			if err == nil {
				t.Errorf("LoadDatabase does not pass. Looking for %v, got %v", nil, err)
			}
	*/
}
