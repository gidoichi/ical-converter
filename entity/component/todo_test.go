package component_test

import (
	"log"
	"reflect"
	"testing"
	"time"

	ical "github.com/arran4/golang-ical"
	"github.com/gidoichi/ical-converter/entity/component"
	"github.com/gidoichi/ical-converter/entity/valuetype"
)

func TestSetDateProperty(t *testing.T) {
	// Given
	var todo component.Todo

	// When
	date := valuetype.NewDate(time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC))
	todo.SetDateProperty(ical.ComponentPropertyDtstamp, date)

	// Then
	log.Println("Set property is type of DATE")
	if !reflect.DeepEqual(todo.GetProperty(ical.ComponentPropertyDtstamp).ICalParameters["VALUE"], []string{"DATE"}) {
		t.Errorf("e.GetProperty(ical.ComponentPropertyDtstamp).ICalParameters[VALUE] = %s; want [DATE]", todo.GetProperty(ical.ComponentPropertyDtstamp).ICalParameters["VALUE"])
	}
	if todo.GetProperty(ical.ComponentPropertyDtstamp).Value != "20060102" {
		t.Errorf("e.GetProperty(ical.ComponentPropertyDtstamp).Value = %s; want 20060102", todo.GetProperty(ical.ComponentPropertyDtstamp).Value)
	}
}

func TestSetDateTimeProperty(t *testing.T) {
	// Given
	var todo component.Todo

	// When
	datetime := valuetype.NewDateTime(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC))
	todo.SetDateTimeProperty(ical.ComponentPropertyDtstamp, datetime)

	// Then
	log.Println("Set property is type of DATE-TIME")
	if len(todo.GetProperty(ical.ComponentPropertyDtstamp).ICalParameters) != 0 {
		t.Errorf("e.GetProperty(ical.ComponentPropertyDtstamp).ICalParameters = %s; want empty", todo.GetProperty(ical.ComponentPropertyDtstamp).ICalParameters)
	}
	if todo.GetProperty(ical.ComponentPropertyDtstamp).Value != "20060102T150405Z" {
		t.Errorf("e.GetProperty(ical.ComponentPropertyDtstamp).Value = %s; want 20060102T150405Z", todo.GetProperty(ical.ComponentPropertyDtstamp).Value)
	}
}
