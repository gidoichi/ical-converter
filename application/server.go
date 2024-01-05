package application

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"cloudeng.io/net/http/httperror"
	"github.com/gidoichi/ical-converter/application/datasource"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	convertService convertService
	icsURL         string
	port           string
}

func NewServer(convertService convertService, icsURL, port string) Server {
	return Server{
		convertService: convertService,
		icsURL:         icsURL,
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

	dataSource := datasource.NewHTTPICalDataSource(s.icsURL)
	if username, password, ok := r.BasicAuth(); ok {
		dataSource.SetBasicAuth(username, password)
	}

	serialized, err := s.convertService.Convert(dataSource)
	if err != nil {
		var httpErr *httperror.T
		if ok := errors.As(err, &httpErr); ok {
			if httperror.IsHTTPError(httpErr, http.StatusUnauthorized) {
				w.Header().Set("WWW-Authenticate", "Basic")
				http.Error(w, "", httpErr.StatusCode)
				return
			} else if httperror.IsHTTPError(httpErr, http.StatusForbidden) {
				http.Error(w, "", httpErr.StatusCode)
				return
			}
		}

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
