package payments

import (
	"fmt"
	"log"
	"os"
	"strconv"
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
