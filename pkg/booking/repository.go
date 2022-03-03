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
	AddTx(tx *sql.Tx, booking *Booking) error
	GetDestinationIdForBookingId(tx *sql.Tx, id string) (string, error)
	GetAll() ([]Booking, error)
	Delete(tx *sql.Tx, id string) error
}

// DestinationRepository repository to access all launchpads
type DestinationRepository interface {
	Exists(id string) (bool, error)
	GetAll() ([]Destination, error)
}

// LaunchpadRepository repository to access all launchpads
type LaunchpadRepository interface {
	ExistsAndIsActive(id string) (bool, error)
	AddOrUpdate(launchpad *Launchpad) error
	GetAllActive() ([]Launchpad, error)
}

// LaunchRepository repository to access all launches
type LaunchRepository interface {
	Exists(date Date, launchpadId string) (bool, error)
	Add(launch *Launch) error
	AddTx(tx *sql.Tx, launch *Launch) error
	GetFromLaunchpadAtDate(launchpadId string, date Date) ([]Launch, error)
	GetWeekLaunches(tx *sql.Tx, launchpadId string, date Date) ([]Launch, error)
	GetAllUpcoming() ([]Launch, error)
	Delete(tx *sql.Tx, id string) error
}

type mainRepository struct {
	db *sql.DB
}

type destinationRepository struct {
	db *sql.DB
}

type launchpadRepository struct {
	db *sql.DB
}

type launchRepository struct {
	db *sql.DB
}

func NewMainRepository(db *sql.DB) MainRepository {
	return &mainRepository{db: db}
}

func (m *mainRepository) GetDestinationIdForBookingId(tx *sql.Tx, id string) (string, error) {
	row := tx.QueryRow("SELECT destination_id FROM booking WHERE id = $1", id)
	destinationId := ""
	err := row.Scan(&destinationId)
	if err != nil {
		return "", err
	}
	return destinationId, nil
}

func (m *mainRepository) AddTx(tx *sql.Tx, booking *Booking) error {
	query := `INSERT INTO booking 
    	(id, first_name, last_name, gender, birthday, launch_date, launchpad_id, destination_id, launch_id)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := tx.Exec(query, booking.Id, booking.FirstName, booking.LastName, booking.Gender, time.Time(booking.Birthday),
		time.Time(booking.LaunchDate), booking.LaunchpadId, booking.DestinationId, booking.LaunchId)
	return err
}

func (m *mainRepository) Delete(tx *sql.Tx, id string) error {
	_, err := tx.Exec("DELETE FROM booking WHERE id = $1", id)
	return err
}

func (m *mainRepository) GetAll() ([]Booking, error) {
	rows, err := m.db.Query(`SELECT 
			id, first_name, last_name, gender, birthday, launch_date, launchpad_id, destination_id, launch_id
		FROM booking ORDER BY id`)
	if err != nil {
		return nil, err
	}

	bookings := make([]Booking, 0, 1)
	for rows.Next() {
		booking := Booking{}
		if err := rows.Scan(&booking.Id, &booking.FirstName, &booking.LastName, &booking.Gender,
			&booking.Birthday, &booking.LaunchDate, &booking.LaunchpadId, &booking.DestinationId, &booking.LaunchId); err != nil {
			return bookings, err
		}
		bookings = append(bookings, booking)
	}
	return bookings, nil
}

// NewDestinationRepository creates new destination repository
func NewDestinationRepository(db *sql.DB) DestinationRepository {
	return &destinationRepository{db: db}
}

func (d *destinationRepository) Exists(id string) (bool, error) {
	row := d.db.QueryRow("SELECT true FROM destination WHERE id = $1", id)
	return exists(row)
}

func (d *destinationRepository) GetAll() ([]Destination, error) {
	rows, err := d.db.Query("SELECT id, name FROM destination ORDER BY name")
	if err != nil {
		return nil, err
	}

	destinations := make([]Destination, 0, 1)
	for rows.Next() {
		destination := Destination{}
		if err := rows.Scan(&destination.Id, &destination.Name); err != nil {
			return destinations, err
		}
		destinations = append(destinations, destination)
	}
	return destinations, nil
}

// NewLaunchpadRepository create new launchpad repository
func NewLaunchpadRepository(db *sql.DB) LaunchpadRepository {
	return &launchpadRepository{db: db}
}

func (l *launchpadRepository) ExistsAndIsActive(id string) (bool, error) {
	row := l.db.QueryRow("SELECT true FROM launchpad WHERE id = $1 AND status = 'active'", id)
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
	time := time.Time(launch.Date)
	year, week := time.ISOWeek()
	_, err = l.db.Exec("INSERT INTO launch (id, launchpad_id, date, year, week) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING",
		newUUID, launch.LaunchpadId, time, year, week)
	return err
}

func (l *launchRepository) AddTx(tx *sql.Tx, launch *Launch) error {
	time := time.Time(launch.Date)
	year, week := time.ISOWeek()
	_, err := tx.Exec("INSERT INTO launch (id, launchpad_id, date, year, week) VALUES ($1, $2, $3, $4, $5)",
		launch.Id, launch.LaunchpadId, time, year, week)
	return err
}

func (l *launchRepository) GetAllUpcoming() ([]Launch, error) {
	rows, err := l.db.Query("SELECT id, launchpad_id, date FROM launch ORDER BY date")
	if err != nil {
		return nil, err
	}
	return l.getLaunches(rows)
}

func (l *launchRepository) GetFromLaunchpadAtDate(launchpadId string, date Date) ([]Launch, error) {
	rows, err := l.db.Query("SELECT id, launchpad_id, date FROM launch WHERE launchpad_id = $1 AND date = $2 ORDER BY id",
		launchpadId, time.Time(date))
	if err != nil {
		return nil, err
	}
	return l.getLaunches(rows)
}

func (l *launchRepository) GetWeekLaunches(tx *sql.Tx, launchpadId string, date Date) ([]Launch, error) {
	year, week := time.Time(date).ISOWeek()
	rows, err := l.db.Query(`SELECT id, launchpad_id, date FROM launch 
		WHERE launchpad_id = $1 AND year = $2 AND week = $3 ORDER BY id`,
		launchpadId, year, week)
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

func (l *launchRepository) Delete(tx *sql.Tx, id string) error {
	_, err := tx.Exec("DELETE FROM launch WHERE id = $1", id)
	return err
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
