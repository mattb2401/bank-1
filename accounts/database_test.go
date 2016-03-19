package accounts

import (
	"reflect"
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
		t.Errorf("DoAccountMeta CreateAccount does not pass. Looking for %v, got %v", nil, err)
	}

	err = doDeleteAccountMeta(&accountHolderDetail, &accountDetail)
	if err != nil {
		t.Errorf("DoAccountMeta DeleteAccount does not pass. Looking for %v, got %v", nil, err)
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
		t.Errorf("GetAccount CreateAccount does not pass. Looking for %v, got %v", nil, err)
	}

	err = doCreateAccountMeta(sqlTime, &accountHolderDetail, &accountDetail)
	if err != nil {
		t.Errorf("GetAccount CreateAccountMeta does not pass. Looking for %v, got %v", nil, err)
	}

	// Get account
	getAccountDetails, err := getAccountDetails(accountDetail.AccountNumber)
	if err != nil {
		t.Errorf("GetAccount does not pass. Looking for %v, got %v", nil, err)
		return
	}
	//Check values
	if getAccountDetails.AccountNumber != accountDetail.AccountNumber {
		t.Errorf("GetAccount does not pass. DETAILS. AccountNumber: Looking for %v, got %v", accountDetail.AccountNumber, getAccountDetails.AccountNumber)
	}
	if getAccountDetails.BankNumber != "" {
		t.Errorf("GetAccount does not pass. DETAILS. BankNumber: Looking for %v, got %v", "", getAccountDetails.BankNumber)
	}
	if getAccountDetails.Overdraft != 0. {
		t.Errorf("GetAccount does not pass. DETAILS. Overdraft: Looking for %v, got %v", 0., getAccountDetails.Overdraft)
	}
	if getAccountDetails.AvailableBalance != 0. {
		t.Errorf("GetAccount does not pass. DETAILS. AvailableBalance: Looking for %v, got %v", 0., getAccountDetails.AvailableBalance)
	}
	if getAccountDetails.AccountBalance != 0. {
		t.Errorf("GetAccount does not pass. DETAILS. AccountBalance: Looking for %v, got %v", 0., getAccountDetails.AccountBalance)
	}
	if getAccountDetails.AccountHolderName != "User,Test" {
		t.Errorf("GetAccount does not pass. DETAILS. AccountHodlerName: Looking for %v, got %v", "User,Test", getAccountDetails.AccountHolderName)
	}

	err = doDeleteAccount(&accountDetail)
	if err != nil {
		t.Errorf("GetAccount DeleteAccount does not pass. Looking for %v, got %v", nil, err)
	}

	err = doDeleteAccountMeta(&accountHolderDetail, &accountDetail)
	if err != nil {
		t.Errorf("GetAccount DeleteAccountMeta does not pass. Looking for %v, got %v", nil, err)
	}

}

func TestGetAccountMeta(t *testing.T) {
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
		t.Errorf("GetAccountMeta CreateAccount does not pass. Looking for %v, got %v", nil, err)
	}

	err = doCreateAccountMeta(sqlTime, &accountHolderDetail, &accountDetail)
	if err != nil {
		t.Errorf("GetAccountMeta CreateAccountMeta does not pass. Looking for %v, got %v", nil, err)
	}

	// Get account
	getAccountDetails, err := getAccountMeta(accountHolderDetail.IdentificationNumber)
	if err != nil {
		t.Errorf("GetAccountMeta does not pass. Looking for %v, got %v", nil, err)
		return
	}
	//Check values
	if getAccountDetails.AccountNumber != accountDetail.AccountNumber {
		t.Errorf("GetAccountMeta does not pass. DETAILS. AccountNumber: Looking for %v, got %v", accountDetail.AccountNumber, getAccountDetails.AccountNumber)
	}
	if getAccountDetails.BankNumber != "" {
		t.Errorf("GetAccountMeta does not pass. DETAILS. BankNumber: Looking for %v, got %v", "", getAccountDetails.BankNumber)
	}
	if getAccountDetails.GivenName != "Test" {
		t.Errorf("GetAccountMeta does not pass. DETAILS. GivenName: Looking for %v, got %v", "Test", getAccountDetails.GivenName)
	}
	if getAccountDetails.FamilyName != "User" {
		t.Errorf("GetAccountMeta does not pass. DETAILS. FamilyName: Looking for %v, got %v", "User", getAccountDetails.FamilyName)
	}
	if getAccountDetails.PostalCode != "22202" {
		t.Errorf("GetAccountMeta does not pass. DETAILS. PostalCode: Looking for %v, got %v", "22202", getAccountDetails.PostalCode)
	}
	if getAccountDetails.IdentificationNumber != "19000101-1000-100" {
		t.Errorf("GetAccountMeta does not pass. DETAILS. IdentificationNumber: Looking for %v, got %v", "19000101-1000-100", getAccountDetails.IdentificationNumber)
	}
	if getAccountDetails.DateOfBirth != "1900-01-01" {
		t.Errorf("GetAccountMeta does not pass. DETAILS. DateOfBirth: Looking for %v, got %v", "1900-01-01", getAccountDetails.DateOfBirth)
	}
	if getAccountDetails.ContactNumber1 != "555-123-1234" {
		t.Errorf("GetAccountMeta does not pass. DETAILS. ContactNumber1: Looking for %v, got %v", "555-123-1234", getAccountDetails.ContactNumber1)
	}
	if getAccountDetails.ContactNumber2 != "" {
		t.Errorf("GetAccountMeta does not pass. DETAILS. ContactNumber2: Looking for %v, got %v", "", getAccountDetails.ContactNumber2)
	}
	if getAccountDetails.AddressLine1 != "Address 1" {
		t.Errorf("GetAccountMeta does not pass. DETAILS. AddressLine1: Looking for %v, got %v", "Address 1", getAccountDetails.AddressLine1)
	}
	if getAccountDetails.AddressLine2 != "Address 2" {
		t.Errorf("GetAccountMeta does not pass. DETAILS. AddressLine2: Looking for %v, got %v", "Address 2", getAccountDetails.AddressLine2)
	}
	if getAccountDetails.AddressLine3 != "Address 3" {
		t.Errorf("GetAccountMeta does not pass. DETAILS. AddressLine3: Looking for %v, got %v", "Address 3", getAccountDetails.AddressLine3)
	}
	if getAccountDetails.EmailAddress != "test@user.com" {
		t.Errorf("GetAccountMeta does not pass. DETAILS. EmailAddress: Looking for %v, got %v", "test@user.com", getAccountDetails.EmailAddress)
	}

	err = doDeleteAccount(&accountDetail)
	if err != nil {
		t.Errorf("GetAccountMeta DeleteAccount does not pass. Looking for %v, got %v", nil, err)
	}

	err = doDeleteAccountMeta(&accountHolderDetail, &accountDetail)
	if err != nil {
		t.Errorf("GetAccountMeta DeleteAccountMeta does not pass. Looking for %v, got %v", nil, err)
	}

}

func TestGetAllAccountDetails(t *testing.T) {
	accounts, err := getAllAccountDetails()

	if err != nil {
		t.Errorf("GetAllAccountDetails does not pass. Looking for %v, got %v", nil, err)
	}

	if reflect.TypeOf(accounts).String() != "[]accounts.AccountDetails" {
		t.Errorf("GetAllAccountDetails does not pass. TYPE. Looking for %v, got %v", "[]accounts.AccountDetails", reflect.TypeOf(accounts).String())
	}
}

func TestGetSingleAccountDetail(t *testing.T) {
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
		t.Errorf("GetSingleAccountDetail CreateAccount does not pass. Looking for %v, got %v", nil, err)
	}

	err = doCreateAccountMeta(sqlTime, &accountHolderDetail, &accountDetail)
	if err != nil {
		t.Errorf("GetSingleAccountDetail CreateAccountMeta does not pass. Looking for %v, got %v", nil, err)
	}

	// Do get account call
	getAccountDetails, err := getSingleAccountDetail(accountDetail.AccountNumber)
	if err != nil {
		t.Errorf("GetSingleAccountDetail does not pass. Looking for %v, got %v", nil, err)
	}

	//Check values
	if getAccountDetails.AccountNumber != accountDetail.AccountNumber {
		t.Errorf("GetSingleAccountDetail does not pass. DETAILS. AccountNumber: Looking for %v, got %v", accountDetail.AccountNumber, getAccountDetails.AccountNumber)
	}
	if getAccountDetails.BankNumber != "" {
		t.Errorf("GetSingleAccountDetail does not pass. DETAILS. BankNumber: Looking for %v, got %v", "", getAccountDetails.BankNumber)
	}
	if getAccountDetails.Overdraft != 0. {
		t.Errorf("GetSingleAccountDetail does not pass. DETAILS. Overdraft: Looking for %v, got %v", 0., getAccountDetails.Overdraft)
	}
	if getAccountDetails.AvailableBalance != 0. {
		t.Errorf("GetSingleAccountDetail does not pass. DETAILS. AvailableBalance: Looking for %v, got %v", 0., getAccountDetails.AvailableBalance)
	}
	if getAccountDetails.AccountBalance != 0. {
		t.Errorf("GetSingleAccountDetail does not pass. DETAILS. AccountBalance: Looking for %v, got %v", 0., getAccountDetails.AccountBalance)
	}
	if getAccountDetails.AccountHolderName != "User,Test" {
		t.Errorf("GetSingleAccountDetail does not pass. DETAILS. AccountHodlerName: Looking for %v, got %v", "User,Test", getAccountDetails.AccountHolderName)
	}

	err = doDeleteAccount(&accountDetail)
	if err != nil {
		t.Errorf("GetSingleAccountDetail DeleteAccount does not pass. Looking for %v, got %v", nil, err)
	}

	err = doDeleteAccountMeta(&accountHolderDetail, &accountDetail)
	if err != nil {
		t.Errorf("GetSingleAccountDetail DeleteAccountMeta does not pass. Looking for %v, got %v", nil, err)
	}
}

func TestGetSingleAccountNumberByID(t *testing.T) {
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
		t.Errorf("GetSingleAccountNumberByID CreateAccount does not pass. Looking for %v, got %v", nil, err)
	}

	err = doCreateAccountMeta(sqlTime, &accountHolderDetail, &accountDetail)
	if err != nil {
		t.Errorf("GetSingleAccountNumberByID CreateAccountMeta does not pass. Looking for %v, got %v", nil, err)
	}

	// Do get account call
	getAccountNumber, err := getSingleAccountNumberByID(accountHolderDetail.IdentificationNumber)
	if err != nil {
		t.Errorf("GetSingleAccountNumberByID does not pass. Looking for %v, got %v", nil, err)
	}

	//Check values
	if getAccountNumber != accountDetail.AccountNumber {
		t.Errorf("GetSingleAccountNumberByID does not pass. DETAILS. AccountNumber: Looking for %v, got %v", getAccountNumber, accountDetail.AccountNumber)
	}

	err = doDeleteAccount(&accountDetail)
	if err != nil {
		t.Errorf("GetSingleAccountNumberByID DeleteAccount does not pass. Looking for %v, got %v", nil, err)
	}

	err = doDeleteAccountMeta(&accountHolderDetail, &accountDetail)
	if err != nil {
		t.Errorf("GetSingleAccountNumberByID DeleteAccountMeta does not pass. Looking for %v, got %v", nil, err)
	}
}
