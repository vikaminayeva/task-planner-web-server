package main

import (
	"log"
	"planner/pck/db"
	"planner/pck/server"
)

func main() {

	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	server.Start()
}
