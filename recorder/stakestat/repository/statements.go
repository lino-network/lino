package repository

const (
	getStakeStatStmt = `
SELECT
    totalConsumptionFriction,
    unclaimedFriction,
    totalLinoStake,
    unclaimedLinoStake,
    timestamp
FROM
    stakeStat
WHERE
    timestamp = ?
`
	insertStakeStatStmt = `
INSERT INTO
stakeStat(totalConsumptionFriction, unclaimedFriction, totalLinoStake, unclaimedLinoStake, timestamp)
VALUES
   (?, ?, ?, ?, ?)
`
)
