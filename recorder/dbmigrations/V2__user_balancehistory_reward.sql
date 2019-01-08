CREATE TABLE `balancehistory`
(
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `username` varchar(45) NOT NULL,
  `fromUser` varchar(45) NOT NULL,
  `toUser` varchar(45) NOT NULL,
  `amount` BIGINT NOT NULL,
  `balance` BIGINT NOT NULL,
  `detailType` int(11) NOT NULL,
  `createdAt` datetime NOT NULL,
  `memo` varchar(45) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `reward`
(
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `username` varchar(45) NOT NULL,
  `totalIncome` BIGINT NOT NULL,
  `originalIncome` BIGINT NOT NULL,
  `frictionIncome` BIGINT NOT NULL,
  `inflationIncome` BIGINT NOT NULL,
  `unclaimReward` BIGINT NOT NULL,
  `createdAt` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `user`
(
  `username` varchar(45) NOT NULL PRIMARY KEY,
  `createdAt` datetime NOT NULL,
  `resetPubKey` varchar(76) NOT NULL,
  `transactionPubKey` varchar(76) NOT NULL,
  `appPubKey` varchar(76) NOT NULL,
  `saving` BIGINT NOT NULL,
  `sequence` BIGINT NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `post`
(
  `author` varchar(45) NOT NULL,
  `postID` varchar(45) NOT NULL,
  `title` varchar(50),
  `content` varchar(200),
  `parentAuthor` varchar(45),
  `parentPostID` varchar(50),
  `sourceAuthor` varchar(45),
  `sourcePostID` varchar(50),
  `links` varchar(200),
  `createdAt` datetime NOT NULL,
  `totalDonateCount` BIGINT NOT NULL,
  `totalReward` BIGINT NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;