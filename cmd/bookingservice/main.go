package main

import (
	"log"

	"github.com/yosadchyi/space-booking/pkg/booking"
	"github.com/yosadchyi/space-booking/pkg/db"
	"github.com/yosadchyi/space-booking/pkg/spacex"
)

func main() {
	database, err := db.Connect("host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal("can't connect to database")
	}
	spaceXClient := spacex.NewClient()
	launchpadRepository := booking.NewLaunchpadRepository(database)
	launchRepository := booking.NewLaunchRepository(database)
	importer := booking.NewDataImporter(spaceXClient, launchpadRepository, launchRepository)

	log.Println("importing launchpad data...")
	err = importer.ImportLaunchpads()
	if err != nil {
		log.Fatalf("error importing launchpad data: %s", err)
	}
}
