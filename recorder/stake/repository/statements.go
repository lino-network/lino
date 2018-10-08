package repository

const (
	getStakeStmt = `
SELECT
    username,
    amount,
    timestamp,
    op
    
FROM
    stake
WHERE
    username = ?
`
	insertStakeStmt = `
INSERT INTO
stake(username, amount, timestamp, op)
VALUES
   (?, ?, ?, ?)
`
)
