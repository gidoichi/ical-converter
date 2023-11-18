package converter

import (
	ical "github.com/arran4/golang-ical"

	"github.com/gidoichi/ical-converter/domain/component"
)

func Convert(todo component.Todo) *component.Event {
	id := todo.GetProperty(ical.ComponentPropertyUniqueId)
	event := component.NewEvent(id.Value)

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
