package application

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"cloudeng.io/net/http/httperror"
	"github.com/gidoichi/ical-converter/application/datasource"
	"github.com/gidoichi/ical-converter/usecase"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.lsp.dev/uri"
)

type Server struct {
	convertService convertService
	icsURL         string
	port           string
	scheme         string
}

func NewServer(convertService convertService, icsURL, port string) (Server, error) {
	location, err := url.Parse(icsURL)
	if err != nil {
		return Server{}, fmt.Errorf("failed to parse url: %w", err)
	}
	if !uriSchemeSupported(location.Scheme) {
		return Server{}, fmt.Errorf("unsupported scheme: %#v", location.Scheme)
	}

	return Server{
		convertService: convertService,
		icsURL:         icsURL,
		scheme:         location.Scheme,
		port:           port,
	}, nil
}

func uriSchemeSupported(scheme string) bool {
	switch scheme {
	case uri.HTTPScheme, uri.HTTPSScheme, uri.FileScheme:
		return true
	default:
		return false
	}
}

func (s *Server) Run() {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/", promhttp.InstrumentHandlerCounter(
		promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
		}, []string{"code"}),
		s,
	))
	server := &http.Server{
		Addr:    ":" + s.port,
		Handler: mux,
	}
	go func() {
		log.Println(server.ListenAndServe())
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("failed to shutdown server: %v", err)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%+v", r)

	var dataSource usecase.DataSource
	switch s.scheme {
	case uri.HTTPScheme, uri.HTTPSScheme:
		username, password, _ := r.BasicAuth()
		dataSource = datasource.NewHTTPICalDataSource(s.icsURL, username, password)
	case uri.FileScheme:
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
