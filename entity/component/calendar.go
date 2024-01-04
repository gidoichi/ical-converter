package component

import (
	ical "github.com/arran4/golang-ical"
)

func NewCalendarFrom(base ical.Calendar) *ical.Calendar {
	cal := ical.NewCalendar()
	cal.CalendarProperties = base.CalendarProperties
	return cal
}
