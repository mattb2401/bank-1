package payments

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ksred/bank/configuration"
)

var Config configuration.Configuration

func SetConfig(config *configuration.Configuration) {
	Config = *config
}

func savePainTransaction(transaction PAINTrans) (err error) {
	db, err := sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
	if err != nil {
		return errors.New("payments.savePainTransaction: " + err.Error())
	}
	// Prepare statement for inserting data
	insertStatement := "INSERT INTO transactions (`transaction`, `type`, `senderAccountNumber`, `senderBankNumber`, `receiverAccountNumber`, `receiverBankNumber`, `transactionAmount`, `feeAmount`, `timestamp`) "
	insertStatement += "VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)"
	stmtIns, err := db.Prepare(insertStatement)
	if err != nil {
		return errors.New("payments.savePainTransaction: " + err.Error())
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	t := time.Now()
	sqlTime := int32(t.Unix())

	// The feePerc is a percentage, convert to amount
	feeAmount := transaction.Amount * transaction.Fee

	_, err = stmtIns.Exec("pain", transaction.PainType, transaction.Sender.AccountNumber, transaction.Sender.BankNumber, transaction.Receiver.AccountNumber, transaction.Receiver.BankNumber,
		transaction.Amount, feeAmount, sqlTime)

	if err != nil {
		return errors.New("payments.savePainTransaction: " + err.Error())
	}
	defer db.Close()

	return
}

//func updateAccounts(sender AccountHolder, receiver AccountHolder, transactionAmount float64, transactionFee float64) {
func updateAccounts(transaction PAINTrans) (err error) {
	t := time.Now()
	sqlTime := int32(t.Unix())

	// The feePerc is a percentage, convert to amount
	feeAmount := transaction.Amount * transaction.Fee

	switch transaction.PainType {
	// Payment
	case 1:
		err = processCreditInitiation(transaction, sqlTime, feeAmount)
		if err != nil {
			return errors.New("payments.updateAccounts: " + err.Error())
		}
		break
	// Deposit
	case 1000:
		err = processDepositInitiation(transaction, sqlTime, feeAmount)
		if err != nil {
			return errors.New("payments.updateAccounts: " + err.Error())
		}
		break
	}

	err = updateBankHoldingAccount(feeAmount, sqlTime)
	if err != nil {
		return errors.New("payments.updateAccounts: " + err.Error())
	}

	return

}

func updateBankHoldingAccount(feeAmount float64, sqlTime int32) (err error) {
	db, err := sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
	if err != nil {
		return errors.New("payments.updateBankHoldingAccount: " + err.Error())
	}

	// Add fees to bank holding account
	// Only one row in this account for now - only holds single holding bank's balance
	updateBank := "UPDATE `bank_account` SET `balance` = (`balance` + ?), `timestamp` = ?"
	stmtUpdBank, err := db.Prepare(updateBank)
	if err != nil {
		return errors.New("payments.updateBankHoldingAccount: " + err.Error())
	}
	defer stmtUpdBank.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtUpdBank.Exec(feeAmount, sqlTime)

	if err != nil {
		return errors.New("payments.updateBankHoldingAccount: " + err.Error())
	}
	defer db.Close()
	return
}

// @TODO Look at using accounts.getAccountDetails here
func checkBalance(account AccountHolder) (balance float64, err error) {
	db, err := sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
	if err != nil {
		return 0., errors.New("payments.checkBalance: " + err.Error())
	}

	rows, err := db.Query("SELECT `availableBalance` FROM `accounts` WHERE `accountNumber` = ?", account.AccountNumber)
	if err != nil {
		return 0., errors.New("payments.checkBalance: " + err.Error())
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		if err := rows.Scan(&balance); err != nil {
			return 0., errors.New("payments.checkBalance: Could not retrieve account details. " + err.Error())
		}
		count++
	}

	if count > 1 {
		return 0., errors.New("payments.checkBalance: More than one account found with uuid")
	}

	return
}

func processCreditInitiation(transaction PAINTrans, sqlTime int32, feeAmount float64) (err error) {
	db, err := sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
	if err != nil {
		return errors.New("payments.processCreditInitiation: " + err.Error())
	}

	// Only update if account local
	if transaction.Sender.BankNumber == "" {
		updateSenderStatement := "UPDATE accounts SET `accountBalance` = (`accountBalance` - ?), `availableBalance` = (`availableBalance` - ?), `timestamp` = ? WHERE `accountNumber` = ? "
		stmtUpdSender, err := db.Prepare(updateSenderStatement)
		if err != nil {
			return errors.New("payments.processCreditInitiation: " + err.Error())
		}
		defer stmtUpdSender.Close() // Close the statement when we leave main() / the program terminates

		_, err = stmtUpdSender.Exec(transaction.Amount+feeAmount, transaction.Amount+feeAmount, sqlTime, transaction.Sender.AccountNumber)

		if err != nil {
			return errors.New("payments.processCreditInitiation: " + err.Error())
		}

	} else {
		// Drop onto ledger
	}

	// Update receiver account
	// Only update if account local
	if transaction.Receiver.BankNumber == "" {
		updateStatementReceiver := "UPDATE accounts SET `accountBalance` = (`accountBalance` + ?), `availableBalance` = (`availableBalance` + ?), `timestamp` = ? WHERE `accountNumber` = ? "
		stmtUpdReceiver, err := db.Prepare(updateStatementReceiver)
		if err != nil {
			return errors.New("payments.processCreditInitiation: " + err.Error())
		}
		defer stmtUpdReceiver.Close() // Close the statement when we leave main() / the program terminates

		_, err = stmtUpdReceiver.Exec(transaction.Amount, transaction.Amount, sqlTime, transaction.Receiver.AccountNumber)

		if err != nil {
			return errors.New("payments.processCreditInitiation: " + err.Error())
		}
	} else {
		// Drop onto ledger
	}
	return
}

func processDepositInitiation(transaction PAINTrans, sqlTime int32, feeAmount float64) (err error) {
	db, err := sql.Open("mysql", Config.MySQLUser+":"+Config.MySQLPass+"@tcp("+Config.MySQLHost+":"+Config.MySQLPort+")/"+Config.MySQLDB)
	if err != nil {
		return errors.New("payments.processDepositInitiation: " + err.Error())
	}

	// We don't update sender as it is deposit
	// Update receiver account
	// The total received amount is the deposited amount minus the fee
	depositTransactionAmount := transaction.Amount - feeAmount
	// Only update if account local
	if transaction.Receiver.BankNumber == "" {
		updateStatementReceiver := "UPDATE accounts SET `accountBalance` = (`accountBalance` + ?), `availableBalance` = (`availableBalance` + ?), `timestamp` = ? WHERE `accountNumber` = ? "
		stmtUpdReceiver, err := db.Prepare(updateStatementReceiver)
		if err != nil {
			return errors.New("payments.processDepositInitiation: " + err.Error())
		}
		defer stmtUpdReceiver.Close() // Close the statement when we leave main() / the program terminates

		_, err = stmtUpdReceiver.Exec(depositTransactionAmount, depositTransactionAmount, sqlTime, transaction.Receiver.AccountNumber)

		if err != nil {
			return errors.New("payments.processDepositInitiation: " + err.Error())
		}
	} else {
		// Drop onto ledger
	}
	return
}
