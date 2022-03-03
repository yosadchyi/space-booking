package booking

import (
	"database/sql"
	"log"

	"github.com/yosadchyi/space-booking/pkg/spacex"
)

type Service struct {
	launchpadRepository   LaunchpadRepository
	launchRepository      LaunchRepository
	destinationRepository DestinationRepository
	importer              DataImporter
	mainRepository        MainRepository
}

func NewService(db *sql.DB) *Service {
	spaceXClient := spacex.NewClient()
	launchpadRepository := NewLaunchpadRepository(db)
	launchRepository := NewLaunchRepository(db)

	return &Service{
		launchpadRepository:   launchpadRepository,
		launchRepository:      launchRepository,
		destinationRepository: NewDestinationRepository(db),
		importer:              NewDataImporter(spaceXClient, launchpadRepository, launchRepository),
		mainRepository:        nil,
	}
}

func (s *Service) Init() {
	log.Println("importing launchpad data...")
	err := s.importer.ImportLaunchpads()
	if err != nil {
		log.Fatalf("error importing launchpad data: %s", err)
	}

	log.Println("importing upcoming launches data...")
	err = s.importer.ImportUpcomingSpaceXLaunches()
	if err != nil {
		log.Fatalf("error importing upcoming launches: %s", err)
	}

	log.Println("data import finished successfully")
}

func (s *Service) GetAllLaunchpads() (AllLaunchpadsResponse, error) {
	launchpads, err := s.launchpadRepository.GetAllActive()
	return launchpads, err
}

func (s *Service) GetAllDestinations() (AllDestinationsResponse, error) {
	destinations, err := s.destinationRepository.GetAll()
	return destinations, err
}
