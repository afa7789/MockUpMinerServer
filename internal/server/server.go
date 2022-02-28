package server

import (
	"fmt"
	"log"
	"net"
	"strings"

	"gitlab.com/afa7789/luxor_challenge/domain"
)

type Server struct {
	port    int
	stratum domain.StratumManager
}

// NewServer, constructor function
func NewServer(cont *domain.ServerContent) *Server {
	return &Server{
		port:    cont.Port,
		stratum: cont.Manager,
	}
}

// Start server function
func (s *Server) Start() {
	log.Printf("starting the server")
	// set's up the listening
listen:
	pl := fmt.Sprintf(":%d", s.port)
	listener, err := net.Listen("tcp", pl)
	log.Printf("listening at %s", pl)
	// handle error of port in use to go the next port
	if err != nil {
		// Using this error treatment to try again on next port
		if strings.Contains(err.Error(), "address already in use") {
			log.Printf("port already in use::%d", s.port)
			s.port++
			log.Printf("trying next port::%d", s.port)
			goto listen
		} else {
			panic(err)
		}
	}

	// handleListener in another function
	s.handleListener(listener)
}

// handleListener uses a for to handle the listener to accepts connections
// this is used to handle them concurrently
func (s *Server) handleListener(l net.Listener) {
	// for each
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go s.stratum.HandleConn(conn)
	}
}
