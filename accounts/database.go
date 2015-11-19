package accounts

/*
@TODO Fix DB repetition
*/

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	//"github.com/satori/go.uuid"
	"os"
	"time"
)

type Configuration struct {
	TimeZone  string
	MySQLUser string
	MySQLPass string
	MySQLHost string
	MySQLPort string
	MySQLDB   string
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

func loadDatabase() (db *sql.DB) {
	configuration := Configuration{}
	loadConfig(&configuration)

	db, err := sql.Open("mysql", configuration.MySQLUser+":"+configuration.MySQLPass+"@tcp("+configuration.MySQLHost+":"+configuration.MySQLPort+")/"+configuration.MySQLDB)
	if err != nil {
		fmt.Println("Could not connect to database")
		return
	}
	defer db.Close()

	// Test connection with ping
	err = db.Ping()
	if err != nil {
		fmt.Println("Ping error: " + err.Error()) // proper error handling instead of panic in your app
		return
	}

	return
}

func createAccount(accountDetails AccountDetails, accountHolderDetails AccountHolderDetails) (newAccountDetails AccountDetails) {
	configuration := Configuration{}
	loadConfig(&configuration)

	db, err := sql.Open("mysql", configuration.MySQLUser+":"+configuration.MySQLPass+"@tcp("+configuration.MySQLHost+":"+configuration.MySQLPort+")/"+configuration.MySQLDB)
	if err != nil {
		fmt.Println("Could not connect to database")
		return
	}

	// Prepare statement for inserting data
	// Create account
	insertStatement := "INSERT INTO accounts (`accountNum`, `bankNum`, `accountHolderName`, `accountBalance`, `overdraft`, `availableBalance`, `timestamp`) "
	insertStatement += "VALUES(?, ?, ?, ?, ?, ?, ?)"
	stmtIns, err := db.Prepare(insertStatement)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// Convert variables
	t := time.Now()
	sqlTime := int32(t.Unix())

	// Generate account number
	//accountDetails.AccountNumber = uuid.NewV4() // @TODO Convert to string. https://github.com/satori/go.uuid/blob/master/uuid.go#L409
	accountDetails.AccountNumber = ""

	_, err = stmtIns.Exec(accountDetails.AccountNumber, accountDetails.BankNumber, accountDetails.AccountHolderName, accountDetails.Overdraft, accountDetails.AvailableBalance, sqlTime)

	if err != nil {
		fmt.Println("Could not save results: " + err.Error())
	}

	// Create account meta
	insertStatement = "INSERT INTO accounts_meta (`accountNumber`, `bankNumber`, `accountHolderGivenName`, `accountHolderFamilyName`, `accountHolderDateOfBirth`, `accountHolderIdentificationNumber`, `accountHolderContactNumber1`, `accountHolderContactNumber2`, `accountHolderEmailAddress`, `accountHolderAddressLine1`, `accountHolderAddressLine2`, `accountHolderAddressLine3`, `accountHolderPostalCode`, `timestamp`) "
	insertStatement += "VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	stmtIns, err = db.Prepare(insertStatement)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	accountHolderDetails.AccountNumber = accountDetails.AccountNumber

	_, err = stmtIns.Exec(accountHolderDetails.AccountNumber, accountHolderDetails.BankNumber, accountHolderDetails.GivenName, accountHolderDetails.FamilyName, accountHolderDetails.DateOfBirth, accountHolderDetails.IdentificationNumber, accountHolderDetails.ContactNumber1, accountHolderDetails.ContactNumber2, accountHolderDetails.EmailAddress, accountHolderDetails.AddressLine1, accountHolderDetails.AddressLine2, accountHolderDetails.AddressLine3,
		accountHolderDetails.PostalCode, sqlTime)

	if err != nil {
		fmt.Println("Could not save results: " + err.Error())
	}

	defer db.Close()

	return
}

func getAccountDetails(id string) (accountDetails AccountDetails) {
	configuration := Configuration{}
	loadConfig(&configuration)
	db, err := sql.Open("mysql", configuration.MySQLUser+":"+configuration.MySQLPass+"@tcp("+configuration.MySQLHost+":"+configuration.MySQLPort+")/"+configuration.MySQLDB)
	if err != nil {
		fmt.Println("Could not connect to database")
		return
	}

	rows, err := db.Query("SELECT `accountNumber`, `bankNumber`, `accountHolderName`, `accountBalance`, `overdraft`, `availableBalance` FROM `accounts` WHERE `accountNumber` = ?", id)
	if err != nil {
		fmt.Println("Error with select query: " + err.Error())
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		if err := rows.Scan(&accountDetails.AccountNumber, &accountDetails.BankNumber, &accountDetails.AccountHolderName, &accountDetails.AccountBalance, &accountDetails.Overdraft, &accountDetails.AvailableBalance); err != nil {
			//@TODO Throw error
			fmt.Println("ERROR: Could not retrieve account details")
			return
		}
		count++
	}

	if count == 0 {
		fmt.Println("Account not found")
		return
	}

	if count > 1 {
		//@TODO Throw error
		fmt.Println("ERROR: More than one account found")
		return
	}

	return
}

func getAccountMeta(id string) (accountDetails AccountHolderDetails) {
	configuration := Configuration{}
	loadConfig(&configuration)
	db, err := sql.Open("mysql", configuration.MySQLUser+":"+configuration.MySQLPass+"@tcp("+configuration.MySQLHost+":"+configuration.MySQLPort+")/"+configuration.MySQLDB)
	if err != nil {
		fmt.Println("Could not connect to database")
		return
	}

	rows, err := db.Query("SELECT `accountNumber`, `bankNumber`, `accountHolderGivenName`, `accountHolderFamilyName`, `accountHolderDateOfBirth`, `accountHolderIdentificationNumber`, `accountHolderContactNumber1`, `accountHolderContactNumber2`, `accountHolderEmailAddress`, `accountHolderAddressLine1`, `accountHolderAddressLine2`, `accountHolderAddressLine3`, `accountHolderPostalCode` FROM `accounts_meta` WHERE `accountHolderIdentificationNumber` = ?", id)
	if err != nil {
		fmt.Println("Error with select query: " + err.Error())
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		if err := rows.Scan(&accountDetails.AccountNumber, &accountDetails.BankNumber, &accountDetails.GivenName, &accountDetails.FamilyName, &accountDetails.DateOfBirth, &accountDetails.IdentificationNumber, &accountDetails.ContactNumber1, &accountDetails.ContactNumber2, &accountDetails.EmailAddress, &accountDetails.AddressLine1, &accountDetails.AddressLine2,
			&accountDetails.AddressLine3, &accountDetails.PostalCode); err != nil {
			//@TODO Throw error
			fmt.Println("ERROR: Could not retrieve account details")
			return
		}
		count++
	}

	if count == 0 {
		fmt.Println("Account not found")
		return
	}

	if count > 1 {
		//@TODO Allow for a customer to have more than one account
		fmt.Println("ERROR: More than one account found")
		return
	}

	return
}
