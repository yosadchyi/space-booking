package booking

// Destination launchpad model
type Destination struct {
	Id   string
	Name string
}

// Launchpad launchpad model
type Launchpad struct {
	Id     string
	Name   string
	Status string
}

// Launch launch model
type Launch struct {
	Id          string
	LaunchpadId string
	Date        Date
}

// Booking booking model
type Booking struct {
	Id            string
	FirstName     string
	LastName      string
	Gender        Gender
	Birthday      Date
	LaunchpadId   string
	DestinationId string
	LaunchDate    Date
}
