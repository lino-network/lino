CREATE TABLE `donation`
(
  `username` varchar(32) NOT NULL,
  `seq` int(10) NOT NULL,
  `dp` BIGINT NOT NULL,
  `permlink` varchar(45) NOT NULL,
  `amount` BIGINT NOT NULL,
  `fromApp` varchar(32),
  `coinDayDonated` BIGINT NOT NULL,
  `reputation` BIGINT NOT NULL,
  `timestamp` BIGINT NOT NULL,
  `evaluateResult` BIGINT NOT NULL,
  PRIMARY KEY (`username`,`seq`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `inflation`
(
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `timestamp` BIGINT NOT NULL,
  `infraPool` BIGINT NOT NULL,
  `devPool` BIGINT NOT NULL,
  `creatorPool` BIGINT NOT NULL,
  `validatorPool` BIGINT NOT NULL,
  `infraInflation` BIGINT NOT NULL,
  `devInflation` BIGINT NOT NULL,
  `creatorInflation` BIGINT NOT NULL,
  `validatorInflation` BIGINT NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `postReward`
(
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `timestamp` BIGINT NOT NULL,
  `permlink` varchar(45) NOT NULL,
  `reward` BIGINT NOT NULL,
  `penaltyScore` varchar(45) NOT NULL,
  `evaluate` BIGINT NOT NULL,
  `original` BIGINT NOT NULL,
  `consumer` varchar(45) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `stake`
(
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `timestamp` BIGINT NOT NULL,
  `username` varchar(45) NOT NULL,
  `amount` BIGINT NOT NULL,
  `op` varchar(45) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `stakeStat`
(
  `timestamp` BIGINT NOT NULL PRIMARY KEY,
  `totalConsumptionFriction` varchar(45) NOT NULL,
  `unclaimedFriction` BIGINT NOT NULL,
  `totalLinoStake` BIGINT NOT NULL,
  `unclaimedLinoStake` BIGINT NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `topContent`
(
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `timestamp` BIGINT NOT NULL,
  `permlink` varchar(45) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;