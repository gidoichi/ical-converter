package two_do

import (
	"fmt"
	"log"
	"time"

	ical "github.com/arran4/golang-ical"
	"github.com/gidoichi/ical-converter/entity/component"
	"github.com/gidoichi/ical-converter/entity/valuetype"
	"github.com/gidoichi/ical-converter/usecase"
)

type twoDoRepository struct {
	timeZone time.Location
}

func NewTwoDoRepository(timeZone time.Location) usecase.Repository {
	return &twoDoRepository{
		timeZone: timeZone,
	}
}

func (r *twoDoRepository) GetICal(source usecase.DataSource) (cal *ical.Calendar, err error) {
	rawCal, err := source.GetICal()
	if err != nil {
		return nil, err
	}

	cal = component.NewCalendarFrom(*rawCal)
	for _, rawTodo := range rawCal.Components {
		todo := component.Todo{
			ComponentBase: ical.ComponentBase{
				Components: rawTodo.SubComponents(),
				Properties: rawTodo.UnknownPropertiesIANAProperties(),
			},
		}

		if metadata, err := parseMetadata(todo, r.timeZone); err == nil && metadata != nil {
			if start, err := metadata.getStartTime(); err == nil && start != nil {
				if start.Hour() == 0 && start.Minute() == 0 && start.Second() == 0 {
					todo.SetDateProperty(ical.ComponentPropertyDtStart, valuetype.NewDate(*start))
				} else {
					todo.SetDateTimeProperty(ical.ComponentPropertyDtStart, valuetype.NewDateTime(*start))
				}
			}
			if url := metadata.getURL(); url != nil {
				todo.SetProperty(ical.ComponentPropertyUrl, *url)
			}
		} else if err != nil {
			log.Println(fmt.Errorf("failed to parse metadata: %w", err))
		}

		for _, targetProp := range []ical.Property{
			ical.PropertyDtstamp,
			ical.PropertyDtstart,
			ical.PropertyLastModified,
			ical.PropertyDue,
		} {
			prop := todo.GetProperty(ical.ComponentProperty(targetProp))
			if prop == nil {
				continue
			}

			if t, err := time.ParseInLocation("20060102T150405", prop.Value, &r.timeZone); err == nil {
				todo.SetDateTimeProperty(ical.ComponentProperty(targetProp), valuetype.NewDateTime(t))
			} else {
				var params []ical.PropertyParameter
				for k, v := range prop.ICalParameters {
					params = append(params, &ical.KeyValues{Key: k, Value: v})
				}
				todo.SetProperty(ical.ComponentProperty(targetProp), prop.Value, params...)
			}
		}
		cal.Components = append(cal.Components, (*ical.VTodo)(&todo))
	}

	return cal, nil
}
