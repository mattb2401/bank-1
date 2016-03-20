package appauth

import (
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ksred/bank/configuration"
	"github.com/satori/go.uuid"
	"gopkg.in/redis.v3"
)

const (
	TOKEN_TTL = time.Hour // One hour
)

var Config configuration.Configuration

func SetConfig(config *configuration.Configuration) {
	Config = *config
}

func ProcessAppAuth(data []string) (result string, err error) {
	//@TODO: Change from []string to something more solid, struct/interface/key-pair
	if len(data) < 3 {
		return "", errors.New("appauth.ProcessAppAuth: Not all required fields present")
	}
	switch data[2] {
	// Auth an existing account
	case "1":
		// TOKEN~appauth~1
		if len(data) < 3 {
			return "", errors.New("appauth.ProcessAppAuth: Not all required fields present")
		}
		err := CheckToken(data[0])
		if err != nil {
			return "", err
		}
		return result, nil
	// Log in
	case "2":
		if len(data) < 5 {
			return "", errors.New("appauth.ProcessAppAuth: Not all required fields present")
		}
		result, err = CreateToken(data[3], data[4])
		if err != nil {
			return "", err
		}
		return result, nil
	// Create an account
	case "3":
		if len(data) < 5 {
			return "", errors.New("appauth.ProcessAppAuth: Not all required fields present")
		}
		result, err = CreateUserPassword(data[3], data[4])
		if err != nil {
			return "", err
		}
		return result, nil
	// Remove an account
	case "4":
		if len(data) < 5 {
			return "", errors.New("appauth.ProcessAppAuth: Not all required fields present")
		}
		result, err = RemoveUserPassword(data[3], data[4])
		if err != nil {
			return "", err
		}
		return result, nil
	}
	return "", errors.New("appauth.ProcessAppAuth: No valid option chosen")
}

func CreateUserPassword(user string, password string) (result string, err error) {
	//TEST 0~appauth~3~181ac0ae-45cb-461d-b740-15ce33e4612f~testPassword
	// Generate hash
	hasher := sha512.New()
	hasher.Write([]byte(password))
	hash := hex.EncodeToString(hasher.Sum(nil))

	db, err := sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
	if err != nil {
		return "", errors.New("appauth.CreateUserPassword: Could not connect to database")
	}

	// Check for existing account
	rows, err := db.Query("SELECT `accountNumber` FROM `accounts_auth` WHERE `accountNumber` = ?", user)
	if err != nil {
		return "", errors.New("appauth.CreateUserPassword: Error with select query. " + err.Error())
	}
	defer rows.Close()

	// @TODO Must be easy way to get row count returned
	count := 0
	for rows.Next() {
		count++
	}

	if count > 0 {
		return "", errors.New("appauth.CreateUserPassword: Account already exists")
	}

	// Prepare statement for inserting data
	insertStatement := "INSERT INTO accounts_auth (`accountNumber`, `password`, `timestamp`) "
	insertStatement += "VALUES(?, ?, ?)"
	stmtIns, err := db.Prepare(insertStatement)
	if err != nil {
		return "", errors.New("appauth.CreateUserPassword: Error with insert. " + err.Error())
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// Convert variables
	t := time.Now()
	sqlTime := int32(t.Unix())

	_, err = stmtIns.Exec(user, hash, sqlTime)

	if err != nil {
		return "", errors.New("appauth.CreateUserPassword: Could not save account. " + err.Error())
	}

	result = "Successfully created account"
	return
}

func RemoveUserPassword(user string, hashedPassword string) (result string, err error) {
	db, err := sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
	if err != nil {
		return "", errors.New("appauth.RemoveUserPassword: Could not connect to database")
	}

	// Check for existing account
	rows, err := db.Query("SELECT `accountNumber` FROM `accounts_auth` WHERE `accountNumber` = ?", user)
	if err != nil {
		return "", errors.New("appauth.RemoveUserPassword: Error with select query. " + err.Error())
	}
	defer rows.Close()

	// @TODO Must be easy way to get row count returned
	count := 0
	for rows.Next() {
		count++
	}

	if count == 0 {
		return "", errors.New("appauth.RemoveUserPassword: Account auth does not exists")
	}

	// Prepare statement for inserting data
	delStatement := "DELETE FROM accounts_auth WHERE `accountNumber` = ? AND `password` = ? "
	stmtDel, err := db.Prepare(delStatement)
	if err != nil {
		return "", errors.New("appauth.RemoveUserPassword: Error with delete. " + err.Error())
	}
	defer stmtDel.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtDel.Exec(user, hashedPassword)

	if err != nil {
		return "", errors.New("appauth.RemoveUserPassword: Could not delete account. " + err.Error())
	}

	result = "Successfully deleted account"
	return
}

func CreateToken(user string, password string) (token string, err error) {
	// Check if username and password match
	db, err := sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
	if err != nil {
		return "", errors.New("appauth.CreateToken: Could not connect to database")
	}

	rows, err := db.Query("SELECT `password` FROM `accounts_auth` WHERE `accountNumber` = ?", user)
	if err != nil {
		return "", errors.New("appauth.CreateToken: Error with select query. " + err.Error())
	}
	defer rows.Close()

	count := 0
	hashedPassword := ""
	for rows.Next() {
		if err := rows.Scan(&hashedPassword); err != nil {
			return "", errors.New("appauth.CreateToken: Could not retreive account details")
		}
		count++
	}

	// Generate hash
	hasher := sha512.New()
	hasher.Write([]byte(password))
	hash := hex.EncodeToString(hasher.Sum(nil))

	if hash != hashedPassword {
		return "", errors.New("appauth.CreateToken: Authentication credentials invalid")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     Config.RedisHost + ":" + Config.RedisPort,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	newUuid := uuid.NewV4()
	token = newUuid.String()

	// @TODO Remove all tokens for this user
	err = client.Set(token, user, TOKEN_TTL).Err()
	if err != nil {
		return "", errors.New("appauth.CreateToken: Could not set token. " + err.Error())
	}

	return
}

func RemoveToken(token string) (result string, err error) {
	//TEST 0~appauth~480e67e3-e2c9-48ee-966c-8d251474b669
	client := redis.NewClient(&redis.Options{
		Addr:     Config.RedisHost + ":" + Config.RedisPort,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err = client.Del(token).Result()

	if err == redis.Nil {
		return "", errors.New("appauth.RemoveToken: Token not found. " + err.Error())
	} else if err != nil {
		return "", errors.New("appauth.RemoveToken: Could not remove token. " + err.Error())
	} else {
		result = "Token removed"
	}

	return
}

func CheckToken(token string) (err error) {
	//TEST 0~appauth~480e67e3-e2c9-48ee-966c-8d251474b669
	client := redis.NewClient(&redis.Options{
		Addr:     Config.RedisHost + ":" + Config.RedisPort,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	user, err := client.Get(token).Result()

	if err == redis.Nil {
		return errors.New("appauth.CheckToken: Token not found. " + err.Error())
	} else if err != nil {
		return errors.New("appauth.CheckToken: Could not get token. " + err.Error())
	} else {
		// Extend token
		err := client.Set(user, token, TOKEN_TTL).Err()
		if err != nil {
			return errors.New("appauth.CheckToken: Could not extend token. " + err.Error())
		}
	}

	return
}

func GetUserFromToken(token string) (user string, err error) {
	//TEST 0~appauth~~181ac0ae-45cb-461d-b740-15ce33e4612f~testPassword
	client := redis.NewClient(&redis.Options{
		Addr:     Config.RedisHost + ":" + Config.RedisPort,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	user, err = client.Get(token).Result()
	if err != nil {
		return "", errors.New("appauth.GetUserFromToken: Could not get token. " + err.Error())
	}

	// If valid then extend
	if user != "" {
		err := client.Set(user, token, TOKEN_TTL).Err()
		if err != nil {
			return "", errors.New("appauth.GetUserFromToken: Could not extend token. " + err.Error())
		}
	}

	return
}
