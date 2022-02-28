package postgres

import (
	"context"
	"fmt"
)

func (c *PGClient) InsertSubscription(ctx context.Context, notification_name string, miner_id int64, bytea []byte) error {
	sqlStatement := `
	INSERT INTO subscriptions (id,miner_id,notification_name,subscribed)
	VALUES ($1,$2,$3,TRUE)
	`
	// insert values with exec
	_, err := c.db.Exec(context.Background(), sqlStatement,
		bytea,
		miner_id,
		notification_name,
	)
	if err != nil {
		return fmt.Errorf("%v: %v", "InsertSubscription at Scan", err)
	}

	return nil
}

func (c *PGClient) Unsubscribe(ctx context.Context, bytea []byte) error {
	sqlStatement := `
	UPDATE subscriptions
	SET subscribed = FALSE
	WHERE id = $1
	`

	// insert values with exec
	_, err := c.db.Exec(context.Background(), sqlStatement,
		bytea,
	)
	if err != nil {
		return fmt.Errorf("%v: %v", "InsertSubscription at Scan", err)
	}
	return nil
}
