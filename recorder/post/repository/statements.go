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
	setRewardStmt = `
    UPDATE post
    SET
      totalReward = ?
    WHERE author = ? AND postID = ?
    `
)
