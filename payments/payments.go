package payments

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type AccountHolder struct {
	accountNumber int64
	bankNumber    int64
}

const TRANSACTION_FEE = 0.0001 // 0.01%

func main() {
}

func CheckPayment() {
	fmt.Println("Payment Check")
}

func ProcessPAIN(data []string) (res string) {
	fmt.Println("Validating PAIN ... ")

	//There must be at least 5 elements
	if len(data) < 5 {
		fmt.Println("ERROR: Not all data is present. Run pain~help to check for needed PAIN data")
		os.Exit(1)
	}

	// Validate input
	painType, err := strconv.ParseInt(data[1], 10, 64)
	if err != nil {
		fmt.Println("Could not get type of PAIN transaction")
		log.Fatal(err)
		return
	}
	sender := parseAccountHolder(data[2])
	receiver := parseAccountHolder(data[3])
	trAmt := strings.TrimRight(data[4], "\x00")
	transactionAmount, err := strconv.ParseFloat(trAmt, 64)
	if err != nil {
		fmt.Println("ERROR: Could not convert transaction amount to float64")
		//log.Fatal(err)
		return
	}

	transaction := PAINTrans{painType, sender, receiver, transactionAmount}

	// Save transaction
	res = processPAINTransaction(transaction, TRANSACTION_FEE)

	return
}

func parseAccountHolder(account string) (accountHolder AccountHolder) {
	accountStr := strings.Split(account, "@")
	accountAccNum, err := strconv.ParseInt(accountStr[0], 10, 64)
	if err != nil {
		fmt.Println("ERROR: Could not convert account details S1")
		os.Exit(1)
	}
	if len(accountStr[1]) == 0 {
		accountStr[1] = "0"
	}
	accountBankNum, err := strconv.ParseInt(accountStr[1], 10, 64)
	if err != nil {
		fmt.Println("ERROR: Could not convert account details S2")
		os.Exit(1)
	}

	accountHolder = AccountHolder{accountAccNum, accountBankNum}
	return
}
