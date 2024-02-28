package component_test

import (
	"reflect"
	"testing"
	"time"

	ical "github.com/arran4/golang-ical"
	"github.com/gidoichi/ical-converter/entity/component"
	"github.com/gidoichi/ical-converter/entity/valuetype"
)

func TestTodoCanSetDateProperty(t *testing.T) {
	// Given
	var todo component.Todo

	// When
	date := valuetype.NewDate(time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC))
	todo.SetDateProperty(ical.ComponentPropertyDtstamp, date)

	// Then
	if !reflect.DeepEqual(todo.GetProperty(ical.ComponentPropertyDtstamp).ICalParameters["VALUE"], []string{"DATE"}) {
		t.Errorf("e.GetProperty(ical.ComponentPropertyDtstamp).ICalParameters[VALUE] = %s; want [DATE]", todo.GetProperty(ical.ComponentPropertyDtstamp).ICalParameters["VALUE"])
	}
	if todo.GetProperty(ical.ComponentPropertyDtstamp).Value != "20060102" {
		t.Errorf("e.GetProperty(ical.ComponentPropertyDtstamp).Value = %s; want 20060102", todo.GetProperty(ical.ComponentPropertyDtstamp).Value)
	}
}

func TestTodoCanSetDateTimeProperty(t *testing.T) {
	// Given
	var todo component.Todo

	// When
	datetime := valuetype.NewDateTime(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC))
	todo.SetDateTimeProperty(ical.ComponentPropertyDtstamp, datetime)

	// Then
	if len(todo.GetProperty(ical.ComponentPropertyDtstamp).ICalParameters) != 0 {
		t.Errorf("e.GetProperty(ical.ComponentPropertyDtstamp).ICalParameters = %s; want empty", todo.GetProperty(ical.ComponentPropertyDtstamp).ICalParameters)
	}
	if todo.GetProperty(ical.ComponentPropertyDtstamp).Value != "20060102T150405Z" {
		t.Errorf("e.GetProperty(ical.ComponentPropertyDtstamp).Value = %s; want 20060102T150405Z", todo.GetProperty(ical.ComponentPropertyDtstamp).Value)
	}
}

func TestTodoCanRemoveProperty(t *testing.T) {
	// Given
	todo := component.Todo(
		ical.VTodo{ComponentBase: ical.ComponentBase{Properties: []ical.IANAProperty{
			{BaseProperty: ical.BaseProperty{
				IANAToken: "DTSTART",
				Value:     "20060102T150405Z",
			}},
		}}},
	)

	// When
	todo.RemoveProperty(ical.ComponentPropertyDtStart)

	// Then
	if len(todo.Properties) != 0 {
		t.Errorf("len(e.Properties) = %d; want 0", len(todo.Properties))
	}
}
