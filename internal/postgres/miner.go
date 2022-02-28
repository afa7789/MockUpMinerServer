package postgres

import (
	"context"
	"fmt"

	"gitlab.com/afa7789/luxor_challenge/domain"
)

// GetMiner will get it from the DB
func (c *PGClient) GetMiner(id int64) (*domain.Miner, error) {
	sqlStatement := `SELECT * FROM miners WHERE id=$1;`
	var miner domain.Miner
	row := c.db.QueryRow(context.Background(), sqlStatement, id)
	if err := row.Scan(
		&miner.ID,
		&miner.Authorized,
		&miner.CreatedAt,
	); err != nil {
		return nil, fmt.Errorf("%v: %v", "GetMiner at Scan", err)
	}
	return &miner, nil
}

// Inserting it as true for now
func (c *PGClient) InsertMiner() (*int64, error) {
	sqlStatement := `
	INSERT INTO miners (authorized)
	VALUES (FALSE) RETURNING id`

	var response int64
	// insert
	err := c.db.QueryRow(context.Background(), sqlStatement).Scan(&response)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", "InsertMiner at Scan", err)
	}

	return &response, nil
}

// Authorize the miner
func (c *PGClient) AuthorizeMiner(id int64) error {
	sqlStatement := `
	UPDATE miners
	SET authorized = TRUE
	WHERE id = $1`

	// Exec to try and authorize the miner
	_, err := c.db.Exec(context.Background(), sqlStatement, id)
	if err != nil {
		return fmt.Errorf("%v: %v", "AuthorizeMiner at Scan", err)
	}

	return nil
}

// Unauthorize the miner
func (c *PGClient) UnauthorizeMiner(id int64) error {
	sqlStatement := `
	UPDATE miners
	SET authorized = FALSE
	WHERE id = $1`

	// Exec to try and authorize the miner
	_, err := c.db.Exec(context.Background(), sqlStatement, id)
	if err != nil {
		return fmt.Errorf("%v: %v", "AuthorizeMiner at Scan", err)
	}

	return nil
}
