/*
Not sure if this will need more fields, for now we keep one line which is the total of the 
bank's holding account.

Transactional data can be found in the transactions account, and for the bank on the bank
transactional account
*/
CREATE TABLE IF NOT EXISTS `bank_account` (
`id` int NOT NULL AUTO_INCREMENT,
`balance` float NOT NULL,
`timestamp` int NOT NULL, 
PRIMARY KEY (`id`)
);

/* This table must be seeded */
INSERT INTO `bank_account` VALUES (1, 0, 0)

CREATE TABLE IF NOT EXISTS bank_transactions (
`id` int NOT NULL AUTO_INCREMENT,
`transaction` varchar(4), 
`type` int, 
`senderBankNumber` int, 
`receiverBankNumber` int, 
`transactionAmount` float, 
`feeAmount` float, 
`timestamp` int, 
PRIMARY KEY (`id`)
);
