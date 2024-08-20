package main

import (
	"github.com/TeslaMode1X/gormTest/connection"
	"github.com/TeslaMode1X/gormTest/router"
	"log"
	"net/http"
)

func main() {
	r := router.Router()

	if connection.DB != nil {
		log.Println("Connected to DB")
	}

	log.Fatal(http.ListenAndServe(":4000", r))
}
