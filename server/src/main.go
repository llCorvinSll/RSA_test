package main

import (
	"decrypt_test"
	"log"
)

func main() {
	routes := decrypt_test.GetRoutes()

	err := routes.Run()

	if err != nil {
		log.Fatal("dont started")
	}
}