package component

import (
	ical "github.com/arran4/golang-ical"
	"github.com/gidoichi/ical-converter/domain/valuetype"
)

type Todo ical.VTodo

func (e *Todo) SetDateProperty(name ical.ComponentProperty, value valuetype.Date) {
	e.SetProperty(name, value.String(), &ical.KeyValues{
		Key:   "VALUE",
		Value: []string{"DATE"},
	})
}

func (e *Todo) SetDateTimeProperty(name ical.ComponentProperty, value valuetype.DateTime) {
	e.SetProperty(name, value.String())
}
