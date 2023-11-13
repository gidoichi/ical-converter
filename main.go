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

	http.Handle("/", &service)
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
VTODO:
	for _, todo := range cal.Components {
		var id string
		for _, prop := range todo.UnknownPropertiesIANAProperties() {
			if prop.IANAToken == string(ical.PropertyUid) {
				id = prop.Value
				break
			}
		}

		event := ical.NewEvent(id)
		for _, prop := range todo.UnknownPropertiesIANAProperties() {
			var params []ical.PropertyParameter
			for k, v := range prop.ICalParameters {
				params = append(params, &ical.KeyValues{Key: k, Value: v})
			}

			switch ical.Property(prop.IANAToken) {
			case ical.PropertyUid:
				continue
			case ical.PropertyCompleted:
				continue
			case ical.PropertyPercentComplete:
				continue
			case ical.PropertyStatus:
				switch ical.ObjectStatus(prop.Value) {
				case ical.ObjectStatusCancelled, ical.ObjectStatusCompleted:
					continue VTODO
				}
			case ical.PropertyDtstamp, ical.PropertyDtstart, ical.PropertyLastModified:
				if t, err := time.ParseInLocation("20060102T150405", prop.Value, &c.timeZone); err == nil {
					event.SetProperty(ical.ComponentProperty(prop.IANAToken), t.UTC().Format("20060102T150405Z"), params...)
				} else {
					event.SetProperty(ical.ComponentProperty(prop.IANAToken), prop.Value, params...)
				}
			case ical.PropertyDue:
				if t, err := time.ParseInLocation("20060102T150405", prop.Value, &c.timeZone); err == nil {
					event.SetProperty(ical.ComponentPropertyDtEnd, t.UTC().Format("20060102T150405Z"), params...)
				} else {
					event.SetProperty(ical.ComponentPropertyDtEnd, prop.Value, params...)
				}
			default:
				event.SetProperty(ical.ComponentProperty(prop.IANAToken), prop.Value, params...)
			}
		}

		if start, err := c.getStartDateFrom2doappMetadata(event); start != nil && err == nil {
			if start.Hour() == 0 && start.Minute() == 0 && start.Second() == 0 {
				event.SetProperty(ical.ComponentPropertyDtStart, start.Format("20060102"), &ical.KeyValues{
					Key:   "VALUE",
					Value: []string{"DATE"},
				})
			} else {
				event.SetProperty(ical.ComponentPropertyDtStart, start.UTC().Format("20060102T150405Z"))
			}
		} else if err != nil {
			log.Println(err)
		}
		if event.GetProperty(ical.ComponentPropertyDtStart) == nil {
			continue
		}

		newcal.Components = append(newcal.Components, event)
	}

	if _, err := fmt.Fprint(w, newcal.Serialize()); err != nil {
		log.Println("failed to write response: ", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (c *convertService) getStartDateFrom2doappMetadata(event *ical.VEvent) (*time.Time, error) {
	prop := event.GetProperty("X-2DOAPP-METADATA")
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

	t := time.Unix(parsed.StartDate, 0).In(&c.timeZone)
	_, offset := t.Zone()
	t = t.Add(-time.Second * time.Duration(offset))
	return &t, nil
}
