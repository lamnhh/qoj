package main

import (
	"os"
	"qoj/server"
	"qoj/server/config"
)

func main() {
	config.InitialiseDatabaseConnection()
	app := server.InitialiseApp()

	port := os.Getenv("PORT")
	// If env.PORT does not exist, run on port 3000
	if port == "" {
		port = "3000"
	}
	_ = app.Run(":" + port)
}