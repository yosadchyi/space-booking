package booking

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

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
	Exists(date Date, launchpadId string) (bool, error)
	Add(launch *Launch) error
	GetAtDate(date Date) ([]Launch, error)
	GetAllUpcoming() ([]Launch, error)
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
	return exists(row)
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

func (l *launchRepository) Exists(date Date, launchpadId string) (bool, error) {
	row := l.db.QueryRow("SELECT true FROM launch WHERE launchpad_id = $1 AND date = $2", launchpadId, date)
	return exists(row)
}

func (l *launchRepository) Add(launch *Launch) error {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	_, err = l.db.Exec("INSERT INTO launch (id, launchpad_id, date) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
		newUUID, launch.LaunchpadId, time.Time(launch.Date))
	return err
}

func (l *launchRepository) GetAllUpcoming() ([]Launch, error) {
	rows, err := l.db.Query("SELECT id, launchpad_id, date FROM launch ORDER BY date")
	if err != nil {
		return nil, err
	}
	return l.getLaunches(rows)
}

func (l *launchRepository) GetAtDate(date Date) ([]Launch, error) {
	rows, err := l.db.Query("SELECT id, launchpad_id, date FROM launch WHERE date = $1 ORDER BY id", date)
	if err != nil {
		return nil, err
	}
	return l.getLaunches(rows)
}

func (l *launchRepository) getLaunches(rows *sql.Rows) ([]Launch, error) {
	launches := make([]Launch, 0, 1)
	for rows.Next() {
		launch := Launch{}
		if err := rows.Scan(&launch.Id, &launch.LaunchpadId, &launch.Date); err != nil {
			return launches, err
		}
		launches = append(launches, launch)
	}
	return launches, nil
}

// exists is an utility function to check if record exists
func exists(row *sql.Row) (bool, error) {
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
