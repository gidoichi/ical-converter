package main

import (
	"log"
	"os"
	"time"

	"github.com/gidoichi/ical-converter/application"
	"github.com/gidoichi/ical-converter/application/datasource"
	"github.com/gidoichi/ical-converter/infrastructure"
	"github.com/gidoichi/ical-converter/usecase"
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
	convertService := application.NewConvertService(&converter)
	server := application.NewServer(convertService, dataSource, port)

	server.Run()
}
