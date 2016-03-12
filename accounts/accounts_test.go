package accounts

import (
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

func TestProcessAccountACMTTypeSet(t *testing.T) {
	tst := []string{"", "", "1"}
	_, err := ProcessAccount(tst)

	if err != nil {
		t.Errorf("ProcessAccount does not pass. Looking for %v, got %v", nil, err)
	}
}

func TestProcessAccountACMTTypeIncorrect(t *testing.T) {
	tst := []string{"", "", "-1000"}
	_, err := ProcessAccount(tst)

	if err == nil {
		t.Errorf("ProcessAccount does not pass. Looking for %v, got %v", "ACMT transaction code invalid", nil)
	}
}
