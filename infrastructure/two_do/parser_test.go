package two_do_test

import (
	"log"
	"testing"
	"time"

	ical "github.com/arran4/golang-ical"
	"github.com/gidoichi/ical-converter/infrastructure/two_do"
	"github.com/gidoichi/ical-converter/usecase/mock_usecase"
	"go.uber.org/mock/gomock"
)

func TestGetICal(t *testing.T) {
	// Given

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dataSource := mock_usecase.NewMockDataSource(ctrl)
	repository := two_do.NewTwoDoRepository(*time.UTC)
	given := ical.Calendar{
		Components: []ical.Component{
			&ical.VTodo{ComponentBase: ical.ComponentBase{Properties: []ical.IANAProperty{
				{BaseProperty: ical.BaseProperty{
					IANAToken: "UID",
					Value:     "1",
				}},
				{BaseProperty: ical.BaseProperty{
					IANAToken:      "X-2DOAPP-METADATA",
					ICalParameters: map[string][]string{"SHARE-SCOPE": {"GLOBAL"}},
					Value:          "<2Do Meta>%7B%22RecurrenceValue%22%3A0%2C%22RecurrenceEndRepetitionsOrig%22%3A0%2C%22StartDate%22%3A1688601600%2C%22RecurrenceEndRepetitions%22%3A0%2C%22TaskType%22%3A0%2C%22TaskDuration%22%3A0%2C%22RecurrenceType%22%3A0%2C%22StartDayDelay%22%3A0%2C%22isExpandedToShowChildProjects%22%3A0%2C%22IsStarred%22%3A0%2C%22RecurrenceFrom%22%3A0%2C%22RecurrenceEndType%22%3A0%2C%22actionType%22%3A9%2C%22RUID%22%3A%22%22%2C%22actionValue%22%3A%22http%3A%5C%2F%5C%2Fexample.com%5C%2Fpub%5C%2Fcalendars%5C%2Fjsmith%5C%2Fmytime.ics%22%2C%22DisplayOrder%22%3A0%7D</2Do Meta>\\n",
				}},
			}}},
			&ical.VTodo{ComponentBase: ical.ComponentBase{Properties: []ical.IANAProperty{
				{BaseProperty: ical.BaseProperty{
					IANAToken: "UID",
					Value:     "2",
				}},
				{BaseProperty: ical.BaseProperty{
					IANAToken:      "X-2DOAPP-METADATA",
					ICalParameters: map[string][]string{"SHARE-SCOPE": {"GLOBAL"}},
					Value:          "<2Do Meta>%7B%22RecurrenceValue%22%3A0%2C%22RecurrenceEndRepetitionsOrig%22%3A0%2C%22RecurrenceEndRepetitions%22%3A0%2C%22TaskType%22%3A0%2C%22TaskDuration%22%3A0%2C%22RecurrenceType%22%3A0%2C%22StartDayDelay%22%3A0%2C%22isExpandedToShowChildProjects%22%3A0%2C%22IsStarred%22%3A0%2C%22RecurrenceFrom%22%3A0%2C%22RecurrenceEndType%22%3A0%2C%22RUID%22%3A%22%22%2C%22DisplayOrder%22%3A0%7D</2Do Meta>\\n",
				}},
			}}},
			&ical.VTodo{ComponentBase: ical.ComponentBase{Properties: []ical.IANAProperty{
				{BaseProperty: ical.BaseProperty{
					IANAToken: "UID",
					Value:     "3",
				}},
			}}},
			&ical.VTodo{ComponentBase: ical.ComponentBase{Properties: []ical.IANAProperty{
				{BaseProperty: ical.BaseProperty{
					IANAToken: "UID",
					Value:     "4",
				}},
				{BaseProperty: ical.BaseProperty{
					IANAToken:      "X-2DOAPP-METADATA",
					ICalParameters: map[string][]string{"SHARE-SCOPE": {"GLOBAL"}},
					Value:          "<2Do Meta>%7B%22RecurrenceValue%22%3A0%2C%22RecurrenceEndRepetitionsOrig%22%3A0%2C%22StartDate%22%3A1705161600%2C%22RecurrenceEndRepetitions%22%3A0%2C%22TaskType%22%3A0%2C%22TaskDuration%22%3A0%2C%22RecurrenceType%22%3A0%2C%22StartDayDelay%22%3A0%2C%22isExpandedToShowChildProjects%22%3A0%2C%22IsStarred%22%3A0%2C%22RecurrenceFrom%22%3A0%2C%22RecurrenceEndType%22%3A0%2C%22RUID%22%3A%22%22%2C%22DisplayOrder%22%3A0%7D</2Do Meta>\\n",
				}},
			}}},
		},
	}
	dataSource.EXPECT().GetICal().Return(&given, nil)

	// When

	cal, err := repository.GetICal(dataSource)

	// Then

	log.Println("iCal is proper")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if len(cal.Components) != 4 {
		t.Errorf("unexpected calendar: %s", cal.Serialize())
	}

	log.Println("Component with 2Do app metadata having StartDate and ActionType are parsed")
	var (
		url     *ical.IANAProperty
		dtstart *ical.IANAProperty
	)
	for _, prop := range cal.Components[0].UnknownPropertiesIANAProperties() {
		prop := prop
		switch prop.IANAToken {
		case string(ical.ComponentPropertyUrl):
			url = &prop
		case string(ical.ComponentPropertyDtStart):
			dtstart = &prop
		}
	}
	if url == nil {
		t.Errorf("URL Property is not found: %s", cal.Serialize())
	}
	if url.Value != "http://example.com/pub/calendars/jsmith/mytime.ics" {
		t.Errorf("unexpected URL: %s", url.Value)
	}
	if dtstart == nil {
		t.Errorf("DTSTART Property is not found: %s", cal.Serialize())
	}
	if dtstart.Value != "20230706" {
		t.Errorf("unexpected DTSTART: %s", dtstart.Value)
	}

	log.Println("Component with 2Do app metadata not having StartDate and ActionType are parsed")
	if len(cal.Components[1].UnknownPropertiesIANAProperties()) != 2 {
		t.Errorf("component does not have 2 properties: %s", cal.Components[1])
	}
	if cal.Components[1].UnknownPropertiesIANAProperties()[0].IANAToken != string(ical.ComponentPropertyUniqueId) {
		t.Errorf("component does not have UID property: %v", cal.Components[1].UnknownPropertiesIANAProperties()[0])
	}
	if cal.Components[1].UnknownPropertiesIANAProperties()[1].IANAToken != "X-2DOAPP-METADATA" {
		t.Errorf("component does not have X-2DOAPP-METADATA property: %v", cal.Components[1].UnknownPropertiesIANAProperties()[1])
	}

	log.Println("Component without 2Do app metadata are parsed")
	if len(cal.Components[2].UnknownPropertiesIANAProperties()) != 1 {
		t.Errorf("component does not have 1 property: %s", cal.Components[2])
	}
	if cal.Components[2].UnknownPropertiesIANAProperties()[0].IANAToken != string(ical.ComponentPropertyUniqueId) {
		t.Errorf("component does not have UID property: %v", cal.Components[2].UnknownPropertiesIANAProperties()[0])
	}

	log.Println("Component with not all-day task is parsed")
	if len(cal.Components[3].UnknownPropertiesIANAProperties()) != 3 {
		t.Errorf("component does not have 2 properties: %s", cal.Components[3])
	}
	if cal.Components[3].UnknownPropertiesIANAProperties()[0].IANAToken != string(ical.ComponentPropertyUniqueId) {
		t.Errorf("component does not have UID property: %v", cal.Components[3].UnknownPropertiesIANAProperties()[0])
	}
	if cal.Components[3].UnknownPropertiesIANAProperties()[1].IANAToken != "X-2DOAPP-METADATA" {
		t.Errorf("component does not have X-2DOAPP-METADATA property: %v", cal.Components[3].UnknownPropertiesIANAProperties()[1])
	}
	if cal.Components[3].UnknownPropertiesIANAProperties()[2].IANAToken != string(ical.PropertyDtstart) {
		t.Errorf("component does not have DTSTART property: %v", cal.Components[3].UnknownPropertiesIANAProperties()[2])
	}
	if cal.Components[3].UnknownPropertiesIANAProperties()[2].Value != "20240113T070000Z" {
		t.Errorf("unexpected DTSTART: %s", cal.Components[3].UnknownPropertiesIANAProperties()[2].Value)
	}
}
