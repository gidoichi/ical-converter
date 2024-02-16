package main_test

import (
	"fmt"
	"log"
	"strings"

	ical "github.com/arran4/golang-ical"
	"github.com/gidoichi/ical-converter/usecase"
)

type stringICalDataSource struct {
	data string
}

func (d stringICalDataSource) GetICal() (*ical.Calendar, error) {
	return ical.ParseCalendar(strings.NewReader(d.data))
}

type simpleRepository struct {
}

func (r *simpleRepository) GetICal(source usecase.DataSource) (*ical.Calendar, error) {
	return source.GetICal()
}

func Example() {
	source := stringICalDataSource{
		data: `
BEGIN:VCALENDAR
PRODID:-//xyz Corp//NONSGML PDA Calendar Version 1.0//EN
VERSION:2.0
BEGIN:VTODO
DTSTAMP:19960704T120000Z
UID:uid1@example.com
DTSTART:19960918T143000Z
DUE:19960920T220000Z
SUMMARY:Networld+Interop Conference
END:VTODO
END:VCALENDAR`,
	}

	repo := simpleRepository{}
	converter := usecase.NewConverter(&repo)
	cal, err := converter.Convert(source)
	if err != nil {
		log.Fatal(err)
	}

	// Serialized calendar has CRLF. To remove CR for this example, strings.Replace is used.
	fmt.Println(strings.Replace(cal.Serialize(), "\r", "", -1))
	// Output:
	// BEGIN:VCALENDAR
	// PRODID:-//xyz Corp//NONSGML PDA Calendar Version 1.0//EN
	// VERSION:2.0
	// BEGIN:VEVENT
	// UID:uid1@example.com
	// DTSTAMP:19960704T120000Z
	// DTSTART:19960918T143000Z
	// DTEND:19960920T220000Z
	// SUMMARY:Networld+Interop Conference
	// END:VEVENT
	// END:VCALENDAR
}
