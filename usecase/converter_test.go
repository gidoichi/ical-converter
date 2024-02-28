package usecase_test

import (
	"reflect"
	"testing"

	ical "github.com/arran4/golang-ical"
	"github.com/gidoichi/ical-converter/usecase"
	"github.com/gidoichi/ical-converter/usecase/mock_usecase"
	"go.uber.org/mock/gomock"
)

func TestConverterCanConvertVtodoToVevent(t *testing.T) {
	// Given

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repository := mock_usecase.NewMockRepository(ctrl)
	given := ical.Calendar{
		Components: []ical.Component{
			&ical.VTodo{
				ComponentBase: ical.ComponentBase{
					Properties: []ical.IANAProperty{
						{BaseProperty: ical.BaseProperty{IANAToken: "UID", Value: "1"}},
						{BaseProperty: ical.BaseProperty{IANAToken: "SUMMARY", Value: "test"}},
						{BaseProperty: ical.BaseProperty{IANAToken: "COMPLETED", Value: "20070707T100000Z"}},
						{BaseProperty: ical.BaseProperty{IANAToken: "STATUS", Value: "COMPLETED"}},
						{BaseProperty: ical.BaseProperty{IANAToken: "DUE", Value: "20070707T100000Z"}},
					},
				},
			},
		},
	}
	dataSource := mock_usecase.NewMockDataSource(ctrl)
	repository.EXPECT().GetICal(dataSource).Return(&given, nil)
	converter := usecase.NewConverter(repository)

	// When

	cal, err := converter.Convert(dataSource)

	// Then

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if len(cal.Components) != 1 {
		t.Errorf("unexpected calendar: %s", cal.Serialize())
	}
	if _, ok := cal.Components[0].(*ical.VEvent); !ok {
		t.Errorf("unexpected component type: %s", reflect.TypeOf(cal.Components[0]))
	}

	if len(cal.Components[0].UnknownPropertiesIANAProperties()) != 4 {
		t.Errorf("unexpected properties: %s", cal.Serialize())
	}
	for _, comp := range cal.Components[0].UnknownPropertiesIANAProperties() {
		switch comp.IANAToken {
		case "UID":
			if comp.Value != "1" {
				t.Errorf("unexpected UID: %s", comp.Value)
			}
		case "SUMMARY":
			if comp.Value != "test" {
				t.Errorf("unexpected SUMMARY: %s", comp.Value)
			}
		case "STATUS":
			if comp.Value != "CANCELLED" {
				t.Errorf("unexpected STATUS: %s", comp.Value)
			}
		case "DTEND":
			if comp.Value != "20070707T100000Z" {
				t.Errorf("unexpected DTEND: %s", comp.Value)
			}
		default:
			t.Errorf("unexpected property: %s", comp.IANAToken)
		}
	}
}

func TestConverterCallingConvertReturnsErrorWithUnknownComponents(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repository := mock_usecase.NewMockRepository(ctrl)
	given := ical.Calendar{
		Components: []ical.Component{
			&ical.VJournal{},
		},
	}
	dataSource := mock_usecase.NewMockDataSource(ctrl)
	repository.EXPECT().GetICal(dataSource).Return(&given, nil)
	converter := usecase.NewConverter(repository)

	// When
	_, err := converter.Convert(dataSource)

	// Then
	if err == nil {
		t.Error("expected error, but not")
	}
}
