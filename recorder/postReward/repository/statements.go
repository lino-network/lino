package repository

const (
	getPostRewardStmt = `
SELECT
    permlink,
    reward,
    penaltyScore,
    timestamp,
    evaluate,
    original,
    consumer
FROM
    postReward
WHERE
    timestamp = ?
`
	insertPostRewardStmt = `
INSERT INTO
postReward(permlink, reward, penaltyScore, timestamp, evaluate, original, consumer)
VALUES
   (?, ?, ?, ?, ?, ?, ?)
`
)
