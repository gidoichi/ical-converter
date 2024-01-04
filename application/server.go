package application

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gidoichi/ical-converter/usecase"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	convertService convertService
	dataSource     usecase.DataSource
	port           string
}

func NewServer(convertService convertService, dataSource usecase.DataSource, port string) Server {
	return Server{
		convertService: convertService,
		dataSource:     dataSource,
		port:           port,
	}
}

func (s *Server) Run() {
	http.Handle("/", promhttp.InstrumentHandlerCounter(
		promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
		}, []string{"code"}),
		s,
	))

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":"+s.port, nil))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
