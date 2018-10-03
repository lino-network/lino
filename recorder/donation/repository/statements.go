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
    coinDayDonated
FROM
    donation
WHERE
    username = ?
`
	insertDonationStmt = `
REPLACE INTO
donation(username, seq, dp, permlink, amount, fromApp, coinDayDonated)
VALUES
   (?, ?, ?, ?, ?, ?, ?)
`
)
