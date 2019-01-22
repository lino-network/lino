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
    totalReward,
    isDeleted
FROM
    post
WHERE
    author = ?
`
	insertPostStmt = `
INSERT INTO
post(author, postID, title, content, parentAuthor, parentPostID, sourceAuthor, sourcePostID, links, createdAt, totalDonateCount, totalReward, isDeleted)
VALUES
   (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0)
`
	setRewardStmt = `
    UPDATE post
    SET
      totalReward = ?
    WHERE author = ? AND postID = ?
    `
	deletePostStmt = `
    UPDATE post
    SET
      isDeleted = 1
    WHERE author = ? AND postID = ?
    `
)
