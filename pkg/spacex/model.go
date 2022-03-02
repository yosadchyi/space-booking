package spacex

// Launchpad represents launchpad response item, only needed fields are defined
type Launchpad struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// Launch is an information about launch (e.g. upcoming), only needed fields are defined
type Launch struct {
	Id        string `json:"id"`
	Launchpad string `json:"launchpad"`
	DateUnix  int64  `json:"date_unix"`
}
