package datasource

import (
	"fmt"
	"net/http"

	"cloudeng.io/net/http/httperror"
	ical "github.com/arran4/golang-ical"
)

type httpICalDataSource struct {
	url               string
	basicAuthUsername string
	basicAuthPassword string
}

func NewHTTPICalDataSource(url string) httpICalDataSource {
	return httpICalDataSource{
		url: url,
	}
}

func (d *httpICalDataSource) SetBasicAuth(username, password string) {
	d.basicAuthUsername = username
	d.basicAuthPassword = password
}

func (d httpICalDataSource) GetICal() (*ical.Calendar, error) {
	var resp *http.Response

	if d.basicAuthUsername != "" {
		client := &http.Client{}
		req, err := http.NewRequest("GET", d.url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
		req.SetBasicAuth(d.basicAuthUsername, d.basicAuthPassword)
		if resp, err = client.Do(req); err != nil {
			return nil, fmt.Errorf("failed to get ical: %w", err)
		}
	} else {
		var err error
		if resp, err = http.Get(d.url); err != nil {
			return nil, fmt.Errorf("failed to get ical: %w", err)
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, &httperror.T{
			Status:     "failed to get ical: unauthorized",
			StatusCode: resp.StatusCode,
		}
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get ical: status code %d", resp.StatusCode)
	}
	cal, err := ical.ParseCalendar(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse calendar: %w", err)
	}

	return cal, nil
}
