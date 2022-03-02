package booking

import "database/sql"

// MainRepository booking repository
type MainRepository interface {
	Add(booking *Booking) (bool, error)
	GetAll() ([]Booking, error)
}

// LaunchpadRepository repository to access all launchpads
type LaunchpadRepository interface {
	Exists(id string) (bool, error)
	Add(launchpad *Launchpad) error
	GetAll() ([]Launchpad, error)
}

// LaunchRepository repository to access all launches
type LaunchRepository interface {
	Exists() (bool, error)
	Add(launch *Launch) error
	GetAtDate() ([]Launch, error)
}

type launchpadRepository struct {
	db *sql.DB
}

func (l *launchpadRepository) Exists(id string) (bool, error) {
	panic("implement me")
}

func (l *launchpadRepository) Add(launchpad *Launchpad) error {
	panic("implement me")
}

func (l *launchpadRepository) GetAll() ([]Launchpad, error) {
	panic("implement me")
}

// NewLaunchpadRepository create new launchpad repository
func NewLaunchpadRepository(db *sql.DB) LaunchpadRepository {
	return &launchpadRepository{
		db: db,
	}
}
