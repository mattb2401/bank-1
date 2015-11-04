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

func processPAINTransaction(transaction PAINTrans) (res string) {
	fmt.Printf("Process transaction %v", transaction)

	// Checks for transaction (avail balance, accounts open, etc)
	// Save in transaction table
	// Amend sender and receiver accounts
	// Amend bank's account with fee addition

	res = "true"
	return

}
