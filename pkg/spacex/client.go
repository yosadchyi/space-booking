package spacex

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Client interface {
	GetAllLaunchpads() ([]Launchpad, error)
	GetUpcomingLaunches() ([]Launch, error)
}

type client struct {
	baseUrl string
}

// NewClient creates new SpaceX API client.
func NewClient() Client {
	return &client{
		baseUrl: "https://api.spacexdata.com/v4",
	}
}

func (c *client) GetAllLaunchpads() ([]Launchpad, error) {
	var result []Launchpad
	resource := "/launchpads"

	err := c.doGet(resource, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *client) GetUpcomingLaunches() ([]Launch, error) {
	var result []Launch
	resource := "/launches/upcoming"

	err := c.doGet(resource, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *client) doGet(resource string, result interface{}) error {
	url := fmt.Sprintf("%s%s", c.baseUrl, resource)

	resp, err := http.Get(url)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			println("error closing response")
		}
	}(resp.Body)

	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, result)
	if err != nil {
		return err
	}

	return nil
}
