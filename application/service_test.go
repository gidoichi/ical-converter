package application_test

import (
	"log"
	"testing"

	ical "github.com/arran4/golang-ical"
	"github.com/gidoichi/ical-converter/application"
	"github.com/gidoichi/ical-converter/usecase/mock_usecase"
	"go.uber.org/mock/gomock"
)

func TestConvert(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	converter := mock_usecase.NewMockConverter(ctrl)
	given := ical.Calendar{
		Components: []ical.Component{
			&ical.VEvent{
				ComponentBase: ical.ComponentBase{
					Properties: []ical.IANAProperty{
						{BaseProperty: ical.BaseProperty{IANAToken: "UID", Value: "1"}},
						{BaseProperty: ical.BaseProperty{IANAToken: "DTSTART", Value: "20060102T150405Z"}},
					},
				},
			},
			&ical.VEvent{
				ComponentBase: ical.ComponentBase{
					Properties: []ical.IANAProperty{
						{BaseProperty: ical.BaseProperty{IANAToken: "UID", Value: "2"}},
						{BaseProperty: ical.BaseProperty{IANAToken: "DTSTART", Value: "20060102T150405Z"}},
						{BaseProperty: ical.BaseProperty{IANAToken: "STATUS", Value: "TENTATIVE"}},
					},
				},
			},
			&ical.VEvent{
				ComponentBase: ical.ComponentBase{
					Properties: []ical.IANAProperty{
						{BaseProperty: ical.BaseProperty{IANAToken: "UID", Value: "10"}},
						{BaseProperty: ical.BaseProperty{IANAToken: "DTSTART", Value: "20060102T150405Z"}},
						{BaseProperty: ical.BaseProperty{IANAToken: "STATUS", Value: "CANCELLED"}},
					},
				},
			},
			&ical.VEvent{},
			&ical.VJournal{},
		},
	}
	dataSource := mock_usecase.NewMockDataSource(ctrl)
	converter.EXPECT().Convert(dataSource).Return(&given, nil)
	service := application.NewConvertService(converter)

	// When
	cal, err := service.Convert(dataSource)

	// Then
	log.Println("Calendar is serialized")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	expectCal := ical.Calendar{
		Components: []ical.Component{
			&ical.VEvent{
				ComponentBase: ical.ComponentBase{
					Properties: []ical.IANAProperty{
						{BaseProperty: ical.BaseProperty{IANAToken: "UID", Value: "1"}},
						{BaseProperty: ical.BaseProperty{IANAToken: "DTSTART", Value: "20060102T150405Z"}},
					},
				},
			},
			&ical.VEvent{
				ComponentBase: ical.ComponentBase{
					Properties: []ical.IANAProperty{
						{BaseProperty: ical.BaseProperty{IANAToken: "UID", Value: "2"}},
						{BaseProperty: ical.BaseProperty{IANAToken: "DTSTART", Value: "20060102T150405Z"}},
					},
				},
			},
		},
	}
	expect := expectCal.Serialize()
	if cal != expect {
		t.Errorf("unexpected serialized calendar: %s", cal)
	}
}
