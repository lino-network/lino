package repository

const (
	getDonationStmt = `
SELECT
    username,
    seq,
    dp,
    permlink,
    amount,
    fromApp,
    coinDayDonated,
    reputation,
    timestamp,
    evaluateResult
FROM
    donation
WHERE
    username = ?
`
	insertDonationStmt = `
INSERT INTO
donation(username, seq, dp, permlink, amount, fromApp, coinDayDonated, reputation, timestamp, evaluateResult)
VALUES
   (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`
)
