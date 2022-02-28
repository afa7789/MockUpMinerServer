package postgres

import (
	"context"
	"fmt"

	"gitlab.com/afa7789/luxor_challenge/domain"
)

// intert the entry
func (c *PGClient) InsertEntry(ser domain.StratumEntryRPC, success bool, ip string) error {
	sqlStatement := `
	INSERT INTO entries (miner_id,method,params,success,ip)
	VALUES ($1,$2,$3,$4,$5)
	`
	// insert values with exec
	_, err := c.db.Exec(context.Background(), sqlStatement,
		ser.ID,
		ser.Method,
		fmt.Sprintf("%v", ser.Params),
		success,
		ip,
	)
	if err != nil {
		return fmt.Errorf("%v: %v", "InsertEntry at Scan", err)
	}
	return nil
}
