package valuetype

import (
	"fmt"
	"time"
)

type Date time.Time

var _ fmt.Stringer = Date{}

func NewDate(t time.Time) Date {
	return Date(t)
}

func (t Date) String() string {
	return time.Time(t).Format("20060102")
}
