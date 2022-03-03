package booking

// Request represents request for booking the flight
type Request struct {
	FirstName     string
	LastName      string
	Gender        Gender
	Birthday      Date
	LaunchpadId   string
	DestinationId string
	LaunchDate    Date
}

// SuccessResponse response in case of successful booking
type SuccessResponse struct {
	Id string
}

// ErrorResponse response in case of booking error
type ErrorResponse struct {
	Code    string
	Message string
}

// AllLaunchpadsResponse represents response to /launchpad/ GET request
type AllLaunchpadsResponse []Launchpad

// AllDestinationsResponse represents response to /destination/ GET request
type AllDestinationsResponse []Destination

// AllBookingsResponse represents response to /booking/ GET request
type AllBookingsResponse []Booking
