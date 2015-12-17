package payments

/*
PAIN transactions are as follows

Payments initiation:
1 - CustomerCreditTransferInitiationV06
2 - CustomerPaymentStatusReportV06
7 - CustomerPaymentReversalV05
8 - CustomerDirectDebitInitiationV05

Payments mandates:
9 - MandateInitiationRequestV04
10 - MandateAmendmentRequestV04
11 - MandateCancellationRequestV04
12 - MandateAcceptanceReportV04

#### Custom payments
1000 - CustomerDepositInitiation (@FIXME Will need to implement this properly, for now we use it to demonstrate functionality)

*/

import (
	"fmt"
	"github.com/ksred/bank/appauth"
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

type PAINTrans struct {
	PainType int64
	Sender   AccountHolder
	Receiver AccountHolder
	Amount   float64
	Fee      float64
}

func main() {
}

func CheckPayment() {
	fmt.Println("Payment Check")
}

func ProcessPAIN(data []string) (result string) {
	fmt.Println("Validating PAIN ... ")

	//There must be at least 3 elements
	if len(data) < 3 {
		fmt.Println("ERROR: Not all data is present. Run pain~help to check for needed PAIN data")
		os.Exit(1)
	}

	// Get type
	painType, err := strconv.ParseInt(data[2], 10, 64)
	if err != nil {
		fmt.Println("Could not get type of PAIN transaction")
		log.Fatal(err)
		return
	}

	switch painType {
	case 1:
		//There must be at least 6 elements
		if len(data) < 6 {
			fmt.Println("ERROR: Not all data is present. Run pain~help to check for needed PAIN data")
			os.Exit(1)
		}

		result = painCreditTransferInitiation(painType, data)
		break
	case 1000:
		//There must be at least 4 elements
		//token~pain~type~amount
		if len(data) < 4 {
			fmt.Println("ERROR: Not all data is present. Run pain~help to check for needed PAIN data")
			os.Exit(1)
		}
		result = customerDepositInitiation(painType, data)
		break
	}

	return
}

func painCreditTransferInitiation(painType int64, data []string) (result string) {

	// Validate input
	sender := parseAccountHolder(data[3])
	receiver := parseAccountHolder(data[4])
	trAmt := strings.TrimRight(data[5], "\x00")
	transactionAmount, err := strconv.ParseFloat(trAmt, 64)
	if err != nil {
		fmt.Println("ERROR: Could not convert transaction amount to float64")
		//log.Fatal(err)
		return
	}

	// Check if sender valid
	tokenUser := appauth.GetUserFromToken(data[0])
	if tokenUser != sender.AccountNumber {
		result = "0~Sender not valid"
		return
	}

	transaction := PAINTrans{painType, sender, receiver, transactionAmount, TRANSACTION_FEE}

	// Checks for transaction (avail balance, accounts open, etc)
	balanceAvailable := checkBalance(transaction.Sender)
	if balanceAvailable < transaction.Amount {
		fmt.Println("ERROR: Insufficient funds available")
		result = "0~Insufficient funds"
		return
	}

	// Save transaction
	result = processPAINTransaction(transaction)

	return
}

func processPAINTransaction(transaction PAINTrans) (res string) {
	fmt.Printf("Process transaction %v", transaction)
	// Test: pain~1~1b2ca241-0373-4610-abad-da7b06c50a7b@~181ac0ae-45cb-461d-b740-15ce33e4612f@~20
	// Save in transaction table
	savePainTransaction(transaction)
	// Amend sender and receiver accounts
	// Amend bank's account with fee addition
	updateAccounts(transaction)

	res = "true"
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

func customerDepositInitiation(painType int64, data []string) (result string) {
	// Validate input
	// Sender is bank
	sender := parseAccountHolder("0@0")
	receiver := parseAccountHolder(data[3])
	trAmt := strings.TrimRight(data[4], "\x00")
	transactionAmount, err := strconv.ParseFloat(trAmt, 64)
	if err != nil {
		fmt.Println("ERROR: Could not convert transaction amount to float64")
		//log.Fatal(err)
		return
	}

	// Check if sender valid
	tokenUser := appauth.GetUserFromToken(data[0])
	if tokenUser != receiver.AccountNumber {
		result = "0~Sender not valid"
		return
	}

	// Issue deposit
	// @TODO This flow show be fixed. Maybe have banks approve deposits before initiation, or
	// immediate approval below a certain amount subject to rate limiting
	transaction := PAINTrans{painType, sender, receiver, transactionAmount, TRANSACTION_FEE}
	// Save transaction
	result = processPAINTransaction(transaction)

	return
}
