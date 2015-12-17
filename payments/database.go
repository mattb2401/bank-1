package payments

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ksred/bank/configuration"
	"time"
)

var Config configuration.Configuration

func SetConfig(config *configuration.Configuration) {
	Config = *config
}

func loadDatabase() (db *sql.DB) {
	db, err := sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
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

func savePainTransaction(transaction PAINTrans) {
	db, err := sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
	if err != nil {
		fmt.Println("Could not connect to database")
		return
	}
	// Prepare statement for inserting data
	insertStatement := "INSERT INTO transactions (`transaction`, `type`, `senderAccountNumber`, `senderBankNumber`, `receiverAccountNumber`, `receiverBankNumber`, `transactionAmount`, `feeAmount`, `timestamp`) "
	insertStatement += "VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)"
	stmtIns, err := db.Prepare(insertStatement)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	t := time.Now()
	sqlTime := int32(t.Unix())

	// The feePerc is a percentage, convert to amount
	feeAmount := transaction.Amount * transaction.Fee

	_, err = stmtIns.Exec("pain", transaction.PainType, transaction.Sender.AccountNumber, transaction.Sender.BankNumber, transaction.Receiver.AccountNumber, transaction.Receiver.BankNumber,
		transaction.Amount, feeAmount, sqlTime)

	if err != nil {
		fmt.Println("Could not save results: " + err.Error())
	}
	defer db.Close()
}

//func updateAccounts(sender AccountHolder, receiver AccountHolder, transactionAmount float64, transactionFee float64) {
func updateAccounts(transaction PAINTrans) (result string) {
	t := time.Now()
	sqlTime := int32(t.Unix())

	// Update sender account
	fmt.Println("Processing transaction...")
	fmt.Println(transaction)

	// The feePerc is a percentage, convert to amount
	feeAmount := transaction.Amount * transaction.Fee

	switch transaction.PainType {
	// Payment
	case 1:
		result = processCreditInitiation(transaction, sqlTime, feeAmount)
		break
	// Deposit
	case 1000:
		result = processDepositInitiation(transaction, sqlTime, feeAmount)
		break
	}

	resBankUpdate := updateBankHoldingAccount(feeAmount, sqlTime)
	if !resBankUpdate {
		result = "0~Could not update bank account"
		return
	}

	return

}

func updateBankHoldingAccount(feeAmount float64, sqlTime int32) (result bool) {
	db, err := sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
	if err != nil {
		fmt.Println("Could not connect to database")
		return
	}

	// Add fees to bank holding account
	fmt.Println("   Processing bank...")
	// Only one row in this account for now - only holds single holding bank's balance
	updateBank := "UPDATE `bank_account` SET `balance` = (`balance` + ?), `timestamp` = ?"
	stmtUpdBank, err := db.Prepare(updateBank)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtUpdBank.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtUpdBank.Exec(feeAmount, sqlTime)

	if err != nil {
		fmt.Println("Could not save results: " + err.Error())
	}
	defer db.Close()
	return
}

// @TODO Look at using accounts.getAccountDetails here
func checkBalance(account AccountHolder) (balance float64) {
	db, err := sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
	if err != nil {
		fmt.Println("Could not connect to database")
		return
	}

	rows, err := db.Query("SELECT `availableBalance` FROM `accounts` WHERE `accountNumber` = ?", account.AccountNumber)
	if err != nil {
		fmt.Println("Error with select query: " + err.Error())
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		if err := rows.Scan(&balance); err != nil {
			//@TODO Throw error
			fmt.Println("ERROR: Could not retrieve account details")
			return
		}
		count++
	}

	if count > 1 {
		fmt.Println("ERROR: More than one account found with uuid")
	}

	return
}

func processCreditInitiation(transaction PAINTrans, sqlTime int32, feeAmount float64) (result string) {
	db, err := sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
	if err != nil {
		fmt.Println("Could not connect to database")
		return
	}

	// Only update if account local
	if transaction.Sender.BankNumber == "" {
		fmt.Println("   Processing sender...")
		updateSenderStatement := "UPDATE accounts SET `accountBalance` = (`accountBalance` - ?), `availableBalance` = (`availableBalance` - ?), `timestamp` = ? WHERE `accountNumber` = ? "
		stmtUpdSender, err := db.Prepare(updateSenderStatement)
		if err != nil {
			//panic(err.Error()) // proper error handling instead of panic in your app
			result = "0~Could not save payment"
			return
		}
		defer stmtUpdSender.Close() // Close the statement when we leave main() / the program terminates

		resUpd, err := stmtUpdSender.Exec(transaction.Amount+feeAmount, transaction.Amount+feeAmount, sqlTime, transaction.Sender.AccountNumber)
		fmt.Println(resUpd)

		if err != nil {
			fmt.Println("Could not save results: " + err.Error())
			result = "0~Could not process payment"
			return
		}

		result = "1~Payment processed"
	} else {
		// Drop onto ledger
	}

	// Update receiver account
	// Only update if account local
	if transaction.Receiver.BankNumber == "" {
		fmt.Println("   Processing receiver...")
		updateStatementReceiver := "UPDATE accounts SET `accountBalance` = (`accountBalance` + ?), `availableBalance` = (`availableBalance` + ?), `timestamp` = ? WHERE `accountNumber` = ? "
		stmtUpdReceiver, err := db.Prepare(updateStatementReceiver)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		defer stmtUpdReceiver.Close() // Close the statement when we leave main() / the program terminates

		_, err = stmtUpdReceiver.Exec(transaction.Amount, transaction.Amount, sqlTime, transaction.Receiver.AccountNumber)

		if err != nil {
			fmt.Println("Could not save results: " + err.Error())
		}
	} else {
		// Drop onto ledger
	}
	return
}

func processDepositInitiation(transaction PAINTrans, sqlTime int32, feeAmount float64) (result string) {
	db, err := sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
	if err != nil {
		fmt.Println("Could not connect to database")
		return
	}

	// We don't update sender as it is deposit
	// Update receiver account
	// The total received amount is the deposited amount minus the fee
	depositTransactionAmount := transaction.Amount - feeAmount
	// Only update if account local
	if transaction.Receiver.BankNumber == "" {
		fmt.Println("   Processing receiver...")
		updateStatementReceiver := "UPDATE accounts SET `accountBalance` = (`accountBalance` + ?), `availableBalance` = (`availableBalance` + ?), `timestamp` = ? WHERE `accountNumber` = ? "
		stmtUpdReceiver, err := db.Prepare(updateStatementReceiver)
		if err != nil {
			//panic(err.Error()) // proper error handling instead of panic in your app
			result = "0~Could not save deposit"
			return
		}
		defer stmtUpdReceiver.Close() // Close the statement when we leave main() / the program terminates

		_, err = stmtUpdReceiver.Exec(transaction.Amount, depositTransactionAmount, sqlTime, transaction.Receiver.AccountNumber)

		if err != nil {
			fmt.Println("Could not save results: " + err.Error())

			result = "0~Deposit not successful"
			return
		}

		result = "1~Successful deposit"
	} else {
		// Drop onto ledger
	}
	return
}
