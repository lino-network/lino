package repository

const (
	getPostRewardStmt = `
SELECT
    id,
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
    permlink = ?
`
	insertPostRewardStmt = `
INSERT INTO
postReward(permlink, reward, penaltyScore, timestamp, evaluate, original, consumer)
VALUES
   (?, ?, ?, ?, ?, ?, ?)
`
)
