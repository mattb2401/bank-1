package configuration

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
	redis "gopkg.in/redis.v3"
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
	Db           *sql.DB
	Redis        *redis.Client
}

var configPath = "./config.json"

func LoadConfig() (configuration Configuration, err error) {
	// Get config
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return Configuration{}, errors.New("configutation.LoadConfig: Could not load config. " + err.Error())
	}
	err = json.Unmarshal([]byte(file), &configuration)
	if err != nil {
		return Configuration{}, errors.New("configutation.LoadConfig: Could not load config. " + err.Error())
	}
	// Load MySQL
	err = loadMySQL(&configuration)
	if err != nil {
		return Configuration{}, errors.New("configuration.LoadConfig: Could not load MySQL. " + err.Error())
	}
	// Load Redis
	loadRedis(&configuration)
	return
}

func loadMySQL(configuration *Configuration) (err error) {
	configuration.Db, err = sql.Open("mysql", configuration.MySQLUser+":"+configuration.MySQLPass+"@tcp("+configuration.MySQLHost+":"+configuration.MySQLPort+")/"+configuration.MySQLDB)
	if err != nil {
		return errors.New("configuration.loadMySQL: Could not connect to database")
	}

	return
}

func loadRedis(configuration *Configuration) {
	configuration.Redis = redis.NewClient(&redis.Options{
		Addr:     configuration.RedisHost + ":" + configuration.RedisPort,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

}
