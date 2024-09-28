package two_do_test

import (
	"testing"
	"time"

	ical "github.com/arran4/golang-ical"
	"github.com/gidoichi/ical-converter/infrastructure/two_do"
	"github.com/gidoichi/ical-converter/usecase/mock_usecase"
	"go.uber.org/mock/gomock"
)

func TestTwoDoRepositoryCallingGetICalParsesSuccessfullyWhenComponentDoesNotHave2DoMetadata(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dataSource := mock_usecase.NewMockDataSource(ctrl)
	repository := two_do.NewTwoDoRepository(*time.FixedZone("JST", int((+9 * time.Hour).Seconds())))
	given := ical.Calendar{
		Components: []ical.Component{
			&ical.VTodo{ComponentBase: ical.ComponentBase{Properties: []ical.IANAProperty{
				{BaseProperty: ical.BaseProperty{
					IANAToken: "UID",
					Value:     "1",
				}},
			}}},
		},
	}
	dataSource.EXPECT().GetICal().Return(&given, nil)

	// When
	cal, err := repository.GetICal(dataSource)

	// Then

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if len(cal.Components) != 1 {
		t.Errorf("unexpected calendar: %s", cal.Serialize())
	}

	if len(cal.Components[0].UnknownPropertiesIANAProperties()) != 1 {
		t.Errorf("component does not have 1 property: %s", cal.Components[0])
	}
	if cal.Components[0].UnknownPropertiesIANAProperties()[0].IANAToken != string(ical.ComponentPropertyUniqueId) {
		t.Errorf("component does not have UID property: %v", cal.Components[0].UnknownPropertiesIANAProperties()[0])
	}
}

func TestTwoDoRepositoryCallingGetICalParsesSuccessfullyWhenComponentHas2DoMetadataContainingStartDateAndActionType(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dataSource := mock_usecase.NewMockDataSource(ctrl)
	repository := two_do.NewTwoDoRepository(*time.FixedZone("JST", int((+9 * time.Hour).Seconds())))
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
					Value:          "<2Do Meta>%7B%22RecurrenceValue%22%3A0%2C%22RecurrenceEndRepetitionsOrig%22%3A0%2C%22StartDate%22%3A1688601600%2C%22RecurrenceEndRepetitions%22%3A0%2C%22TaskType%22%3A0%2C%22TaskDuration%22%3A0%2C%22RecurrenceType%22%3A0%2C%22StartDayDelay%22%3A0%2C%22isExpandedToShowChildProjects%22%3A0%2C%22IsStarred%22%3A0%2C%22RecurrenceFrom%22%3A0%2C%22RecurrenceEndType%22%3A0%2C%22actionType%22%3A9%2C%22RUID%22%3A%22%22%2C%22actionValue%22%3A%22http%3A%5C%2F%5C%2Fexample.com%5C%2Fpub%5C%2Fcalendars%5C%2Fjsmith%5C%2Fmytime.ics%22%2C%22DisplayOrder%22%3A0%7D</2Do Meta>\n",
				}},
			}}},
		},
	}
	dataSource.EXPECT().GetICal().Return(&given, nil)

	// When
	cal, err := repository.GetICal(dataSource)

	// Then

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if len(cal.Components) != 1 {
		t.Errorf("unexpected calendar: %s", cal.Serialize())
	}

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
}

func TestTwoDoRepositoryCallingGetICalParsesSuccessfullyWhenComponentHas2DoMetadataNeitherContainingStartDateNorActionType(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dataSource := mock_usecase.NewMockDataSource(ctrl)
	repository := two_do.NewTwoDoRepository(*time.FixedZone("JST", int((+9 * time.Hour).Seconds())))
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
					Value:          "<2Do Meta>%7B%22RecurrenceValue%22%3A0%2C%22RecurrenceEndRepetitionsOrig%22%3A0%2C%22RecurrenceEndRepetitions%22%3A0%2C%22TaskType%22%3A0%2C%22TaskDuration%22%3A0%2C%22RecurrenceType%22%3A0%2C%22StartDayDelay%22%3A0%2C%22isExpandedToShowChildProjects%22%3A0%2C%22IsStarred%22%3A0%2C%22RecurrenceFrom%22%3A0%2C%22RecurrenceEndType%22%3A0%2C%22RUID%22%3A%22%22%2C%22DisplayOrder%22%3A0%7D</2Do Meta>\n",
				}},
			}}},
		},
	}
	dataSource.EXPECT().GetICal().Return(&given, nil)

	// When
	cal, err := repository.GetICal(dataSource)

	// Then

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if len(cal.Components) != 1 {
		t.Errorf("unexpected calendar: %s", cal.Serialize())
	}

	if len(cal.Components[0].UnknownPropertiesIANAProperties()) != 2 {
		t.Errorf("component does not have 2 properties: %s", cal.Components[0])
	}
	if cal.Components[0].UnknownPropertiesIANAProperties()[0].IANAToken != string(ical.ComponentPropertyUniqueId) {
		t.Errorf("component does not have UID property: %v", cal.Components[0].UnknownPropertiesIANAProperties()[0])
	}
	if cal.Components[0].UnknownPropertiesIANAProperties()[1].IANAToken != "X-2DOAPP-METADATA" {
		t.Errorf("component does not have X-2DOAPP-METADATA property: %v", cal.Components[0].UnknownPropertiesIANAProperties()[1])
	}
}

func TestTwoDoRepositoryCallingGetICalParsesSuccessfullyWhenComponentIsNotAllDayTask(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dataSource := mock_usecase.NewMockDataSource(ctrl)
	repository := two_do.NewTwoDoRepository(*time.FixedZone("JST", int((+9 * time.Hour).Seconds())))
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
					Value:          "<2Do Meta>%7B%22RecurrenceValue%22%3A0%2C%22RecurrenceEndRepetitionsOrig%22%3A0%2C%22StartDate%22%3A1705161600%2C%22RecurrenceEndRepetitions%22%3A0%2C%22TaskType%22%3A0%2C%22TaskDuration%22%3A0%2C%22RecurrenceType%22%3A0%2C%22StartDayDelay%22%3A0%2C%22isExpandedToShowChildProjects%22%3A0%2C%22IsStarred%22%3A0%2C%22RecurrenceFrom%22%3A0%2C%22RecurrenceEndType%22%3A0%2C%22RUID%22%3A%22%22%2C%22DisplayOrder%22%3A0%7D</2Do Meta>\n",
				}},
			}}},
		},
	}
	dataSource.EXPECT().GetICal().Return(&given, nil)

	// When
	cal, err := repository.GetICal(dataSource)

	// Then

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if len(cal.Components) != 1 {
		t.Errorf("unexpected calendar: %s", cal.Serialize())
	}

	if len(cal.Components[0].UnknownPropertiesIANAProperties()) != 3 {
		t.Errorf("component does not have 2 properties: %s", cal.Components[0])
	}
	if cal.Components[0].UnknownPropertiesIANAProperties()[0].IANAToken != string(ical.ComponentPropertyUniqueId) {
		t.Errorf("component does not have UID property: %v", cal.Components[0].UnknownPropertiesIANAProperties()[0])
	}
	if cal.Components[0].UnknownPropertiesIANAProperties()[1].IANAToken != "X-2DOAPP-METADATA" {
		t.Errorf("component does not have X-2DOAPP-METADATA property: %v", cal.Components[0].UnknownPropertiesIANAProperties()[1])
	}
	if cal.Components[0].UnknownPropertiesIANAProperties()[2].IANAToken != string(ical.PropertyDtstart) {
		t.Errorf("component does not have DTSTART property: %v", cal.Components[0].UnknownPropertiesIANAProperties()[2])
	}
	if cal.Components[0].UnknownPropertiesIANAProperties()[2].Value != "20240113T070000Z" {
		t.Errorf("unexpected DTSTART: %s", cal.Components[0].UnknownPropertiesIANAProperties()[2].Value)
	}
}

func TestTwoDoRepositoryCallingGetICalRemovesTimeRangeWhenComponentHasParent(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dataSource := mock_usecase.NewMockDataSource(ctrl)
	repository := two_do.NewTwoDoRepository(*time.FixedZone("JST", int((+9 * time.Hour).Seconds())))
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
					Value:          "<2Do Meta>%7B%7D</2Do Meta>\n",
				}},
				{BaseProperty: ical.BaseProperty{
					IANAToken: "RELATED-TO",
					Value:     "0",
				}},
				{BaseProperty: ical.BaseProperty{
					IANAToken: "DTSTART",
					Value:     "20070514T110000Z",
				}},
				{BaseProperty: ical.BaseProperty{
					IANAToken: "DUE",
					Value:     "20070709T130000Z",
				}},
			}}},
		},
	}
	dataSource.EXPECT().GetICal().Return(&given, nil)

	// When
	cal, err := repository.GetICal(dataSource)

	// Then

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if len(cal.Components) != 1 {
		t.Errorf("unexpected calendar: %s", cal.Serialize())
	}

	var (
		dtstart *ical.IANAProperty
		due     *ical.IANAProperty
	)
	for _, prop := range cal.Components[0].UnknownPropertiesIANAProperties() {
		prop := prop
		switch prop.IANAToken {
		case string(ical.ComponentPropertyDtStart):
			dtstart = &prop
		case string(ical.ComponentPropertyDue):
			due = &prop
		}
	}
	if dtstart != nil {
		t.Errorf("DTSTART Property is not removed: %s", cal.Serialize())
	}
	if due != nil {
		t.Errorf("DUE Property is not removed: %s", cal.Serialize())
	}
}
