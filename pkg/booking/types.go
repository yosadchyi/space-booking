package booking

import (
	"encoding/json"
	"strings"
	"time"
)

// Gender represents gender of the passenger
type Gender string

const (
	MALE   = Gender("Male")
	FEMALE = Gender("Female")
)

// Date represents date without time
type Date time.Time

// UnmarshalJSON Implement Unmarshaler interface
func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

// MarshalJSON Implement Marshaler interface
func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d))
}

// Format Maybe a Format function for printing your date
func (d Date) Format(s string) string {
	t := time.Time(d)
	return t.Format(s)
}
