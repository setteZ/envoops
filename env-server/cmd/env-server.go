package main

import (
	"env-server/internal/app/client"
	"env-server/internal/app/server"
	"env-server/internal/database"
	"log"
)

func main() {
	err := database.CheckDB()
	if err != nil {
		log.Panicf("error while checking the database: '%s'", err.Error())
	}

	// mqtt client
	go client.Run()

	// web server
	server.Run()
}
