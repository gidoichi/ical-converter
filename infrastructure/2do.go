package infrastructure

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
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

		if metadata, err := parseMetadata(todo); err == nil {
			if start, err := metadata.getStartTime(); err == nil {
				start := start.UTC()
				if start.Hour() == 0 && start.Minute() == 0 && start.Second() == 0 {
					todo.SetDateProperty(ical.ComponentPropertyDtStart, valuetype.NewDate(start))
				} else {
					todo.SetDateTimeProperty(ical.ComponentPropertyDtStart, valuetype.NewDateTime(start))
				}
			}
			if url := metadata.getURL(); url != nil {
				todo.SetProperty(ical.ComponentPropertyUrl, *url)
			}
		} else {
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

type ActionType int

var (
	ActionTypeCall    ActionType = 5
	ActionTypeMessage ActionType = 6
	ActionTypeMail    ActionType = 7
	ActionTypeVisit   ActionType = 8
	ActionTypeURL     ActionType = 9
	ActionTypeSearch  ActionType = 10
)

type metadata struct {
	ActionType  ActionType `json:"actionType"`
	ActionValue string     `json:"actionValue"`
	StartDate   int64      `json:"StartDate"`
}

func parseMetadata(todo component.Todo) (*metadata, error) {
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

	var parsed metadata
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		return nil, fmt.Errorf("failed to unmarshall json: %w", err)
	}

	return &parsed, nil
}

func (m metadata) getURL() (url *string) {
	if m.ActionType != ActionTypeURL {
		return nil
	}

	return &m.ActionValue
}

func (m metadata) getStartTime() (*time.Time, error) {
	t := time.Unix(m.StartDate, 0)
	return &t, nil
}
