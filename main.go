package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
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
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal("failed to parse location: %w", err)
	}
	service := newConvertService(icsURL, *loc)

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
	resp, err := http.Get(c.icsURL)
	if err != nil {
		log.Println(fmt.Errorf("failed to get ical: %w", err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	cal, err := ical.ParseCalendar(resp.Body)
	if err != nil {
		log.Println(fmt.Errorf("failed to parse calendar: %w", err))
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

			switch prop.IANAToken {
			case string(ical.PropertyUid):
				continue
			case string(ical.PropertyCompleted):
				continue
			case string(ical.PropertyPercentComplete):
				continue
			case string(ical.PropertyStatus):
				switch prop.Value {
				case string(ical.ObjectStatusCancelled), string(ical.ObjectStatusCompleted):
					continue VTODO
				}
			case string(ical.PropertyDtstart):
				if t, err := time.ParseInLocation("20060102T150405", prop.Value, &c.timeZone); err == nil {
					event.SetProperty(ical.ComponentPropertyDtStart, t.UTC().Format("20060102T150405Z"), params...)
				} else {
					event.SetProperty(ical.ComponentPropertyDtStart, prop.Value, params...)
				}
			case string(ical.PropertyDue):
				if t, err := time.ParseInLocation("20060102T150405", prop.Value, &c.timeZone); err == nil {
					event.SetProperty(ical.ComponentPropertyDtEnd, t.UTC().Format("20060102T150405Z"), params...)
				} else {
					event.SetProperty(ical.ComponentPropertyDtEnd, prop.Value, params...)
				}
			default:
				event.SetProperty(ical.ComponentProperty(prop.IANAToken), prop.Value, params...)
			}
		}

		if event.GetProperty(ical.ComponentPropertyDtStart) == nil {
			continue
		}

		newcal.Components = append(newcal.Components, event)
	}

	if _, err := fmt.Fprint(w, newcal.Serialize()); err != nil {
		log.Println(fmt.Errorf("failed to write response: %v", err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
