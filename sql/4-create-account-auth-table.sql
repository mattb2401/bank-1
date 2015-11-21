CREATE TABLE IF NOT EXISTS accounts_auth (
`id` int NOT NULL AUTO_INCREMENT,
`accountNumber` char(36) NOT NULL, 
`password` varchar(255) NOT NULL,
`timestamp` int NOT NULL,
PRIMARY KEY (`id`)
);

