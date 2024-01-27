package application

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"cloudeng.io/net/http/httperror"
	"github.com/gidoichi/ical-converter/application/datasource"
	"github.com/gidoichi/ical-converter/usecase"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	convertService convertService
	icsURL         string
	port           string
	scheme         string
}

func NewServer(convertService convertService, icsURL, port string) Server {
	location, err := url.Parse(icsURL)
	if err != nil {
		log.Fatal("failed to parse ics url: ", err)
	}

	return Server{
		convertService: convertService,
		icsURL:         icsURL,
		scheme:         location.Scheme,
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

	var dataSource usecase.DataSource
	switch s.scheme {
	case "http", "https":
		username, password, _ := r.BasicAuth()
		dataSource = datasource.NewHTTPICalDataSource(s.icsURL, username, password)
	case "file":
		parsed, err := url.Parse(s.icsURL)
		if err != nil {
			log.Println("failed to parse url: ", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		dataSource = datasource.NewFileICalDataSource(parsed.Path)
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
