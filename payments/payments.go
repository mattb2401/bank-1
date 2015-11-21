package payments

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const TRANSACTION_FEE = 0.0001 // 0.01%

// @TODO Have this struct not repeat in payments and accounts
type AccountHolder struct {
	AccountNumber string
	BankNumber    string
}

func main() {
}

func CheckPayment() {
	fmt.Println("Payment Check")
}

func ProcessPAIN(data []string) (res string) {
	fmt.Println("Validating PAIN ... ")

	//There must be at least 6 elements
	if len(data) < 6 {
		fmt.Println("ERROR: Not all data is present. Run pain~help to check for needed PAIN data")
		os.Exit(1)
	}

	// Validate input
	painType, err := strconv.ParseInt(data[2], 10, 64)
	if err != nil {
		fmt.Println("Could not get type of PAIN transaction")
		log.Fatal(err)
		return
	}
	sender := parseAccountHolder(data[3])
	receiver := parseAccountHolder(data[4])
	trAmt := strings.TrimRight(data[5], "\x00")
	transactionAmount, err := strconv.ParseFloat(trAmt, 64)
	if err != nil {
		fmt.Println("ERROR: Could not convert transaction amount to float64")
		//log.Fatal(err)
		return
	}

	transaction := PAINTrans{painType, sender, receiver, transactionAmount, TRANSACTION_FEE}

	// Save transaction
	res = processPAINTransaction(transaction)

	return
}

func parseAccountHolder(account string) (accountHolder AccountHolder) {
	accountStr := strings.Split(account, "@")

	if len(accountStr) < 2 {
		fmt.Println("ERROR: Could not parse account holders")
		return
	}

	accountHolder = AccountHolder{accountStr[0], accountStr[1]}
	return
}
