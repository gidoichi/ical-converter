//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE
package usecase

import (
	ical "github.com/arran4/golang-ical"
)

type DataSource interface {
	GetICal() (*ical.Calendar, error)
}

type Repository interface {
	GetICal(source DataSource) (*ical.Calendar, error)
}
