package datasource

import (
	"fmt"
	"net/http"

	ical "github.com/arran4/golang-ical"
)

type httpICalDataSource struct {
	url string
}

func NewHTTPICalDataSource(url string) httpICalDataSource {
	return httpICalDataSource{
		url: url,
	}
}

func (d httpICalDataSource) GetICal() (*ical.Calendar, error) {
	resp, err := http.Get(d.url)
	if err != nil {
		return nil, fmt.Errorf("failed to get ical: %w", err)
	}

	cal, err := ical.ParseCalendar(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse calendar: %w", err)
	}

	return cal, nil
}
