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
	painType int64
	sender   AccountHolder
	receiver AccountHolder
	amount   float64
}

func processPAINTransaction(transaction PAINTrans, TRANSACTION_FEE float64) (res string) {
	fmt.Printf("Process transaction %v", transaction)

	// Checks for transaction (avail balance, accounts open, etc)
	balanceAvailable := checkBalance(transaction.sender)
	if balanceAvailable < transaction.amount {
		fmt.Println("ERROR: Insufficient funds available")
		res = "0~Insufficient funds"
		return
	}
	// Save in transaction table
	savePainTransaction(transaction, TRANSACTION_FEE)
	// Amend sender and receiver accounts
	// Amend bank's account with fee addition

	res = "true"
	return

}
