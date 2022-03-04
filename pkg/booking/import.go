package booking

import (
	"time"

	"github.com/google/uuid"
	"github.com/yosadchyi/space-booking/pkg/spacex"
)

// DataImporter imports the data from external sources
type DataImporter interface {
	ImportLaunchpads() error
	ImportUpcomingSpaceXLaunches() error
}

type dataImporter struct {
	client        spacex.Client
	launchpadRepo LaunchpadRepository
	launchRepo    LaunchRepository
}

// NewDataImporter creates new data importer
func NewDataImporter(client spacex.Client, launchpadRepo LaunchpadRepository, launchRepo LaunchRepository) DataImporter {
	return &dataImporter{
		client:        client,
		launchpadRepo: launchpadRepo,
		launchRepo:    launchRepo,
	}
}

// ImportLaunchpads imports launchpads information, duplicates are updated
func (d *dataImporter) ImportLaunchpads() error {
	launchpads, err := d.client.GetAllLaunchpads()
	if err != nil {
		return err
	}
	for _, launchpad := range launchpads {
		err := d.launchpadRepo.AddOrUpdate(&Launchpad{
			Id:     launchpad.Id,
			Name:   launchpad.Name,
			Status: launchpad.Status,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// ImportUpcomingSpaceXLaunches fetches and stores information about upcoming launches, duplicates are ignored
func (d *dataImporter) ImportUpcomingSpaceXLaunches() error {
	launches, err := d.client.GetUpcomingLaunches()
	if err != nil {
		return err
	}
	for _, launch := range launches {
		newUUID, err := uuid.NewUUID()
		if err != nil {
			return err
		}
		err = d.launchRepo.Add(&Launch{
			Id:          newUUID.String(),
			LaunchpadId: launch.Launchpad,
			Date:        Date(time.Unix(launch.DateUnix, 0)),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
