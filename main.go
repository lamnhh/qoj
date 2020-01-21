package main

import (
	"qoj/server"
	"qoj/server/config"
)

func main() {
	config.InitialiseDatabaseConnection()
	app := server.InitialiseApp()
	_ = app.Run(":3000")
}