package domain

import "time"

// look at the file in etc/postgres/0_init.sql
// Mine representation in GO
type Miner struct {
	ID         int64     `db:"id"`
	Authorized bool      `db:"authorized"`
	CreatedAt  time.Time `db:"created_at"`
}

// Entry representation in GO
type Entry struct {
	ID         int64     `db:"id"`
	MinerID    int64     `db:"miner_id"`
	Authorized bool      `db:"authorized"`
	CreatedAt  time.Time `db:"created_at"`
}

// BlackList representation in GO
type BlackList struct {
	ID int64  `db:"id"`
	IP string `db:"ip"`
}
