package valuetype

import (
	"fmt"
	"time"
)

type DateTime time.Time

var _ fmt.Stringer = Date{}

func NewDateTime(t time.Time) DateTime {
	return DateTime(t)
}

func (t DateTime) String() string {
	return time.Time(t).UTC().Format("20060102T150405Z")
}
