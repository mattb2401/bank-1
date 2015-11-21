package configuration

import (
	"encoding/json"
	"fmt"
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

func LoadConfig() (configuration Configuration) {
	// Get config
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}

	return
}
