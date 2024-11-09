package component

import (
	ical "github.com/arran4/golang-ical"
	"github.com/gidoichi/ical-converter/entity/valuetype"
)

type Todo struct {
	ical.VTodo
}

func (e *Todo) SetDateProperty(name ical.ComponentProperty, value valuetype.Date) {
	e.SetProperty(name, value.String(), &ical.KeyValues{
		Key:   "VALUE",
		Value: []string{"DATE"},
	})
}

func (e *Todo) SetDateTimeProperty(name ical.ComponentProperty, value valuetype.DateTime) {
	e.SetProperty(name, value.String())
}

func (e *Todo) RemoveProperty(name ical.ComponentProperty) {
	for i, _ := range e.Properties {
		if ical.ComponentProperty(e.Properties[i].IANAToken) == name {
			e.Properties = append(e.Properties[:i], e.Properties[i+1:]...)
			break
		}
	}
}
