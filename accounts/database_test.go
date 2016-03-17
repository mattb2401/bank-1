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
		"1111",
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
}

func TestDoDeleteAccount(t *testing.T) {
	accountDetail := AccountDetails{
		"1111",
		"",
		"User,Test",
		0.,
		0.,
		0.,
		0,
	}

	err := doDeleteAccount(&accountDetail)

	if err != nil {
		t.Errorf("DoDeleteAccount does not pass. Looking for %v, got %v", nil, err)
	}
}

func TestDoAccountMeta(t *testing.T) {
	accountDetail := AccountDetails{
		"1111",
		"",
		"User,Test",
		0.,
		0.,
		0.,
		0,
	}

	accountHolderDetail := AccountHolderDetails{
		"1111",
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
}

func TestDoDeleteAccountMeta(t *testing.T) {
	accountDetail := AccountDetails{
		"1111",
		"",
		"User,Test",
		0.,
		0.,
		0.,
		0,
	}

	accountHolderDetail := AccountHolderDetails{
		"1111",
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

	err := doDeleteAccountMeta(&accountHolderDetail, &accountDetail)

	if err != nil {
		t.Errorf("DoDeleteAccountMeta does not pass. Looking for %v, got %v", nil, err)
	}
}
