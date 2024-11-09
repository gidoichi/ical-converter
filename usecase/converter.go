//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE
package usecase

import (
	"fmt"
	"log"
	"reflect"

	ical "github.com/arran4/golang-ical"
	"github.com/gidoichi/ical-converter/entity/component"
	eerror "github.com/gidoichi/ical-converter/entity/error"
)

type Converter interface {
	Convert(source DataSource) (converted *ical.Calendar, err error)
}

type converter struct {
	repository Repository
}

func NewConverter(repository Repository) converter {
	return converter{
		repository: repository,
	}
}

func (c *converter) Convert(source DataSource) (converted *ical.Calendar, err error) {
	cal, err := c.repository.GetICal(source)
	switch err.(type) {
	case nil:
	case *eerror.ComponentsError:
		log.Printf("failed to get ical: %v", err)
	default:
		return nil, err
	}

	converted = component.NewCalendarFrom(*cal)
	for _, comp := range cal.Components {
		var event *ical.VEvent

		switch v := comp.(type) {
		case *ical.VTodo:
			event = c.convertFromTodo(*v)
		case *ical.VTimezone:
			converted.AddVTimezone(v)
			continue
		default:
			return nil, fmt.Errorf("component type not supported: %s", reflect.TypeOf(comp))
		}

		converted.AddVEvent(event)
	}

	return converted, nil
}

func (c *converter) convertFromTodo(todo ical.VTodo) *ical.VEvent {
	id := todo.GetProperty(ical.ComponentPropertyUniqueId)
	event := ical.NewEvent(id.Value)

	for _, prop := range todo.UnknownPropertiesIANAProperties() {
		var params []ical.PropertyParameter
		for k, v := range prop.ICalParameters {
			params = append(params, &ical.KeyValues{Key: k, Value: v})
		}

		switch ical.Property(prop.IANAToken) {
		case ical.PropertyUid:
			continue
		case ical.PropertyCompleted, ical.PropertyPercentComplete:
			continue
		case ical.PropertyStatus:
			var status ical.ObjectStatus
			switch ical.ObjectStatus(prop.Value) {
			case ical.ObjectStatusCompleted, ical.ObjectStatusCancelled:
				status = ical.ObjectStatusCancelled
			case ical.ObjectStatusNeedsAction, ical.ObjectStatusInProcess:
				status = ical.ObjectStatusConfirmed
			}
			event.SetProperty(ical.ComponentPropertyStatus, string(status), params...)
		case ical.PropertyDue:
			event.SetProperty(ical.ComponentPropertyDtEnd, prop.Value, params...)
		default:
			event.SetProperty(ical.ComponentProperty(prop.IANAToken), prop.Value, params...)
		}
	}

	return event
}
