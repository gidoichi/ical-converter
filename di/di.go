package di

import (
	"log"
	"os"
	"time"

	"github.com/gidoichi/ical-converter/application"
	twoDo "github.com/gidoichi/ical-converter/infrastructure/two_do"
	"github.com/gidoichi/ical-converter/usecase"
)

func DI() *application.Server {
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	icsURL := os.Getenv("ICAL_CONVERTER_ICS_URL")
	if icsURL == "" {
		log.Fatal("failed to get env: ICAL_CONVERTER_ICS_URL")
	}

	tz := time.FixedZone("JST", int((+9 * time.Hour).Seconds()))
	repository := twoDo.NewTwoDoRepository(*tz)
	converter := usecase.NewConverter(repository)
	convertService := application.NewConvertService(&converter)
	server := application.NewServer(convertService, icsURL, port)

	return &server
}
