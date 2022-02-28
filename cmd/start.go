package cmd

import (
	"os"

	"gitlab.com/afa7789/luxor_challenge/domain"
	"gitlab.com/afa7789/luxor_challenge/internal/middleware"
	"gitlab.com/afa7789/luxor_challenge/internal/postgres"
	"gitlab.com/afa7789/luxor_challenge/internal/server"
	"gitlab.com/afa7789/luxor_challenge/internal/stratum"
)

// Start is the main command of our project it sets part of the code and starts the process
func Start(port int) {

	// creation of the database connection
	d := postgres.NewPGClient(
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
	)

	// creation of the middleware
	m := middleware.NewMiddleware(d)

	// creation of the server struct
	c := &domain.ServerContent{
		Port:    port,
		Manager: stratum.NewManager(m, d),
	}

	// New server using the struct.
	s := server.NewServer(c)

	// Starting the server
	s.Start()

}
