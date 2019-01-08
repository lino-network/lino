package repository

const (
	insertUserStmt = `
INSERT INTO
user(username, createdAt, resetPubKey, transactionPubKey, appPubKey, saving, sequence)
VALUES
   (?, ?, ?, ?, ?, ?, ?)
`
	getUserStmt = `
    SELECT
        username,
        createdAt,
        resetPubKey,
        transactionPubKey,
        appPubKey,
        saving,
        sequence
    FROM
        user
    WHERE
        username = ?
`
	increaseSeqByOneStmt = `
    UPDATE user
    SET
        sequence=sequence+1
    WHERE
        username = ?
    `

	updatePubKeyStmt = `
    UPDATE user
    SET
        resetPubKey=?, transactionPubKey=?, appPubKey=?
    WHERE
        username = ?
    `
	updateBalanceStmt = `
    UPDATE user
    SET
        saving=?
    WHERE
        username = ?
    `
)
