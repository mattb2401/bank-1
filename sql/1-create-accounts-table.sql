/*
Accounts table should have the following fields:

id
account num
bank num
account holder name
account type
account balance
overdraft
available balance

*/
CREATE TABLE IF NOT EXISTS accounts (
`id` int NOT NULL AUTO_INCREMENT,
`accountNumber` char(36) NOT NULL, 
`bankNumber` char(36) NOT NULL, 
`accountHolderName` text NOT NULL, 
`accountBalance` float NOT NULL, 
`overdraft` float NOT NULL, 
`availableBalance` float NOT NULL, 
`timestamp` int NOT NULL, 
PRIMARY KEY (`id`)
);
