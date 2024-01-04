package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	ical "github.com/arran4/golang-ical"
	"github.com/gidoichi/ical-converter/domain/component"
	"github.com/gidoichi/ical-converter/domain/converter"
	"github.com/gidoichi/ical-converter/domain/valuetype"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	icsURL := os.Getenv("ICAL_CONVERTER_ICS_URL")
	if icsURL == "" {
		log.Fatal("failed to get env: ICAL_CONVERTER_ICS_URL")
	}
	tz := time.FixedZone("JST", int((+9 * time.Hour).Seconds()))
	service := newConvertService(icsURL, *tz)

	http.Handle("/", promhttp.InstrumentHandlerCounter(
		promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
		}, []string{"code"}),
		&service,
	))

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

type convertService struct {
	icsURL   string
	timeZone time.Location
}

func newConvertService(icsURL string, timeZone time.Location) convertService {
	return convertService{
		icsURL:   icsURL,
		timeZone: timeZone,
	}
}

func (c *convertService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%+v", r)

	resp, err := http.Get(c.icsURL)
	if err != nil {
		log.Println("failed to get ical: ", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	cal, err := ical.ParseCalendar(resp.Body)
	if err != nil {
		log.Println("failed to parse calendar: ", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	newcal := ical.NewCalendar()
	newcal.CalendarProperties = cal.CalendarProperties
	for _, rawTodo := range cal.Components {
		// infrastructure logic

		todo := component.Todo{
			ComponentBase: ical.ComponentBase{
				Components: rawTodo.SubComponents(),
				Properties: rawTodo.UnknownPropertiesIANAProperties(),
			},
		}
		if start, err := c.getStartDateFrom2doappMetadata(&todo); start != nil && err == nil {
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

			if t, err := time.ParseInLocation("20060102T150405", prop.Value, &c.timeZone); err == nil {
				todo.SetDateTimeProperty(ical.ComponentProperty(targetProp), valuetype.NewDateTime(t))
			} else {
				var params []ical.PropertyParameter
				for k, v := range prop.ICalParameters {
					params = append(params, &ical.KeyValues{Key: k, Value: v})
				}
				todo.SetProperty(ical.ComponentProperty(targetProp), prop.Value, params...)
			}
		}

		// business logic

		event := converter.Convert(todo)

		// application logic

		if event.GetProperty(ical.ComponentPropertyDtStart) == nil {
			continue
		}
		status := event.GetProperty(ical.ComponentPropertyStatus).Value
		switch ical.ObjectStatus(status) {
		case ical.ObjectStatusTentative, ical.ObjectStatusConfirmed:
			event.RemoveProperty(ical.ComponentPropertyStatus)
		case ical.ObjectStatusCancelled, ical.ObjectStatusCompleted:
			continue
		}

		vevent := ical.VEvent(*event)
		newcal.AddVEvent(&vevent)
	}

	if _, err := fmt.Fprint(w, newcal.Serialize()); err != nil {
		log.Println("failed to write response: ", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (c *convertService) getStartDateFrom2doappMetadata(todo *component.Todo) (*time.Time, error) {
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
