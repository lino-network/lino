package repository

const (
	getTopContentStmt = `
SELECT
    permlink,
    timestamp
    
FROM
    topContent
WHERE
    timestamp = ?
`
	insertTopContentStmt = `
INSERT INTO
topContent(permlink, timestamp)
VALUES
   (?, ?)
`
)
