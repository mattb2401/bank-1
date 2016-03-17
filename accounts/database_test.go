package accounts

import (
	"testing"
	"time"

	"github.com/ksred/bank/configuration"
)

func TestLoadConfiguration(t *testing.T) {
	// Load app config
	_, err := configuration.LoadConfig()
	if err != nil {
		t.Errorf("loadDatabase does not pass. Configuration does not load, looking for %v, got %v", nil, err)
	}
}

func TestLoadDatabase(t *testing.T) {
	// Load app config
	Config, err := configuration.LoadConfig()
	if err != nil {
		t.Errorf("loadDatabase does not pass. Configuration does not load, looking for %v, got %v", nil, err)
	}
	// Set config in packages
	SetConfig(&Config)

	_, err = loadDatabase()

	if err != nil {
		t.Errorf("LoadDatabase does not pass. Looking for %v, got %v", nil, err)
	}
}

func TestDoCreateAccount(t *testing.T) {
	//accountDetails AccountDetails, accountHolderDetails AccountHolderDetails
	accountDetail := AccountDetails{
		"",
		"",
		"User,Test",
		0.,
		0.,
		0.,
		0,
	}

	ti := time.Now()
	sqlTime := int32(ti.Unix())
	err := doCreateAccount(sqlTime, &accountDetail)

	if err != nil {
		t.Errorf("DoCreateAccount does not pass. Looking for %v, got %v", nil, err)
	}

	err = doDeleteAccount(&accountDetail)
	if err != nil {
		t.Errorf("DoDeleteAccount does not pass. Looking for %v, got %v", nil, err)
	}
}

func TestDoAccountMeta(t *testing.T) {
	accountDetail := AccountDetails{
		"",
		"",
		"User,Test",
		0.,
		0.,
		0.,
		0,
	}

	accountHolderDetail := AccountHolderDetails{
		"",
		"",
		"Test",
		"User",
		"1900-01-01",
		"19000101-1000-100",
		"555-123-1234",
		"",
		"test@user.com",
		"Address 1",
		"Address 2",
		"Address 3",
		"22202",
	}

	ti := time.Now()
	sqlTime := int32(ti.Unix())
	err := doCreateAccountMeta(sqlTime, &accountHolderDetail, &accountDetail)

	if err != nil {
		t.Errorf("DoCreateAccountMeta does not pass. Looking for %v, got %v", nil, err)
	}

	err = doDeleteAccountMeta(&accountHolderDetail, &accountDetail)
	if err != nil {
		t.Errorf("DoDeleteAccountMeta does not pass. Looking for %v, got %v", nil, err)
	}
}

func TestGetAccount(t *testing.T) {
	accountDetail := AccountDetails{
		"",
		"",
		"User,Test",
		0.,
		0.,
		0.,
		0,
	}

	accountHolderDetail := AccountHolderDetails{
		"",
		"",
		"Test",
		"User",
		"1900-01-01",
		"19000101-1000-100",
		"555-123-1234",
		"",
		"test@user.com",
		"Address 1",
		"Address 2",
		"Address 3",
		"22202",
	}

	ti := time.Now()
	sqlTime := int32(ti.Unix())

	err := doCreateAccount(sqlTime, &accountDetail)
	if err != nil {
		t.Errorf("GetAccountCreateAccount does not pass. Looking for %v, got %v", nil, err)
	}

	err = doCreateAccountMeta(sqlTime, &accountHolderDetail, &accountDetail)
	if err != nil {
		t.Errorf("GetAccountCreateAccountMeta does not pass. Looking for %v, got %v", nil, err)
	}

	// Get account
	getAccountDetails, err := getAccountDetails(accountDetail.AccountNumber)
	if err != nil {
		t.Errorf("GetAccountDetails does not pass. Looking for %v, got %v", nil, err)
		return
	}
	//Check values
	if getAccountDetails.AccountNumber != accountDetail.AccountNumber {
		t.Errorf("GetAccountDetails does not pass. DETAILS. AccountNumber: Looking for %v, got %v", accountDetail.AccountNumber, getAccountDetails.AccountNumber)
	}
	if getAccountDetails.BankNumber != "" {
		t.Errorf("GetAccountDetails does not pass. DETAILS. BankNumber: Looking for %v, got %v", "", getAccountDetails.BankNumber)
	}
	if getAccountDetails.Overdraft != 0. {
		t.Errorf("GetAccountDetails does not pass. DETAILS. Overdraft: Looking for %v, got %v", 0., getAccountDetails.Overdraft)
	}
	if getAccountDetails.AvailableBalance != 0. {
		t.Errorf("GetAccountDetails does not pass. DETAILS. AvailableBalance: Looking for %v, got %v", 0., getAccountDetails.AvailableBalance)
	}
	if getAccountDetails.AccountBalance != 0. {
		t.Errorf("GetAccountDetails does not pass. DETAILS. AccountBalance: Looking for %v, got %v", 0., getAccountDetails.AccountBalance)
	}
	if getAccountDetails.AccountHolderName != "User,Test" {
		t.Errorf("GetAccountDetails does not pass. DETAILS. AccountHodlerName: Looking for %v, got %v", "User,Test", getAccountDetails.AccountHolderName)
	}

	err = doDeleteAccount(&accountDetail)
	if err != nil {
		t.Errorf("GetAccountDeleteAccount does not pass. Looking for %v, got %v", nil, err)
	}

	err = doDeleteAccountMeta(&accountHolderDetail, &accountDetail)
	if err != nil {
		t.Errorf("GetAccountDeleteAccountMeta does not pass. Looking for %v, got %v", nil, err)
	}

}
