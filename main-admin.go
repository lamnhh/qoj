package main

import (
	"qoj/server"
	"qoj/server/config"
)

func main() {
	config.InitialiseDatabaseConnection()
	app := server.InitialiseAdminApp()

	_ = app.Run(":3001")
}
