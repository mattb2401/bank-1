package accounts

/*
@TODO Fix DB repetition
*/

import (
	"database/sql"
	"errors"
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

func deleteAccount(accountDetails AccountDetails, accountHolderDetails AccountHolderDetails) (result bool, err error) {
	err = doDeleteAccount(&accountDetails)
	if err != nil {
		return
	}

	err = doDeleteAccountMeta(&accountHolderDetails, &accountDetails)
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

func doDeleteAccount(accountDetails *AccountDetails) (err error) {
	db, err := sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
	if err != nil {
		return
	}
	// Create account
	deleteStatement := "DELETE FROM accounts WHERE `accountNumber` = ? AND `bankNumber` = ? AND `accountHolderName` = ?) "
	stmtDel, err := db.Prepare(deleteStatement)
	if err != nil {
		return
	}

	// Prepare statement for inserting data
	defer stmtDel.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtDel.Exec(accountDetails.AccountNumber, accountDetails.BankNumber, accountDetails.AccountHolderName)
	if err != nil {
		return
	}
	// Can use db.RowsAffected() to check
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

func doDeleteAccountMeta(accountHolderDetails *AccountHolderDetails, accountDetails *AccountDetails) (err error) {
	db, err := sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
	if err != nil {
		return
	}
	// Create account meta
	deleteStatement := "DELETE FROM accounts_meta WHERE `accountNumber` = ? AND `bankNumber` = ? AND `accountHolderGivenName` = ? AND `accountHolderFamilyName` = ? AND `accountHolderDateOfBirth` = ? AND `accountHolderIdentificationNumber` = ? AND `accountHolderContactNumber1` = ? AND `accountHolderContactNumber2` = ? AND `accountHolderEmailAddress` = ? AND `accountHolderAddressLine1` = ? AND `accountHolderAddressLine2` = ? AND `accountHolderAddressLine3` = ? AND `accountHolderPostalCode` = ?) "
	deleteStatement += "VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	stmtDel, err := db.Prepare(deleteStatement)
	if err != nil {
		return
	}
	defer stmtDel.Close() // Close the statement when we leave main() / the program terminates

	accountHolderDetails.AccountNumber = accountDetails.AccountNumber

	_, err = stmtDel.Exec(accountHolderDetails.AccountNumber, accountHolderDetails.BankNumber, accountHolderDetails.GivenName, accountHolderDetails.FamilyName, accountHolderDetails.DateOfBirth, accountHolderDetails.IdentificationNumber, accountHolderDetails.ContactNumber1, accountHolderDetails.ContactNumber2, accountHolderDetails.EmailAddress, accountHolderDetails.AddressLine1, accountHolderDetails.AddressLine2, accountHolderDetails.AddressLine3,
		accountHolderDetails.PostalCode)

	if err != nil {
		return
	}

	defer db.Close()

	return
}

func getAccountDetails(id string) (accountDetails AccountDetails, err error) {
	db, err := sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
	if err != nil {
		return
	}

	rows, err := db.Query("SELECT `accountNumber`, `bankNumber`, `accountHolderName`, `accountBalance`, `overdraft`, `availableBalance` FROM `accounts` WHERE `accountNumber` = ?", id)
	if err != nil {
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		err := rows.Scan(&accountDetails.AccountNumber, &accountDetails.BankNumber, &accountDetails.AccountHolderName, &accountDetails.AccountBalance, &accountDetails.Overdraft, &accountDetails.AvailableBalance)
		if err != nil {
			break
		}
		count++
	}

	if count == 0 {
		return AccountDetails{}, errors.New("Account not found")
	}

	if count > 1 {
		//@TODO: Allow user to have multiple accounts
		return AccountDetails{}, errors.New("More than one account found")
	}

	return
}

func getAccountMeta(id string) (accountDetails AccountHolderDetails, err error) {
	db, err := sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
	if err != nil {
		return
	}

	rows, err := db.Query("SELECT `accountNumber`, `bankNumber`, `accountHolderGivenName`, `accountHolderFamilyName`, `accountHolderDateOfBirth`, `accountHolderIdentificationNumber`, `accountHolderContactNumber1`, `accountHolderContactNumber2`, `accountHolderEmailAddress`, `accountHolderAddressLine1`, `accountHolderAddressLine2`, `accountHolderAddressLine3`, `accountHolderPostalCode` FROM `accounts_meta` WHERE `accountHolderIdentificationNumber` = ?", id)
	if err != nil {
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		if err := rows.Scan(&accountDetails.AccountNumber, &accountDetails.BankNumber, &accountDetails.GivenName, &accountDetails.FamilyName, &accountDetails.DateOfBirth, &accountDetails.IdentificationNumber, &accountDetails.ContactNumber1, &accountDetails.ContactNumber2, &accountDetails.EmailAddress, &accountDetails.AddressLine1, &accountDetails.AddressLine2,
			&accountDetails.AddressLine3, &accountDetails.PostalCode); err != nil {
			//@TODO Throw error
			break
		}
		count++
	}

	if count == 0 {
		return AccountHolderDetails{}, errors.New("Account not found")
	}

	if count > 1 {
		//@TODO: Allow user to have multiple accounts
		return AccountHolderDetails{}, errors.New("More than one account found")
	}

	return
}

func getAllAccountDetails() (allAccounts []AccountDetails, err error) {
	db, err := sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
	if err != nil {
		return
	}

	rows, err := db.Query("SELECT `accountNumber`, `bankNumber`, `accountHolderName` FROM `accounts`")
	if err != nil {
		return
	}
	defer rows.Close()

	count := 0
	allAccounts = make([]AccountDetails, 0)

	for rows.Next() {
		accountDetailsSingle := AccountDetails{}
		if err := rows.Scan(&accountDetailsSingle.AccountNumber, &accountDetailsSingle.BankNumber, &accountDetailsSingle.AccountHolderName); err != nil {
			break
		}

		allAccounts = append(allAccounts, accountDetailsSingle)
		count++
	}

	return
}

func getSingleAccountDetail(accountNumber string) (account AccountDetails, err error) {
	db, err := sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
	if err != nil {
		return
	}

	rows, err := db.Query("SELECT `accountNumber`, `bankNumber`, `accountHolderName`, `accountBalance`, `overdraft`, `availableBalance` FROM `accounts` WHERE `accountNumber` = ?", accountNumber)
	if err != nil {
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		if err := rows.Scan(&account.AccountNumber, &account.BankNumber, &account.AccountHolderName, &account.AccountBalance, &account.Overdraft, &account.AvailableBalance); err != nil {
			break
		}

		count++
	}

	return
}

func getSingleAccountNumberByID(userID string) (accountID string, err error) {
	db, err := sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
	if err != nil {
		return
	}

	rows, err := db.Query("SELECT `accountNumber` FROM `accounts_meta` WHERE `accountHolderIdentificationNumber` = ?", userID)
	if err != nil {
		return
	}
	defer rows.Close()

	count := 0
	// @TODO Right now this will return the latest account only, if there are two accounts
	for rows.Next() {
		if err := rows.Scan(&accountID); err != nil {
			break
		}
		count++
	}

	return
}
