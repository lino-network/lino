package repository

const (
	getInflationStmt = `
SELECT
    infraPool,
    devPool,
    creatorPool,
    validatorPool,
    timestamp
FROM
    inflation
WHERE
    timestamp = ?
`
	insertInflationStmt = `
INSERT INTO
inflation(infraPool, devPool, creatorPool, validatorPool, timestamp)
VALUES
   (?, ?, ?, ?, ?)
`
)
