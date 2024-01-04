package infrastructure

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	ical "github.com/arran4/golang-ical"
	dcomponent "github.com/gidoichi/ical-converter/domain/component"
	"github.com/gidoichi/ical-converter/domain/valuetype"
	"github.com/gidoichi/ical-converter/entity/component"
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
		todo := dcomponent.Todo{
			ComponentBase: ical.ComponentBase{
				Components: rawTodo.SubComponents(),
				Properties: rawTodo.UnknownPropertiesIANAProperties(),
			},
		}
		if start, err := r.getStartDateFrom2doappMetadata(&todo); start != nil && err == nil {
			start := start.UTC()
			if start.Hour() == 0 && start.Minute() == 0 && start.Second() == 0 {
				todo.SetDateProperty(ical.ComponentPropertyDtStart, valuetype.NewDate(start))
			} else {
				todo.SetDateTimeProperty(ical.ComponentPropertyDtStart, valuetype.NewDateTime(start))
			}
		} else if err != nil {
			log.Println(err)
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

func (c *twoDoRepository) getStartDateFrom2doappMetadata(todo *dcomponent.Todo) (*time.Time, error) {
	prop := todo.GetProperty("X-2DOAPP-METADATA")
	if prop == nil {
		return nil, nil
	}

	raw := prop.Value
	content := strings.TrimSuffix(strings.TrimPrefix(raw, "<2Do Meta>"), "</2Do Meta>\\n")
	content, err := url.QueryUnescape(content)
	if err != nil {
		return nil, fmt.Errorf("failed to unescape percent-encoding: %w", err)
	}

	var parsed struct {
		StartDate int64
	}
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		return nil, fmt.Errorf("failed to unmarshall json: %w", err)
	}
	if parsed.StartDate == 0 {
		return nil, nil
	}

	t := time.Unix(parsed.StartDate, 0)
	return &t, nil
}
