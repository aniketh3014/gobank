package main

import (
	"fmt"
	"log"
)

func main() {

	db, err := NewPostgresDb()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", db)
	// server := NewApiServer(":8080")
	// server.Run()
}
