package booking

// MainRepository booking repository
type MainRepository interface {
	Add(booking *Booking) (bool, error)
	GetAll() ([]Booking, error)
}

// LaunchpadRepository repository to access all launchpads
type LaunchpadRepository interface {
	GetAll() ([]Launchpad, error)
}

// LaunchRepository repository to access all launches
type LaunchRepository interface {
	GetAtDate() ([]Launch, error)
}
