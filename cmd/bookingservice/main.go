package main

import (
	"log"
	"net/http"
	"os"

	"github.com/yosadchyi/space-booking/pkg/booking"
	"github.com/yosadchyi/space-booking/pkg/db"
)

const defaultConnInfo = "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
const defaultBindAddr = ":8080"

func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	database, err := db.Connect(getenv("DB_CONN_INFO", defaultConnInfo))
	if err != nil {
		log.Fatal("can't connect to database")
	}
	service := booking.NewService(database)
	service.Init()
	handler := booking.NewHandler(service)
	bindAddr := getenv("HTTP_BIND_ADDR", defaultBindAddr)
	log.Printf("listening at '%s'...", bindAddr)
	err = http.ListenAndServe(bindAddr, handler)
	if err != nil {
		log.Fatal("can't start http server")
	}
}
