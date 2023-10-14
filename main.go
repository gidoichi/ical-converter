package main

import (
	"errors"
	"fmt"
	ical "github.com/arran4/golang-ical"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", convert)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func convert(w http.ResponseWriter, r *http.Request) {
	icsURL := os.Getenv("ICAL_CONVERTER_ICS_URL")
	if icsURL == "" {
		log.Println(errors.New("failed to get env: ICAL_CONVERTER_ICS_URL"))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	resp, err := http.Get(icsURL)
	if err != nil {
		log.Println(fmt.Errorf("failed to get ical: %w", err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	cal, err := ical.ParseCalendar(resp.Body)
	if err != nil {
		log.Println(fmt.Errorf("failed to parse calendar", err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	for i, component := range cal.Components {
		var id string
		for _, prop := range component.UnknownPropertiesIANAProperties() {
			if prop.IANAToken == "UID" {
				id = prop.Value
				break
			}
		}

		event := ical.NewEvent(id)
		for _, prop := range component.UnknownPropertiesIANAProperties() {
			// var params []ical.PropertyParameter
			// for _, param := range prop.ICalParameters {
			// 	params = append(params, ical.KeyValues{})
			// }
			switch prop.IANAToken {
			case "SUMMARY":
				event.SetSummary(prop.Value)
			case "DTSTAMP":
				event.SetProperty(ical.ComponentPropertyDtstamp, prop.Value)
			case "DTSTART":
				event.SetProperty(ical.ComponentPropertyDtStart, prop.Value)
				// "DTSTAMP"
				// "LAST-MODIFIED"
				// "SEQUENCE"
				// "DESCRIPTION"
				// "STATUS"
				// "PERCENT-COMPLETE"
				// "COMPLETED"
				// "CREATED"
				// "DTSTART"
				// "DUE"
				// "X-2DOAPP-METADATA"
				// "X-2DOAPP-SYNC-IMAGE"
				// "X-2DOAPP-SYNC-AUDIO"
				// "X-2DOAPP-AUDIO-DURATION"
			}
		}

		// fmt.Println("==================================================")
		// fmt.Println(component)
		// fmt.Println("--------------------------------------------------")
		// fmt.Println(event.Serialize())
		cal.Components[i] = event
	}

	if _, err := fmt.Fprint(w, cal.Serialize()); err != nil {
		log.Println(fmt.Errorf("failed to write response: %v", err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
