package booking

import (
	"database/sql"
	"log"
	"time"

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
	return s.launchpadRepository.GetAllActive()
}

func (s *Service) GetAllDestinations() (AllDestinationsResponse, error) {
	return s.destinationRepository.GetAll()
}

func (s *Service) GetAllBookings() (AllBookingsResponse, error) {
	return s.mainRepository.GetAll()
}

func (s *Service) DeleteBooking(id string) error {
	tx, _ := s.db.Begin()
	err := s.mainRepository.DeleteTx(tx, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	err = s.launchRepository.Delete(tx, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	_ = tx.Commit()
	return nil
}

func (s *Service) AddBooking(request Request) (interface{}, error) {
	if time.Time(request.LaunchDate).Before(time.Now()) {
		return &ErrorResponse{
			Code:    "LAUNCH_DATE_IS_IN_PAST",
			Message: "launch date is in past",
		}, nil
	}
	active, err := s.launchpadRepository.ExistsAndIsActive(request.LaunchpadId)
	if err != nil {
		return nil, err
	}
	if !active {
		return &ErrorResponse{
			Code:    "LAUNCHPAD_NOT_AVAILABLE",
			Message: "launchpad does not exists or is inactive",
		}, nil
	}
	_, err = uuid.Parse(request.DestinationId)
	if err != nil {
		// bad uuid
		return &ErrorResponse{
			Code:    "DESTINATION_ID_IS_INVALID",
			Message: "malformed destination id",
		}, nil
	}
	exists, err := s.destinationRepository.Exists(request.DestinationId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return &ErrorResponse{
			Code:    "DESTINATION_IS_INVALID",
			Message: "destination does not exists",
		}, nil
	}
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	tx, _ := s.db.Begin()
	launches, err := s.launchRepository.GetAllFromLaunchpadAtDate(request.LaunchpadId, request.LaunchDate)

	if err != nil {
		return nil, err
	}
	if len(launches) > 0 {
		_ = tx.Rollback()
		return &ErrorResponse{
			Code:    "LAUNCHPAD_BUSY",
			Message: "launchpad is busy at given date",
		}, nil
	}
	booking := Booking{
		Id:            newUUID.String(),
		FirstName:     request.FirstName,
		LastName:      request.LastName,
		Gender:        request.Gender,
		Birthday:      request.Birthday,
		LaunchpadId:   request.LaunchpadId,
		DestinationId: request.DestinationId,
		LaunchDate:    request.LaunchDate,
		LaunchId:      newUUID.String(), // for simplicity use same id for launch as for booking
	}

	weekLaunches, err := s.launchRepository.GetWeekLaunches(tx, booking.LaunchpadId, booking.LaunchDate)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	for _, launch := range weekLaunches {
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
				Message: "launchpad already used/booked for this destination during requested week",
			}, nil
		}
	}
	err = s.launchRepository.AddTx(tx, &Launch{
		Id:          booking.Id, // using same id as for booking for simplicity
		LaunchpadId: booking.LaunchpadId,
		Date:        booking.LaunchDate,
	})
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	err = s.mainRepository.AddTx(tx, &booking)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	_ = tx.Commit()

	return &SuccessResponse{Id: booking.Id}, nil
}

func (s *Service) PingDb() error {
	_, err := s.db.Exec("SELECT * FROM launchpad")
	return err
}
