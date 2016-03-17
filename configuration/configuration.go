package configuration

import (
	"encoding/json"
	"errors"
	"os"
)

type Configuration struct {
	TimeZone     string
	MySQLUser    string
	MySQLPass    string
	MySQLHost    string
	MySQLPort    string
	MySQLDB      string
	RedisHost    string
	RedisPort    string
	PasswordSalt string
}

var configPath = "/Users/ksred/golang/projects/src/github.com/ksred/bank/config.json"

func LoadConfig() (configuration Configuration, err error) {
	// Get config
	file, _ := os.Open(configPath)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		return Configuration{}, errors.New("configuration.LoadConfig: Could not load config. " + err.Error())
	}

	return
}
