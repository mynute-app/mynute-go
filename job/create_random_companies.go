package job
// This script creates a specified number of 
// companies with random data for development purposes only.
// To run the script, execute command below in 
// terminal at <repo_root> level:
// go run job/create_random_companies.go 5

import (
	"mynute-go/core"
	srcLib "mynute-go/core/src/lib"
	"mynute-go/job/lib"
	"os"
	"strconv"
)

func main(qtty int) {
	srcLib.LoadEnv()
	app_env := os.Getenv("APP_ENV")
	if app_env != "dev" {
		panic("This script is intended to run only in development environment. Set APP_ENV to 'dev' in your .env file.")
	}
	SetupServer := core.NewServer().Run("parallel")
	lib.CreateCompaniesRandomly(qtty)
	SetupServer.Shutdown()
}

// Get quantity of companies to create from command line argument
func GetQuantityFromArgs() int {
	if len(os.Args) < 2 {
		panic("Please provide the number of companies to create as a command line argument.")
	}
	qtty, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic("Invalid number format. Please provide a valid integer.")
	}
	return qtty
}

func init() {
	qtty := GetQuantityFromArgs()
	main(qtty)
}
