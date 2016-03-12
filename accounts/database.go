package accounts

/*
@TODO Fix DB repetition
*/

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ksred/bank/configuration"
	"github.com/satori/go.uuid"
)

var Config configuration.Configuration

func SetConfig(config *configuration.Configuration) {
	Config = *config
}

func loadDatabase() (db *sql.DB, err error) {
	db, err = sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
	if err != nil {
		return
	}
	defer db.Close()

	// Test connection with ping
	err = db.Ping()
	if err != nil {
		return
	}

	return
}

func createAccount(accountDetails AccountDetails, accountHolderDetails AccountHolderDetails) (AccountDetails AccountDetails, err error) {
	// Convert variables
	t := time.Now()
	sqlTime := int32(t.Unix())

	err = doCreateAccount(sqlTime, &accountDetails)
	if err != nil {
		return
	}

	err = doCreateAccountMeta(sqlTime, &accountHolderDetails, &accountDetails)
	if err != nil {
		return
	}

	return
}

func doCreateAccount(sqlTime int32, accountDetails *AccountDetails) (err error) {
	db, err := sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
	if err != nil {
		return
	}
	// Create account
	insertStatement := "INSERT INTO accounts (`accountNumber`, `bankNumber`, `accountHolderName`, `accountBalance`, `overdraft`, `availableBalance`, `timestamp`) "
	insertStatement += "VALUES(?, ?, ?, ?, ?, ?, ?)"
	stmtIns, err := db.Prepare(insertStatement)
	if err != nil {
		return
	}

	// Prepare statement for inserting data
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// Generate account number
	newUuid := uuid.NewV4()
	accountDetails.AccountNumber = newUuid.String()

	_, err = stmtIns.Exec(accountDetails.AccountNumber, accountDetails.BankNumber, accountDetails.AccountHolderName, accountDetails.AccountBalance, accountDetails.Overdraft, accountDetails.AvailableBalance, sqlTime)
	if err != nil {
		return
	}
	return
}

func doCreateAccountMeta(sqlTime int32, accountHolderDetails *AccountHolderDetails, accountDetails *AccountDetails) (err error) {
	db, err := sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
	if err != nil {
		return
	}
	// Create account meta
	insertStatement := "INSERT INTO accounts_meta (`accountNumber`, `bankNumber`, `accountHolderGivenName`, `accountHolderFamilyName`, `accountHolderDateOfBirth`, `accountHolderIdentificationNumber`, `accountHolderContactNumber1`, `accountHolderContactNumber2`, `accountHolderEmailAddress`, `accountHolderAddressLine1`, `accountHolderAddressLine2`, `accountHolderAddressLine3`, `accountHolderPostalCode`, `timestamp`) "
	insertStatement += "VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	stmtIns, err := db.Prepare(insertStatement)
	if err != nil {
		return
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	accountHolderDetails.AccountNumber = accountDetails.AccountNumber

	_, err = stmtIns.Exec(accountHolderDetails.AccountNumber, accountHolderDetails.BankNumber, accountHolderDetails.GivenName, accountHolderDetails.FamilyName, accountHolderDetails.DateOfBirth, accountHolderDetails.IdentificationNumber, accountHolderDetails.ContactNumber1, accountHolderDetails.ContactNumber2, accountHolderDetails.EmailAddress, accountHolderDetails.AddressLine1, accountHolderDetails.AddressLine2, accountHolderDetails.AddressLine3,
		accountHolderDetails.PostalCode, sqlTime)

	if err != nil {
		return
	}

	defer db.Close()

	return
}

func getAccountDetails(id string) (accountDetails AccountDetails) {
	db, err := sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
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
		//@TODO Allow multiple accounts
		fmt.Println("ERROR: More than one account found")
		return
	}

	return
}

func getAccountMeta(id string) (accountDetails AccountHolderDetails) {
	db, err := sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
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

func getAllAccountDetails() (allAccounts []AccountDetails) {
	db, err := sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
	if err != nil {
		fmt.Println("Could not connect to database")
		return
	}

	rows, err := db.Query("SELECT `accountNumber`, `bankNumber`, `accountHolderName` FROM `accounts`")
	if err != nil {
		fmt.Println("Error with select query: " + err.Error())
	}
	defer rows.Close()

	count := 0
	allAccounts = make([]AccountDetails, 0)

	for rows.Next() {
		accountDetailsSingle := AccountDetails{}
		if err := rows.Scan(&accountDetailsSingle.AccountNumber, &accountDetailsSingle.BankNumber, &accountDetailsSingle.AccountHolderName); err != nil {
			//@TODO Throw error
			fmt.Println("ERROR: Could not retrieve account details")
			return
		}

		allAccounts = append(allAccounts, accountDetailsSingle)
		count++
	}

	if count == 0 {
		fmt.Println("No accounts")
		return
	}

	return
}

func getSingleAccountDetail(accountNumber string) (account AccountDetails) {
	db, err := sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
	if err != nil {
		fmt.Println("Could not connect to database")
		return
	}

	rows, err := db.Query("SELECT `accountNumber`, `bankNumber`, `accountHolderName`, `accountBalance`, `overdraft`, `availableBalance` FROM `accounts` WHERE `accountNumber` = ?", accountNumber)
	if err != nil {
		fmt.Println("Error with select query: " + err.Error())
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		if err := rows.Scan(&account.AccountNumber, &account.BankNumber, &account.AccountHolderName, &account.AccountBalance, &account.Overdraft, &account.AvailableBalance); err != nil {
			//@TODO Throw error
			fmt.Println("ERROR: Could not retrieve account details")
			return
		}

		count++
	}

	if count == 0 {
		fmt.Println("No accounts")
		return
	}

	return
}

func getSingleAccountNumberByID(userID string) (accountID string) {
	db, err := sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
	if err != nil {
		fmt.Println("Could not connect to database")
		return
	}

	rows, err := db.Query("SELECT `accountNumber` FROM `accounts_meta` WHERE `accountHolderIdentificationNumber` = ?", userID)
	if err != nil {
		fmt.Println("Error with select query: " + err.Error())
	}
	defer rows.Close()

	count := 0
	// @TODO Right now this will return the latest account only, if there are two accounts
	for rows.Next() {
		if err := rows.Scan(&accountID); err != nil {
			//@TODO Throw error
			fmt.Println("ERROR: Could not retrieve account details")
			return
		}

		count++
	}

	if count == 0 {
		fmt.Println("No accounts")
		return
	}

	return
}
