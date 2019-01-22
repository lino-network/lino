package repository

const (
	insertGrantPermissionStmt = `
INSERT INTO
grantpermission(username, authTo, permission, createdAt, expiresAt, amount)
VALUES
   (?, ?, ?, ?, ?, ?)
`
	getGrantPermissionStmt = `
    SELECT
        username,
        authTo,
        permission,
        createdAt,
        expiresAt,
        amount
    FROM
        grantpermission
    WHERE
        username = ? AND authTo = ?
`
	updateAmountStmt = `
    UPDATE grantpermission
    SET
        amount=?
    WHERE
        username = ? AND authTo = ?
    `

	updateGrantPermissionStmt = `
    UPDATE grantpermission
    SET
    permission=?, createdAt=?, expiresAt=?, amount=?
    WHERE
        username = ? AND authTo = ?
    `
	deletePermissionStmt = `
    delete from grantpermission WHERE
    username = ? AND authTo = ?
    `
)
