package payments

import (
	"fmt"
	"github.com/ksred/bank/appauth"
	"strconv"
	"strings"
)

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

type PAINTrans struct {
	PainType int64
	Sender   AccountHolder
	Receiver AccountHolder
	Amount   float64
	Fee      float64
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

	// Save transaction
	result = processPAINTransaction(transaction)

	return
}

func processPAINTransaction(transaction PAINTrans) (res string) {
	fmt.Printf("Process transaction %v", transaction)

	// Checks for transaction (avail balance, accounts open, etc)
	balanceAvailable := checkBalance(transaction.Sender)
	if balanceAvailable < transaction.Amount {
		fmt.Println("ERROR: Insufficient funds available")
		res = "0~Insufficient funds"
		return
	}
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
	if tokenUser != sender.AccountNumber {
		result = "0~Sender not valid"
		return
	}

	// Issue deposit
	// @TODO This flow show be fixed. Maybe have banks approve deposits before initiation, or
	// immediate approval below a certain amount subject to rate limiting
	// @TODO Make sure fees are deducted off deposit amount
	transaction := PAINTrans{painType, sender, receiver, transactionAmount, TRANSACTION_FEE}
	// Save transaction
	result = processPAINTransaction(transaction)

	return
}
