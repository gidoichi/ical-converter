package component

import (
	ical "github.com/arran4/golang-ical"
)

type Event ical.VEvent

func NewEvent(id string) *Event {
	return (*Event)(ical.NewEvent(id))
}

func (e *Event) RemoveProperty(property ical.ComponentProperty) {
	for i := 0; i < len(e.Properties); i++ {
		if ical.ComponentProperty(e.Properties[i].IANAToken) == property {
			e.Properties = append(e.Properties[:i], e.Properties[i+1:]...)
			i--
		}
	}
}
