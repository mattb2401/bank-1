package accounts

import (
	"fmt"
	"log"
	"strconv"
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
	GivenName            string
	FamilyName           string
	DateOfBirth          int
	IdentificationNumber string
	ContactNumber1       string
	ContactNumber2       string
	EmailAddress         string
	AddressLine1         string
	AddressLine2         string
	AddressLine3         string
	PostalCode           int
}

func ProcessAccount(data []string) (result string) {
	acmtType, err := strconv.ParseInt(data[1], 10, 64)
	if err != nil {
		fmt.Println("Could not get type of ACMT transaction")
		log.Fatal(err)
		return
	}

	// Switch on the acmt type
	switch acmtType {
	case 1:
	case 7:
		/*
		   @TODO
		   The differences between AccountOpeningInstructionV05 and AccountOpeningRequestV02 will be explored in detail, for now we treat the same - open an account
		*/
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
		fmt.Println("ERROR: Not all fields present for acocunt creation")
		result = "ERROR: acmt transactions must be as follows:acmt~AcmtType~AccountHolderGivenName~AccountHolderFamilyName~AccountHolderDateOfBirth~AccountHolderIdentificationNumber~AccountHolderContactNumber1~AccountHolderContactNumber2~AccountHolderEmailAddress~AccountHolderAddressLine1~AccountHolderAddressLine2~AccountHolderAddressLine3~AccountHolderPostalCode"
		return
	}
	// Check if account already exists
	// Create account
	return
}
