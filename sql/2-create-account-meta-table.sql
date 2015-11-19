CREATE TABLE IF NOT EXISTS accounts_meta (
`id` int NOT NULL AUTO_INCREMENT,
`accountNumber` char(36) NOT NULL,
`bankNumber` char(36) NOT NULL,
`accountHolderGivenName` text NOT NULL, 
`accountHolderFamilyName` text NOT NULL, 
`accountHolderDateOfBirth` text NOT NULL, 
`accountHolderIdentificationNumber` text NOT NULL, 
`accountHolderContactNumber1` text NOT NULL, 
`accountHolderContactNumber2` text NULL, 
`accountHolderEmailAddress` text NOT NULL, 
`accountHolderAddressLine1` text NOT NULL, 
`accountHolderAddressLine2` text NULL, 
`accountHolderAddressLine3` text NULL, 
`accountHolderPostalCode` text NOT NULL, 
`timestamp` int NOT NULL, 
PRIMARY KEY (`id`)
);
