package postgres

import (
	"context"
	"fmt"
	"log"

	"gitlab.com/afa7789/luxor_challenge/domain"
)

// GetBlackList will return a list of entries with IPS that can be used to block unautorize etc ( just an example//idea )
func (c *PGClient) GetBlackList() ([]domain.BlackList, error) {
	rows, err := c.db.Query(context.Background(), "SELECT * FROM blacklist")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var rowSlice []domain.BlackList
	for rows.Next() {
		var r domain.BlackList
		err := rows.Scan(&r.ID, &r.IP)
		if err != nil {
			return nil, fmt.Errorf("%v: %v", "GetBlackList at Scan in Next", err)
		}
		rowSlice = append(rowSlice, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%v: %v", "GetBlackList at Error in get Rows", err)
	}

	return rowSlice, nil
}
