package main

import (
	"log"
)

func main() {

	db, err := NewPostgresDb()
	if err != nil {
		log.Fatal(err)
	}

	server := NewApiServer(":8001", db)
	server.Run()
}
