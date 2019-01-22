CREATE TABLE `grantpermission`
(
  `username` varchar(45) NOT NULL,
  `authTo` varchar(76) NOT NULL,
  `permission` int(10) NOT NULL,
  `createdAt` datetime NOT NULL,
  `expiresAt` datetime NOT NULL,
  `amount` char(64) NOT NULL DEFAULT "0000000000000000000000000000000000000000000000000000000000000000",
  PRIMARY KEY (`username`,`authto`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE post ADD COLUMN `isDeleted` boolean NOT NULL DEFAULT 0;