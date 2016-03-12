package accounts

import (
	"reflect"
	"testing"
)

func TestProcessAccountTooFewFields(t *testing.T) {
	tst := []string{"", ""}
	_, err := ProcessAccount(tst)

	if err == nil {
		t.Errorf("ProcessAccount does not pass. Looking for %v, got %v", "Data string does not have enough fields", nil)
	}
}

func TestProcessAccountACMTTypeNotSet(t *testing.T) {
	tst := []string{"", "", ""}
	_, err := ProcessAccount(tst)

	if err == nil {
		t.Errorf("ProcessAccount does not pass. Looking for %v, got %v", "Could not get type of ACMT transaction", nil)
	}
}

func TestProcessAccountACMTTypeIncorrect(t *testing.T) {
	tst := []string{"", "", "-1000"}
	_, err := ProcessAccount(tst)

	if err == nil {
		t.Errorf("ProcessAccount does not pass. Looking for %v, got %v", "ACMT transaction code invalid", nil)
	}
}

//@TODO Implement valid ACMT tests

func TestOpenAccount(t *testing.T) {
	tst := []string{"", "", ""}
	_, err := openAccount(tst)

	if err == nil {
		t.Errorf("OpenAccount does not pass. Looking for %v, got %v", "Not all fields present", nil)
	}
}

func TestSetAccountDetails(t *testing.T) {
	tst := []string{"", "", "", "John", "Doe"}
	accountDetails, err := setAccountDetails(tst)

	if err != nil {
		t.Errorf("SetAccountDetails does not pass. ERROR. Looking for %v, got %v", nil, err)
	}

	if reflect.TypeOf(accountDetails).String() != "accounts.AccountDetails" {
		t.Errorf("SetAccountDetails does not pass. TYPE. Looking for %v, got %v", "accounts.AccountDetails", reflect.TypeOf(accountDetails).String())
	}

	if accountDetails.BankNumber != BANK_NUMBER {
		t.Errorf("SetAccountDetails does not pass. DETAILS. Looking for %v, got %v", BANK_NUMBER, accountDetails.BankNumber)
	}

	if accountDetails.Overdraft != OPENING_OVERDRAFT {
		t.Errorf("SetAccountDetails does not pass. DETAILS. Looking for %v, got %v", OPENING_OVERDRAFT, accountDetails.Overdraft)
	}

	if accountDetails.AccountBalance != OPENING_BALANCE {
		t.Errorf("SetAccountDetails does not pass. DETAILS. Looking for %v, got %v", OPENING_BALANCE, accountDetails.AccountBalance)
	}

	if accountDetails.AvailableBalance != (OPENING_BALANCE + OPENING_OVERDRAFT) {
		t.Errorf("SetAccountDetails does not pass. DETAILS. Looking for %v, got %v", (OPENING_BALANCE + OPENING_OVERDRAFT), accountDetails.AvailableBalance)
	}

	if accountDetails.AccountHolderName != "Doe,John" {
		t.Errorf("SetAccountDetails does not pass. DETAILS. Looking for %v, got %v", "Doe,John", accountDetails.AccountHolderName)
	}
}
