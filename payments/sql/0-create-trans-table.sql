CREATE TABLE IF NOT EXISTS transactions (
`id` int NOT NULL AUTO_INCREMENT,
`transaction` varchar(4), 
`type` int, 
`senderAccountNum` int, 
`senderBankNum` int, 
`receiverAccountNum` int, 
`receiverBankNum` int, 
`transactionAmount` float, 
`feeAmount` float, 
`timestamp` int, 
PRIMARY KEY (`id`)
);
