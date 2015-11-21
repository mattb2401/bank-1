package appauth

import (
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/satori/go.uuid"
	"gopkg.in/redis.v3"
	"os"
	"time"
)

const (
	TOKEN_TTL = 60 * 60 * 1000 // One hour
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

func loadConfig(configuration *Configuration) {
	// Get config
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
}

func ProcessAppAuth(data []string) (result string) {
	switch data[2] {
	// Auth an existing account
	case "1":
		if len(data) < 4 {
			result = "0~Not all required fields present"
		}
		CheckToken(data[3])
		break
		// Log in
	case "2":
		if len(data) < 5 {
			result = "0~Not all required fields present"
		}
		result = CreateToken(data[3], data[4])
		break
	// Create an account
	case "3":
		if len(data) < 5 {
			result = "0~Not all required fields present"
		}
		CreateUserPassword(data[3], data[4])
		break
	}

	return
}

func CreateUserPassword(user string, password string) (result string) {
	//TEST 0~appauth~2~181ac0ae-45cb-461d-b740-15ce33e4612f~testPassword
	configuration := Configuration{}
	loadConfig(&configuration)

	// Generate hash
	hasher := sha512.New()
	hasher.Write([]byte(password))
	hash := hex.EncodeToString(hasher.Sum(nil))

	db, err := sql.Open("mysql", configuration.MySQLUser+":"+configuration.MySQLPass+"@tcp("+configuration.MySQLHost+":"+configuration.MySQLPort+")/"+configuration.MySQLDB)
	if err != nil {
		fmt.Println("Could not connect to database")
		return
	}

	// Check for existing account
	rows, err := db.Query("SELECT `accountNumber` FROM `accounts_auth` WHERE `accountNumber` = ?", user)
	if err != nil {
		fmt.Println("Error with select query: " + err.Error())
	}
	defer rows.Close()

	// @TODO Must be easy way to get row count returned
	count := 0
	for rows.Next() {
		count++
	}

	if count > 0 {
		result = "0~Account already exists"
		return
	}

	// Prepare statement for inserting data
	insertStatement := "INSERT INTO accounts_auth (`accountNumber`, `password`, `timestamp`) "
	insertStatement += "VALUES(?, ?, ?)"
	stmtIns, err := db.Prepare(insertStatement)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// Convert variables
	t := time.Now()
	sqlTime := int32(t.Unix())

	_, err = stmtIns.Exec(user, hash, sqlTime)

	if err != nil {
		fmt.Println("Could not save account: " + err.Error())
	}

	result = "1~Successfully created account"
	return
}

func CreateToken(user string, password string) (token string) {
	configuration := Configuration{}
	loadConfig(&configuration)

	// Check if username and password match
	db, err := sql.Open("mysql", configuration.MySQLUser+":"+configuration.MySQLPass+"@tcp("+configuration.MySQLHost+":"+configuration.MySQLPort+")/"+configuration.MySQLDB)
	if err != nil {
		fmt.Println("Could not connect to database")
		return
	}

	rows, err := db.Query("SELECT `password` FROM `accounts_auth` WHERE `accountNumber` = ?", user)
	if err != nil {
		fmt.Println("Error with select query: " + err.Error())
	}
	defer rows.Close()

	count := 0
	hashedPassword := ""
	for rows.Next() {
		if err := rows.Scan(&hashedPassword); err != nil {
			//@TODO Throw error
			fmt.Println("ERROR: Could not retrieve account details")
			return
		}
		count++
	}

	// Generate hash
	hasher := sha512.New()
	hasher.Write([]byte(password))
	hash := hex.EncodeToString(hasher.Sum(nil))

	if hash != hashedPassword {
		token = "0~Authentication credentials invalid"
		return
	}

	client := redis.NewClient(&redis.Options{
		Addr:     configuration.RedisHost + ":" + configuration.RedisPort,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	newUuid := uuid.NewV4()
	token = newUuid.String()

	// @TODO Remove all tokens for this user
	err = client.Set(token, user, TOKEN_TTL).Err()
	if err != nil {
		panic(err)
	}

	return
}

func CheckToken(token string) (res bool) {
	//TEST 0~appauth~480e67e3-e2c9-48ee-966c-8d251474b669
	configuration := Configuration{}
	loadConfig(&configuration)

	client := redis.NewClient(&redis.Options{
		Addr:     configuration.RedisHost + ":" + configuration.RedisPort,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Get(token).Result()

	if err == redis.Nil {
		res = false
	} else if err != nil {
		panic(err)
	}

	res = true

	return
}

func GetUserFromToken(token string) (user string) {
	//TEST 0~appauth~~181ac0ae-45cb-461d-b740-15ce33e4612f~testPassword
	configuration := Configuration{}
	loadConfig(&configuration)

	client := redis.NewClient(&redis.Options{
		Addr:     configuration.RedisHost + ":" + configuration.RedisPort,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	user, err := client.Get(token).Result()
	if err != nil {
		panic(err)
	}

	// If valid then extend
	if user != "" {
		err := client.Set(user, token, TOKEN_TTL).Err()
		if err != nil {
			panic(err)
		}
	}

	return
}
