package component_test

import (
	"testing"

	ical "github.com/arran4/golang-ical"
	"github.com/gidoichi/ical-converter/entity/component"
)

func TestNewCalendarFromReturnsCopiedCalendarWithPropertiesAndWithoutComponents(t *testing.T) {
	// Given

	base := ical.NewCalendar()
	base.CalendarProperties = []ical.CalendarProperty{
		{BaseProperty: ical.BaseProperty{IANAToken: "calendar"}},
	}
	base.Components = []ical.Component{
		&ical.VEvent{
			ComponentBase: ical.ComponentBase{
				Properties: []ical.IANAProperty{
					{BaseProperty: ical.BaseProperty{IANAToken: "component"}},
				},
			},
		},
	}

	// When

	cal := component.NewCalendarFrom(*base)

	// Then

	if len(cal.CalendarProperties) != 1 {
		t.Errorf("len(cal.CalendarProperties) = %d; want 1", len(cal.CalendarProperties))
	}
	if cal.CalendarProperties[0].IANAToken != "calendar" {
		t.Errorf("cal.CalendarProperties[0].IANAToken = %s; want token", cal.CalendarProperties[0].IANAToken)
	}

	if len(cal.Components) != 0 {
		t.Errorf("len(cal.Components) = %d; want 0", len(cal.Components))
	}
}
