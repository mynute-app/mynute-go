package main

import (
	"agenda-kaki-go/core/lib"
	utilsT "agenda-kaki-go/core/test/utils"
	"os"
)

func main() {
	lib.LoadEnv()
	app_env := os.Getenv("APP_ENV")
	if app_env != "dev" {
		panic("This script is intended to run only in development environment. Set APP_ENV to 'dev' in your .env file.")
	}
	utilsT.CreateCompaniesRandomly(10)
}