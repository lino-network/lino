package repository

const (
	getPostStmt = `
SELECT
    author,
    postID,
    title,
    content,
    parentAuthor,
    parentPostID,
    sourceAuthor,
    sourcePostID,
    links,
    createdAt,
    totalDonateCount,
    totalReward
FROM
    post
WHERE
    author = ?
`
	insertPostStmt = `
INSERT INTO
post(author, postID, title, content, parentAuthor, parentPostID, sourceAuthor, sourcePostID, links, createdAt, totalDonateCount, totalReward)
VALUES
   (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`
)
