package repository

const (
	getInflationStmt = `
SELECT
    id,
    infraPool,
    devPool,
    creatorPool,
    validatorPool,
    infraInflation,
    devInflation,
    creatorInflation,
    validatorInflation,
    timestamp
FROM
    inflation
WHERE
    timestamp = ?
`
	insertInflationStmt = `
INSERT INTO
inflation(infraPool, devPool, creatorPool, validatorPool, infraInflation, devInflation, creatorInflation, validatorInflation, timestamp)
VALUES
   (?, ?, ?, ?, ?, ?, ?, ?, ?)
`
)
