package payments

import (
	"fmt"
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

*/

type PAINTrans struct {
	PainType int64
	Sender   AccountHolder
	Receiver AccountHolder
	Amount   float64
	Fee      float64
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
