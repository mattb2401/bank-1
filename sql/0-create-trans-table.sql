CREATE TABLE IF NOT EXISTS transactions (
`id` int NOT NULL AUTO_INCREMENT,
`transaction` varchar(4) NOT NULL, 
`type` int NOT NULL, 
`senderAccountNumber` CHAR(36) NOT NULL, 
`senderBankNumber` CHAR(36) NOT NULL, 
`receiverAccountNumber` CHAR(36) NOT NULL, 
`receiverBankNumber` CHAR(36) NOT NULL, 
`transactionAmount` float NOT NULL, 
`feeAmount` float NOT NULL, 
`timestamp` int NOT NULL, 
PRIMARY KEY (`id`)
);
