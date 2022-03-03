package booking

import (
	"database/sql"

	"github.com/yosadchyi/space-booking/pkg/spacex"
)

type DataImporter interface {
	ImportLaunchpads() error
	ImportSpaceXLaunches(db *sql.DB) error
}

type dataImporter struct {
	client        spacex.Client
	launchpadRepo LaunchpadRepository
	launchRepo    LaunchRepository
}

func NewDataImporter(client spacex.Client, launchpadRepo LaunchpadRepository, launchRepo LaunchRepository) DataImporter {
	return &dataImporter{
		client:        client,
		launchpadRepo: launchpadRepo,
		launchRepo:    launchRepo,
	}
}

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

func (d *dataImporter) ImportSpaceXLaunches(db *sql.DB) error {
	panic("implement me")
}
