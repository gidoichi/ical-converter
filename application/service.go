package application

import (
	"log"
	"reflect"

	ical "github.com/arran4/golang-ical"
	"github.com/gidoichi/ical-converter/entity/component"
	"github.com/gidoichi/ical-converter/usecase"
)

type convertService struct {
	converter usecase.Converter
}

func NewConvertService(converter usecase.Converter) convertService {
	return convertService{
		converter: converter,
	}
}

func (s *convertService) Convert(dataSource usecase.DataSource) (string, error) {
	cal, err := s.converter.Convert(dataSource)
	if err != nil {
		return "", err
	}

	newCal := component.NewCalendarFrom(*cal)
	for _, event := range cal.Components {
		var vevent *ical.VEvent
		var ok bool
		if vevent, ok = event.(*ical.VEvent); !ok {
			log.Printf("unexpected component type: %s", reflect.TypeOf(event))
			continue
		}

		if vevent.GetProperty(ical.ComponentPropertyDtStart) == nil {
			continue
		}

		if status := vevent.GetProperty(ical.ComponentPropertyStatus); status != nil {
			switch ical.ObjectStatus(status.Value) {
			case ical.ObjectStatusTentative, ical.ObjectStatusConfirmed:
				vevent.ComponentBase = s.removeProperty(vevent.ComponentBase, ical.ComponentPropertyStatus)
			case ical.ObjectStatusCancelled, ical.ObjectStatusCompleted:
				continue
			}
		}

		newCal.AddVEvent(vevent)
	}

	return newCal.Serialize(), nil
}

func (s *convertService) removeProperty(component ical.ComponentBase, property ical.ComponentProperty) ical.ComponentBase {
	for i := 0; i < len(component.Properties); i++ {
		if ical.ComponentProperty(component.Properties[i].IANAToken) == property {
			component.Properties = append(component.Properties[:i], component.Properties[i+1:]...)
			i--
		}
	}

	return component
}
