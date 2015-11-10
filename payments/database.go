package payments

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
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

func checkBalance(account AccountHolder) (balance int) {
	return
}

func savePainTransaction(transaction PAINTrans, feePerc float64) {
	configuration := Configuration{}
	loadConfig(&configuration)

	db, err := sql.Open("mysql", configuration.MySQLUser+":"+configuration.MySQLPass+"@tcp("+configuration.MySQLHost+":"+configuration.MySQLPort+")/"+configuration.MySQLDB)
	if err != nil {
		fmt.Println("Could not connect to database")
		return
	}
	// Prepare statement for inserting data
	insertStatement := "INSERT INTO transactions (`transaction`, `type`, `senderAccountNum`, `senderBankNum`, `receiverAccountNum`, `receiverBankNum`, `transactionAmount`, `feeAmount`, `timestamp`) "
	insertStatement += "VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)"
	stmtIns, err := db.Prepare(insertStatement)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// Convert variables
	//sqlChange, _ := strconv.ParseFloat(strings.Replace(stock.Change, ",", "", -1), 64)
	t := time.Now()
	sqlTime := int32(t.Unix())

	// The feePerc is a percentage, convert to amount
	feeAmount := transaction.amount * feePerc

	_, err = stmtIns.Exec("pain", transaction.painType, transaction.sender.accountNumber, transaction.sender.bankNumber, transaction.receiver.accountNumber, transaction.receiver.bankNumber,
		transaction.amount, feeAmount, sqlTime)

	if err != nil {
		fmt.Println("Could not save results: " + err.Error())
	}
	defer db.Close()
}

func updateAccounts(sender AccountHolder, receiver AccountHolder, transactionFee int) {
	// Update sender account
	// Update receiver account
	// Add fees to bank holding account
}
