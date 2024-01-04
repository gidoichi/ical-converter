package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gidoichi/ical-converter/application/datasource"
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
	dataSource := datasource.NewHTTPICalDataSource(icsURL)
	converter := usecase.NewConverter(repository)
	convertService := NewConvertService(&converter)
	server := newServer(convertService, dataSource)

	http.Handle("/", promhttp.InstrumentHandlerCounter(
		promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
		}, []string{"code"}),
		&server,
	))

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

type server struct {
	convertService convertService
	dataSource     usecase.DataSource
}

func newServer(convertService convertService, dataSource usecase.DataSource) server {
	return server{
		convertService: convertService,
		dataSource:     dataSource,
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%+v", r)

	serialized, err := s.convertService.Convert(s.dataSource)
	if err != nil {
		log.Println("failed to convert: ", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if _, err := fmt.Fprint(w, serialized); err != nil {
		log.Println("failed to write response: ", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
