package two_do

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gidoichi/ical-converter/entity/component"
)

type ActionType int

var (
	ActionTypeCall    ActionType = 5
	ActionTypeMessage ActionType = 6
	ActionTypeMail    ActionType = 7
	ActionTypeVisit   ActionType = 8
	ActionTypeURL     ActionType = 9
	ActionTypeSearch  ActionType = 10
)

type metadata struct {
	timeZone    time.Location
	ActionType  *ActionType `json:"actionType"`
	ActionValue *string     `json:"actionValue"`
	StartDate   *int64      `json:"StartDate"`
}

func parseMetadata(todo component.Todo, tz time.Location) (*metadata, error) {
	prop := todo.GetProperty("X-2DOAPP-METADATA")
	if prop == nil {
		return nil, nil
	}

	raw := prop.Value
	// Some content has a trailing `\n`, while others don't.
	content := strings.TrimSuffix(raw, "\n")
	content = strings.TrimSuffix(strings.TrimPrefix(content, "<2Do Meta>"), "</2Do Meta>")
	content, err := url.QueryUnescape(content)
	if err != nil {
		return nil, fmt.Errorf("unescape percent-encoding: %w", err)
	}

	var parsed metadata
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		return nil, fmt.Errorf("unmarshal json: %w", err)
	}

	parsed.timeZone = tz
	return &parsed, nil
}

func (m metadata) getURL() (url *string) {
	if m.ActionType == nil || *m.ActionType != ActionTypeURL || m.ActionValue == nil {
		return nil
	}

	return m.ActionValue
}

func (m metadata) getStartTime() (*time.Time, error) {
	if m.StartDate == nil {
		return nil, nil
	}
	// When startDate is equal to 0 it means that the task has no start date, not that the start date is 1990-01-01.
	if *m.StartDate == 0 {
		return nil, nil
	}

	// The startDate registered by 2Do is in the local time zone, not UTC.
	unix := time.Unix(*m.StartDate, 0)
	t, err := time.ParseInLocation(time.DateTime, unix.UTC().Format(time.DateTime), &m.timeZone)
	return &t, err
}
