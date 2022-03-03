package main

import (
	"log"
	"net/http"

	"github.com/yosadchyi/space-booking/pkg/booking"
	"github.com/yosadchyi/space-booking/pkg/db"
)

func main() {
	database, err := db.Connect("host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal("can't connect to database")
	}
	service := booking.NewService(database)
	service.Init()
	handler := booking.NewHandler(service)
	log.Println("listening at *:8080...")
	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal("can't start http server")
	}
}
