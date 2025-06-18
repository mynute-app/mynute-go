package main

import (
	"agenda-kaki-go/core"
	"agenda-kaki-go/core/lib"
	utilsT "agenda-kaki-go/core/test/utils"
	"os"
)

var CompaniesToCreate = 1
var ShouldCreateCompanies = true

func main() {
	lib.LoadEnv()
	app_env := os.Getenv("APP_ENV")
	if app_env != "dev" {
		panic("This script is intended to run only in development environment. Set APP_ENV to 'dev' in your .env file.")
	}
	if ShouldCreateCompanies {
		SetupServer := core.NewServer().Run("parallel")
		utilsT.CreateCompaniesRandomly(CompaniesToCreate)
		SetupServer.Shutdown()
	}
	core.NewServer().Run("listen")
}
