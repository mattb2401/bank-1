package accounts

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ksred/bank/appauth"
)

/*
Accounts package to deal with all account related queries.

@TODO Implement the ISO20022 standard
http://www.iso20022.org/full_catalogue.page - acmt

@TODO Consider moving checkBalances, updateBalance to here

Accounts (acmt) transactions are as follows:
1  - AccountOpeningInstructionV05
2  - AccountDetailsConfirmationV05
3  - AccountModificationInstructionV05
5  - RequestForAccountManagementStatusReportV03
6  - AccountManagementStatusReportV04
7  - AccountOpeningRequestV02
8  - AccountOpeningAmendmentRequestV02
9  - AccountOpeningAdditionalInformationRequestV02
10 - AccountRequestAcknowledgementV02
11 - AccountRequestRejectionV02
12 - AccountAdditionalInformationRequestV02
13 - AccountReportRequestV02
14 - AccountReportV02
15 - AccountExcludedMandateMaintenanceRequestV02
16 - AccountExcludedMandateMaintenanceAmendmentRequestV02
17 - AccountMandateMaintenanceRequestV02
18 - AccountMandateMaintenanceAmendmentRequestV02
19 - AccountClosingRequestV02
20 - AccountClosingAmendmentRequestV02
21 - AccountClosingAdditionalInformationRequestV02
22 - IdentificationModificationAdviceV02
23 - IdentificationVerificationRequestV02
24 - IdentificationVerificationReportV02

### Custom functionality
1000 - ListAllAccounts (@FIXME Used for now by anyone, close down later)
1001 - ListSingleAccount
1002 - CheckAccountByID

*/

/* acmt~1~
   AccountHolderGivenName~
   AccountHolderFamilyName~
   AccountHolderDateOfBirth~
   AccountHolderIdentificationNumber~
   AccountHolderContactNumber1~
   AccountHolderContactNumber2~
   AccountHolderEmailAddress~
   AccountHolderAddressLine1~
   AccountHolderAddressLine2~
   AccountHolderAddressLine3~
   AccountHolderPostalCode
*/
type AccountHolder struct {
	AccountNumber string
	BankNumber    string
}

type AccountHolderDetails struct {
	AccountNumber        string
	BankNumber           string
	GivenName            string
	FamilyName           string
	DateOfBirth          string
	IdentificationNumber string
	ContactNumber1       string
	ContactNumber2       string
	EmailAddress         string
	AddressLine1         string
	AddressLine2         string
	AddressLine3         string
	PostalCode           string
}

type AccountDetails struct {
	AccountNumber     string
	BankNumber        string
	AccountHolderName string
	AccountBalance    float64
	Overdraft         float64
	AvailableBalance  float64
	Timestamp         int
}

// Set up some defaults
const (
	BANK_NUMBER       = "a0299975-b8e2-4358-8f1a-911ee12dbaac"
	OPENING_BALANCE   = 100.
	OPENING_OVERDRAFT = 0.
)

func ProcessAccount(data []string) (result string, err error) {
	if len(data) < 3 {
		return "", errors.New("accounts.ProcessAccount: Not enough fields, minimum 3")
	}

	acmtType, err := strconv.ParseInt(data[2], 10, 64)
	if err != nil {
		return "", errors.New("accounts.ProcessAccount: Could not get ACMT type")
	}

	// Switch on the acmt type
	switch acmtType {
	case 1, 7:
		/*
		   @TODO
		   The differences between AccountOpeningInstructionV05 and AccountOpeningRequestV02 will be explored in detail, for now we treat the same - open an account
		*/
		result, err = openAccount(data)
		break
	case 1000:
		result, err = fetchAccounts(data)
		break
	case 1001:
		result, err = fetchSingleAccount(data)
		break
	case 1002:
		if len(data) < 3 {
			err = errors.New("accounts.ProcessAccount: Not all fields present")
			return
		}
		result = fetchSingleAccountByID(data)
		break
	default:
		err = errors.New("accounts.ProcessAccount: ACMT transaction code invalid")
		break
	}

	return
}

func openAccount(data []string) (result string, err error) {
	// Validate string against required info/length
	if len(data) < 14 {
		err = errors.New("accounts.openAccount: Not all fields present")
		//@TODO Add to documentation rather than returning here
		//result = "ERROR: acmt transactions must be as follows:acmt~AcmtType~AccountHolderGivenName~AccountHolderFamilyName~AccountHolderDateOfBirth~AccountHolderIdentificationNumber~AccountHolderContactNumber1~AccountHolderContactNumber2~AccountHolderEmailAddress~AccountHolderAddressLine1~AccountHolderAddressLine2~AccountHolderAddressLine3~AccountHolderPostalCode"
		return
	}

	// Test: acmt~1~Kyle~Redelinghuys~19000101~190001011234098~1112223456~~email@domain.com~Physical Address 1~~~1000
	// Check if account already exists, check on ID number
	accountHolder := getAccountMeta(data[6])
	fmt.Println(accountHolder)
	if accountHolder.AccountNumber != "" {
		return "", errors.New("accounts.openAccount" + accountHolder.AccountNumber + "~Account already open")
	}

	// @FIXME: Remove new line from data
	data[len(data)-1] = strings.Replace(data[len(data)-1], "\n", "", -1)

	// Create account
	accountHolderObject, err := setAccountDetails(data)
	if err != nil {
		return "", err
	}
	accountHolderDetailsObject, err := setAccountHolderDetails(data)
	if err != nil {
		return "", err
	}
	createdAccountHolder, err := createAccount(accountHolderObject, accountHolderDetailsObject)
	if err != nil {
		return "", err
	}

	fmt.Println(createdAccountHolder)
	result = createdAccountHolder.AccountNumber
	return
}

func setAccountDetails(data []string) (accountDetails AccountDetails, err error) {
	if data[4] == "" {
		return AccountDetails{}, errors.New("accounts.setAccountDetails: Family name cannot be empty")
	}
	if data[3] == "" {
		return AccountDetails{}, errors.New("accounts.setAccountDetails: Given name cannot be empty")
	}
	accountDetails.BankNumber = BANK_NUMBER
	accountDetails.AccountHolderName = data[4] + "," + data[3] // Family Name, Given Name
	accountDetails.AccountBalance = OPENING_BALANCE
	accountDetails.Overdraft = OPENING_OVERDRAFT
	accountDetails.AvailableBalance = OPENING_BALANCE + OPENING_OVERDRAFT

	return
}

func setAccountHolderDetails(data []string) (accountHolderDetails AccountHolderDetails, err error) {
	if len(data) < 14 {
		return AccountHolderDetails{}, errors.New("accounts.setAccountHolderDetails: Not all field values present")
	}
	//@TODO: Test date parsing in format ddmmyyyy
	if data[4] == "" {
		return AccountHolderDetails{}, errors.New("accounts.setAccountHolderDetails: Family name cannot be empty")
	}
	if data[3] == "" {
		return AccountHolderDetails{}, errors.New("accounts.setAccountHolderDetails: Given name cannot be empty")
	}

	// @TODO Integrity checks
	accountHolderDetails.BankNumber = BANK_NUMBER
	accountHolderDetails.GivenName = data[3]
	accountHolderDetails.FamilyName = data[4]
	accountHolderDetails.DateOfBirth = data[5]
	accountHolderDetails.IdentificationNumber = data[6]
	accountHolderDetails.ContactNumber1 = data[7]
	accountHolderDetails.ContactNumber2 = data[8]
	accountHolderDetails.EmailAddress = data[9]
	accountHolderDetails.AddressLine1 = data[10]
	accountHolderDetails.AddressLine2 = data[11]
	accountHolderDetails.AddressLine3 = data[12]
	accountHolderDetails.PostalCode = data[13]

	return
}

func fetchAccounts(data []string) (result string, err error) {
	// Fetch all accounts. This fetches non-sensitive information (no balances)
	accounts := getAllAccountDetails()

	// Parse into nice result string
	jsonAccounts, err := json.Marshal(accounts)
	if err != nil {
		return "", err
	}

	result = "1~" + string(jsonAccounts)
	return
}

func fetchSingleAccount(data []string) (result string, err error) {
	// Fetch user account. Must be user logged in
	tokenUser := appauth.GetUserFromToken(data[0])
	account := getSingleAccountDetail(tokenUser)

	// Parse into nice result string
	jsonAccount, err := json.Marshal(account)
	if err != nil {
		return "", err
	}

	result = "1~" + string(jsonAccount)
	return
}

func fetchSingleAccountByID(data []string) (result string) {
	// Format: appauth~1002~USERID
	userID := data[2]

	userAccountNumber := getSingleAccountNumberByID(userID)

	result = userAccountNumber
	return
}
