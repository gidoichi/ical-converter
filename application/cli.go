package application

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/gidoichi/ical-converter/application/datasource"
	"github.com/gidoichi/ical-converter/usecase"
	"go.lsp.dev/uri"
)

type CLI struct {
	convertService      convertService
	icsURL              string
	icsBasicAuthUser    string
	icsBasicAuthPass    string
	scheme              string
}

func NewCLI(convertService convertService, icsURL, icsBasicAuthUser, icsBasicAuthPass string) (CLI, error) {
	location, err := url.Parse(icsURL)
	if err != nil {
		return CLI{}, fmt.Errorf("failed to parse url: %w", err)
	}
	if !uriSchemeCLISupported(location.Scheme) {
		return CLI{}, fmt.Errorf("unsupported scheme: %#v", location.Scheme)
	}

	return CLI{
		convertService:   convertService,
		icsURL:           icsURL,
		icsBasicAuthUser: icsBasicAuthUser,
		icsBasicAuthPass: icsBasicAuthPass,
		scheme:           location.Scheme,
	}, nil
}

func uriSchemeCLISupported(scheme string) bool {
	switch scheme {
	case uri.HTTPScheme, uri.HTTPSScheme, uri.FileScheme:
		return true
	default:
		return false
	}
}

func (s *CLI) Run() {
	var dataSource usecase.DataSource
	switch s.scheme {
	case uri.HTTPScheme, uri.HTTPSScheme:
		username, password := s.icsBasicAuthUser, s.icsBasicAuthPass
		dataSource = datasource.NewHTTPICalDataSource(s.icsURL, username, password)
	case uri.FileScheme:
		parsed, err := url.Parse(s.icsURL)
		if err != nil {
			log.Println("failed to parse url: ", err)
			os.Exit(1)
		}
		dataSource = datasource.NewFileICalDataSource(parsed.Path)
	}

	serialized, err := s.convertService.Convert(dataSource)
	if err != nil {
		log.Println("failed to convert: ", err)
		os.Exit(1)
	}

	fmt.Println(serialized)
}
