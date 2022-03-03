package booking

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/yosadchyi/space-booking/pkg/spacex"
)

type Service struct {
	db                    *sql.DB
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
		db:                    db,
		launchpadRepository:   launchpadRepository,
		launchRepository:      launchRepository,
		destinationRepository: NewDestinationRepository(db),
		importer:              NewDataImporter(spaceXClient, launchpadRepository, launchRepository),
		mainRepository:        NewMainRepository(db),
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

func (s *Service) AddBooking(request Request) (interface{}, error) {
	active, err := s.launchpadRepository.ExistsAndIsActive(request.LaunchpadId)
	if err != nil {
		return nil, err
	}
	if !active {
		return &ErrorResponse{
			Code:    "LAUNCHPAD_NOT_AVAILABLE",
			Message: "Launchpad does not exists or is inactive",
		}, nil
	}
	_, err = uuid.Parse(request.DestinationId)
	if err != nil {
		// bad uuid
		return &ErrorResponse{
			Code:    "DESTINATION_ID_IS_INVALID",
			Message: "Malformed destination id",
		}, nil
	}
	exists, err := s.destinationRepository.Exists(request.DestinationId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return &ErrorResponse{
			Code:    "DESTINATION_IS_INVALID",
			Message: "Destination does not exists",
		}, nil
	}
	launches, err := s.launchRepository.GetFromLaunchpadAtDate(request.LaunchpadId, request.LaunchDate)
	if err != nil {
		return nil, err
	}
	if len(launches) > 0 {
		return &ErrorResponse{
			Code:    "LAUNCHPAD_BUSY",
			Message: "Launchpad is busy at given date",
		}, nil
	}
	booking := Booking{
		Id:            "",
		FirstName:     request.FirstName,
		LastName:      request.LastName,
		Gender:        request.Gender,
		Birthday:      request.Birthday,
		LaunchpadId:   request.LaunchpadId,
		DestinationId: request.DestinationId,
		LaunchDate:    request.LaunchDate,
	}

	tx, _ := s.db.Begin()
	weekLaunches, err := s.launchRepository.GetWeekLaunches(tx, booking.LaunchpadId, booking.LaunchDate)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	for _, launch := range weekLaunches {
		// as we're using same Id for Launch record as for booking, it's quite simple
		id, err := s.mainRepository.GetDestinationIdForBookingId(tx, launch.Id)
		switch {
		case err == sql.ErrNoRows:
			continue // skip SpaceX launch
		case err != nil:
			return nil, err
		}
		if request.DestinationId == id {
			return &ErrorResponse{
				Code:    "SAME_DESTINATION_IN_WEEK",
				Message: "Launchpad already used/booked for this destination during requested week",
			}, nil
		}
	}
	err = s.mainRepository.AddTx(tx, &booking)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	err = s.launchRepository.AddTx(tx, &Launch{
		Id:          booking.Id,
		LaunchpadId: booking.LaunchpadId, // same as booking ID, unique in launch table
		Date:        booking.LaunchDate,
	})
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	_ = tx.Commit()

	return &SuccessResponse{Id: booking.Id}, nil
}
