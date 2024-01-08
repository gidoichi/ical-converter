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
	ActionType  *ActionType `json:"actionType"`
	ActionValue *string     `json:"actionValue"`
	StartDate   *int64      `json:"StartDate"`
}

func parseMetadata(todo component.Todo) (*metadata, error) {
	prop := todo.GetProperty("X-2DOAPP-METADATA")
	if prop == nil {
		return nil, nil
	}

	raw := prop.Value
	content := strings.TrimSuffix(strings.TrimPrefix(raw, "<2Do Meta>"), "</2Do Meta>\\n")
	content, err := url.QueryUnescape(content)
	if err != nil {
		return nil, fmt.Errorf("failed to unescape percent-encoding: %w", err)
	}

	var parsed metadata
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		return nil, fmt.Errorf("failed to unmarshall json: %w", err)
	}

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
	t := time.Unix(*m.StartDate, 0)
	return &t, nil
}
