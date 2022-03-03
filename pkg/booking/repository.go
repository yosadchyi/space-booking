package booking

import "database/sql"

type Repositories struct {
	Booking   MainRepository
	Launchpad LaunchpadRepository
	Launch    LaunchRepository
}

// MainRepository booking repository
type MainRepository interface {
	Add(booking *Booking) (bool, error)
	GetAll() ([]Booking, error)
}

// LaunchpadRepository repository to access all launchpads
type LaunchpadRepository interface {
	Exists(id string) (bool, error)
	AddOrUpdate(launchpad *Launchpad) error
	GetAllActive() ([]Launchpad, error)
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

type launchRepository struct {
	db *sql.DB
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Booking:   nil,
		Launchpad: NewLaunchpadRepository(db),
		Launch:    NewLaunchRepository(db),
	}
}

// NewLaunchpadRepository create new launchpad repository
func NewLaunchpadRepository(db *sql.DB) LaunchpadRepository {
	return &launchpadRepository{db: db}
}

func (l *launchpadRepository) Exists(id string) (bool, error) {
	row := l.db.QueryRow("SELECT true FROM launchpad WHERE id = $1", id)
	exists := false
	if err := row.Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, err
		}
	}
	return exists, nil
}

func (l *launchpadRepository) AddOrUpdate(launchpad *Launchpad) error {
	query := `
		INSERT INTO launchpad (id, name, status) VALUES ($1, $2, $3)
		ON CONFLICT (id) DO 
		UPDATE SET name = $2, status = $3`
	_, err := l.db.Exec(query, launchpad.Id, launchpad.Name, launchpad.Status)
	return err
}

func (l *launchpadRepository) GetAllActive() ([]Launchpad, error) {
	rows, err := l.db.Query("SELECT id, name, status FROM launchpad WHERE status = 'active' ORDER BY id")
	if err != nil {
		return nil, err
	}

	launchpads := make([]Launchpad, 0, 1)
	for rows.Next() {
		launchpad := Launchpad{}
		if err := rows.Scan(&launchpad.Id, &launchpad.Name, &launchpad.Status); err != nil {
			return launchpads, err
		}
		launchpads = append(launchpads, launchpad)
	}
	return launchpads, nil
}

// NewLaunchRepository creates new launch repository
func NewLaunchRepository(db *sql.DB) LaunchRepository {
	return &launchRepository{db: db}
}

func (l *launchRepository) Exists() (bool, error) {
	panic("implement me")
}

func (l *launchRepository) Add(launch *Launch) error {
	panic("implement me")
}

func (l *launchRepository) GetAtDate() ([]Launch, error) {
	panic("implement me")
}
