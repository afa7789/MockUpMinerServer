package domain

// ServerContent is a struct to be passed on the Server Creation
type ServerContent struct {
	Port    int // Port used in the server
	Manager StratumManager
}

// EnvVariables
type EnvVariables struct {
	DatabaseHost     string
	DatabasePort     int
	DatabaseName     string
	DatabaseUsername string
	DatabasePassword string
}
