package repository

const (
	getTopContentStmt = `
SELECT
    id,
    permlink,
    timestamp
FROM
    topContent
WHERE
    permlink = ?
`
	insertTopContentStmt = `
INSERT INTO
topContent(permlink, timestamp)
VALUES
   (?, ?)
`
)
