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
`accountNum` int, 
`bankNum` int, 
`accountHolderName` text, 
`accountBalance` float, 
`overdraft` float, 
`availableBalance` float, 
`timestamp` int, 
PRIMARY KEY (`id`)
);
