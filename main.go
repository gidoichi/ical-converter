package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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

	service := newConvertService(icsURL)

	http.Handle("/", &service)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

type convertService struct {
	icsURL string
}

func newConvertService(icsURL string) convertService {
	return convertService{
		icsURL: icsURL,
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
	for _, component := range cal.Components {
		var id string
		for _, prop := range component.UnknownPropertiesIANAProperties() {
			if prop.IANAToken == string(ical.PropertyUid) {
				id = prop.Value
				break
			}
		}

		event := ical.NewEvent(id)
		for _, prop := range component.UnknownPropertiesIANAProperties() {
			var params []ical.PropertyParameter
			for k, v := range prop.ICalParameters {
				params = append(params, &ical.KeyValues{Key: k, Value: v})
			}

			switch prop.IANAToken {
			case string(ical.PropertyUid):
				continue
			case string(ical.PropertyPercentComplete):
				continue
			case string(ical.PropertyDue):
				event.SetProperty(ical.ComponentPropertyDtEnd, prop.Value, params...)
			default:
				event.SetProperty(ical.ComponentProperty(prop.IANAToken), prop.Value, params...)
			}
		}

		if event.GetProperty(ical.ComponentPropertyDtStart) == nil {
			continue
		}
		if event.GetProperty(ical.ComponentProperty(ical.PropertyCompleted)) != nil {
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
