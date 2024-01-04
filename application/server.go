package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	ical "github.com/arran4/golang-ical"
	"github.com/gidoichi/ical-converter/entity/component"
	"github.com/gidoichi/ical-converter/infrastructure"
	"github.com/gidoichi/ical-converter/usecase"
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
	repository := infrastructure.NewTwoDoRepository(*tz)
	dataSource := NewHTTPICalDataSource(icsURL)
	converter := usecase.NewConverter(repository)
	service := newConvertService(icsURL, dataSource, &converter)

	http.Handle("/", promhttp.InstrumentHandlerCounter(
		promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
		}, []string{"code"}),
		&service,
	))

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

type httpICalDataSource struct {
	url string
}

func NewHTTPICalDataSource(url string) httpICalDataSource {
	return httpICalDataSource{
		url: url,
	}
}

func (d httpICalDataSource) GetICal() (*ical.Calendar, error) {
	resp, err := http.Get(d.url)
	if err != nil {
		return nil, fmt.Errorf("failed to get ical: %w", err)
	}

	cal, err := ical.ParseCalendar(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse calendar: %w", err)
	}

	return cal, nil
}

type convertService struct {
	icsURL     string
	dataSource usecase.DataSource
	converter  usecase.Converter
}

func newConvertService(icsURL string, dataSource usecase.DataSource, converter usecase.Converter) convertService {
	return convertService{
		icsURL:     icsURL,
		dataSource: dataSource,
		converter:  converter,
	}
}

func (c *convertService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%+v", r)

	cal, err := c.converter.Convert(c.dataSource)
	if err != nil {
		log.Println("failed to convert: ", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	newCal := component.NewCalendarFrom(*cal)
	for _, event := range cal.Components {
		var vevent *ical.VEvent
		var ok bool
		if vevent, ok = event.(*ical.VEvent); !ok {
			log.Printf("component type not supported: %s", reflect.TypeOf(event))
			continue
		}

		if vevent.GetProperty(ical.ComponentPropertyDtStart) == nil {
			continue
		}

		status := vevent.GetProperty(ical.ComponentPropertyStatus).Value
		switch ical.ObjectStatus(status) {
		case ical.ObjectStatusTentative, ical.ObjectStatusConfirmed:
			vevent.ComponentBase = c.RemoveProperty(vevent.ComponentBase, ical.ComponentPropertyStatus)
		case ical.ObjectStatusCancelled, ical.ObjectStatusCompleted:
			continue
		}

		newCal.AddVEvent(vevent)
	}

	if _, err := fmt.Fprint(w, newCal.Serialize()); err != nil {
		log.Println("failed to write response: ", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (s *convertService) RemoveProperty(component ical.ComponentBase, property ical.ComponentProperty) ical.ComponentBase {
	for i := 0; i < len(component.Properties); i++ {
		if ical.ComponentProperty(component.Properties[i].IANAToken) == property {
			component.Properties = append(component.Properties[:i], component.Properties[i+1:]...)
			i--
		}
	}

	return component
}
