package accounts

import (
	"fmt"
	"log"
	"strconv"
	"strings"
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
type AccountHolderDetails struct {
	AccountNumber        string
	BankNumber           string
	GivenName            string
	FamilyName           string
	DateOfBirth          int64
	IdentificationNumber string
	ContactNumber1       string
	ContactNumber2       string
	EmailAddress         string
	AddressLine1         string
	AddressLine2         string
	AddressLine3         string
	PostalCode           int64
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

func ProcessAccount(data []string) (result string) {
	acmtType, err := strconv.ParseInt(data[1], 10, 64)
	if err != nil {
		fmt.Println("Could not get type of ACMT transaction")
		log.Fatal(err)
		return
	}

	// Switch on the acmt type
	switch acmtType {
	case 1, 7:
		/*
		   @TODO
		   The differences between AccountOpeningInstructionV05 and AccountOpeningRequestV02 will be explored in detail, for now we treat the same - open an account
		*/
		fmt.Println("Processing accountdata")
		result = openAccount(data)
		break
	default:
		break
	}

	return
}

func openAccount(data []string) (result string) {
	// Validate string against required info/length
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
	if len(data) < 13 {
		fmt.Println("ERROR: Not all fields present for account creation")
		result = "ERROR: acmt transactions must be as follows:acmt~AcmtType~AccountHolderGivenName~AccountHolderFamilyName~AccountHolderDateOfBirth~AccountHolderIdentificationNumber~AccountHolderContactNumber1~AccountHolderContactNumber2~AccountHolderEmailAddress~AccountHolderAddressLine1~AccountHolderAddressLine2~AccountHolderAddressLine3~AccountHolderPostalCode"
		return
	}

	// Test: acmt~1~Kyle~Redelinghuys~19000101~190001011234098~1112223456~~email@domain.com~Physical Address 1~~~1000
	// Check if account already exists, check on ID number
	accountHolder := getAccountMeta(data[5])
	fmt.Println(accountHolder)
	if accountHolder.AccountNumber != "" {
		return "1~" + accountHolder.AccountNumber + "~Account already open."
	}

	// Remove new line from data
	data[len(data)-1] = strings.Replace(data[len(data)-1], "\n", "", -1)

	// Create account
	accountHolderObject := setAccountDetails(data)
	accountHolderDetailsObject := setAccountHolderDetails(data)
	createdAccountHolder := createAccount(accountHolderObject, accountHolderDetailsObject)

	fmt.Println(createdAccountHolder)
	result = createdAccountHolder.AccountNumber
	return
}

func setAccountDetails(data []string) (accountDetails AccountDetails) {
	// @TODO Integrity checks
	accountDetails.BankNumber = BANK_NUMBER
	accountDetails.AccountHolderName = data[3] + "," + data[2] // Family Name, Given Name
	accountDetails.AccountBalance = OPENING_BALANCE
	accountDetails.Overdraft = OPENING_OVERDRAFT
	accountDetails.AvailableBalance = OPENING_BALANCE + OPENING_OVERDRAFT

	return
}

func setAccountHolderDetails(data []string) (accountHolderDetails AccountHolderDetails) {
	dob, err := strconv.ParseInt(data[4], 10, 64)
	if err != nil {
		fmt.Println("ERROR: Could not convert date")
		return
	}

	postalCode, err := strconv.ParseInt(data[12], 10, 64)
	if err != nil {
		fmt.Println("ERROR: Could not convert postal code")
		return
	}

	// @TODO Integrity checks
	accountHolderDetails.BankNumber = BANK_NUMBER
	accountHolderDetails.GivenName = data[2]
	accountHolderDetails.FamilyName = data[3]
	accountHolderDetails.DateOfBirth = dob
	accountHolderDetails.IdentificationNumber = data[5]
	accountHolderDetails.ContactNumber1 = data[6]
	accountHolderDetails.ContactNumber2 = data[7]
	accountHolderDetails.EmailAddress = data[8]
	accountHolderDetails.AddressLine1 = data[9]
	accountHolderDetails.AddressLine2 = data[10]
	accountHolderDetails.AddressLine3 = data[11]
	accountHolderDetails.PostalCode = postalCode

	return
}
