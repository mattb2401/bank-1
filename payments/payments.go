package payments

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PAINTrans struct {
	sender   AccountHolder
	receiver AccountHolder
	amount   float64
}

type AccountHolder struct {
	accountNumber int64
	bankNumber    int64
}

func main() {
}

func CheckPayment() {
	fmt.Println("Payment Check")
}

func ProcessPAIN(data []string) {
	fmt.Println("Validating PAIN ... ")

	//There must be at least 4 elements
	if len(data) < 4 {
		fmt.Println("ERROR: Not all data is present. Run pain~help to check for needed PAIN data")
		os.Exit(1)
	}

	// Validate input
	sender := parseAccountHolder(data[1])
	receiver := parseAccountHolder(data[2])
	transactionAmount, err := strconv.ParseFloat(data[3], 64)
	if err != nil {
		fmt.Println("ERROR: Could not convert transaction amount to float64")
		os.Exit(1)
	}

	transaction := PAINTrans{sender, receiver, transactionAmount}

	res := false
	// Save transaction
	res = savePAINTransaction(transaction)
	// Send transaction
	res = sendPAINTransaction(transaction)
	// Notify
	if res {
		fmt.Println("1")
	}
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

func savePAINTransaction(transaction PAINTrans) (res bool) {
	// @TODO Return transaction guid
	fmt.Printf("Save transaction %v", transaction)

	res = true
	return
}

func sendPAINTransaction(transaction PAINTrans) (res bool) {
	fmt.Printf("Send transaction %v", transaction)
	// Check if sender is local
	// Check if recipient is local
	// Adjust balances if local
	// Notify

	res = true
	return

}
