package main

import (
	"qoj/server"
	"qoj/server/config"
	"qoj/server/src/queue"
)

func main() {
	queue.InitQueue()
	config.InitialiseDatabaseConnection()
	app := server.InitialiseApp()

	_ = app.Run(":3000")
}