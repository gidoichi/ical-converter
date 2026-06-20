package main

import (
	"flag"
	"log"
	"time"

	"github.com/gidoichi/ical-converter/application"
	twoDo "github.com/gidoichi/ical-converter/infrastructure/two_do"
	"github.com/gidoichi/ical-converter/usecase"
)

func di() application.CLI{
	icsURL := flag.String("ics-url", "", "Remote ical file server. Supported schemes are http, https, file.")
	icsBasicAuthUser := flag.String("basic-auth-user", "", "Username for basic auth when fetching ical file from remote server.")
	icsBasicAuthPassword := flag.String("basic-auth-password", "", "Password for basic auth when fetching ical file from remote server.")

	flag.Parse()

	if *icsURL == "" {
		log.Fatal("failed to get env: ICAL_CONVERTER_ICS_URL")
	}

	tz := time.FixedZone("JST", int((+9 * time.Hour).Seconds()))
	repository := twoDo.NewTwoDoRepository(*tz)
	converter := usecase.NewConverter(repository)
	convertService := application.NewConvertService(&converter)
	cli, err := application.NewCLI(convertService, *icsURL, *icsBasicAuthUser, *icsBasicAuthPassword)
	if err != nil {
		log.Fatal(err)
	}

	return cli
}
