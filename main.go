package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
	"gitlab.com/afa7789/luxor_challenge/cmd"
)

// port is used to define which port to connect to
var port int

// init is a initialization function being used to get flag values, variables that enters through the command line
func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("error loading .env file")
	}
	flag.IntVar(&port, "port", 8080, "port number")
	flag.Parse()
}

// main
func main() {
	cmd.Start(port)
}
