package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
)

type PGClient struct {
	db *pgx.Conn
}

func NewPGClient(port, host, dbname, user, password string) *PGClient {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	psqlconn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbname)

	db, err := pgx.Connect(context.Background(), psqlconn)
	if err != nil {
		log.Fatalf("error creating database: %v", err)
	}
	err = db.Ping(context.Background())
	if err != nil {
		log.Fatalf("error creating database: %v", err)
	}
	return &PGClient{db: db}
}

func (c *PGClient) Close() {
	c.db.Close(context.Background())
}
