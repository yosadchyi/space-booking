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

// AllLaunchpadsResponse represents response to /launchpad/ request
type AllLaunchpadsResponse []Launchpad

// AllDestinationsResponse represents response to /destination/ request
type AllDestinationsResponse []Destination
