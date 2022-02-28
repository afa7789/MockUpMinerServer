package middleware

import (
	"log"
	"net"

	"gitlab.com/afa7789/luxor_challenge/internal/postgres"
)

type Middleware struct {
	db *postgres.PGClient
}

func NewMiddleware(db *postgres.PGClient) *Middleware {
	return &Middleware{
		db: db,
	}
}

// Authorize will parse through a list of
func (mid *Middleware) Authorize(conn net.Conn, id int64) error {
	addr, ok := conn.RemoteAddr().(*net.TCPAddr)
	// THIS was not requested but I wanted to do something to
	if ok {
		list, err := mid.db.GetBlackList()
		if err != nil {
			log.Printf("error getting the blacklist at Middleware in the authorize function")
		} else {
			// this only consumes processing if list is not empty
			for _, v := range list {
				if addr.IP.String() == v.IP {
					log.Printf("IP blacklisted")
					// mid.db.UnauthorizeMiner(id) doesn't exist
					return nil // should be an error , but it's not required by the challenge this
				}
			}
		}
	}
	m, err := mid.db.GetMiner(id)
	if err != nil {
		log.Printf("\nerror at Authorize in GetMiner:%s\n", err)
	}
	return mid.db.AuthorizeMiner(m.ID)
}

// check if it's authorized
func (mid *Middleware) IsAuthorized(id int64) bool {
	// Do Something here
	m, err := mid.db.GetMiner(id)
	if err != nil {
		log.Printf("error at IsAuthorized in GetMiner:%s", err)
	}
	if m == nil {
		print("m é nulo\n")
	}
	return m.Authorized
}
