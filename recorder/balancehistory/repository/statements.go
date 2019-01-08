package repository

const (
	insertBalanceHistoryStmt = `
INSERT INTO
balancehistory(username, fromUser, toUser, amount, balance, detailType, createdAt, memo)
VALUES
   (?, ?, ?, ?, ?, ?, ?, ?)
`
	getBalanceHistoryStmt = `
    SELECT
        id,
        username,
        fromUser,
        toUser,
        amount,
        balance,
        detailType,
        createdAt,
        memo
    FROM
        balancehistory
    WHERE username = ?
`
)
