package repository

const (
	insertRewardStmt = `
INSERT INTO
reward(username, totalIncome, originalIncome, frictionIncome, inflationIncome, unclaimReward, createdAt)
VALUES
   (?, ?, ?, ?, ?, ?, ?)
`
	getRewardStmt = `
    SELECT
        id,
        username,
        totalIncome,
        originalIncome,
        frictionIncome,
        inflationIncome,
        unclaimReward,
        createdAt
    FROM
        reward
    WHERE
        username = ?
`
)
