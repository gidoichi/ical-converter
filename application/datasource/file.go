package datasource

import (
	"fmt"
	"os"

	ical "github.com/arran4/golang-ical"
	"github.com/gidoichi/ical-converter/usecase"
)

type fileICalDataSource struct {
	path string
}

func NewFileICalDataSource(path string) usecase.DataSource {
	return fileICalDataSource{
		path: path,
	}
}

func (d fileICalDataSource) GetICal() (*ical.Calendar, error) {
	file, err := os.Open(d.path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	cal, err := ical.ParseCalendar(file)
	if err != nil {
		return nil, fmt.Errorf("parse calendar: %w", err)
	}

	return cal, nil
}
