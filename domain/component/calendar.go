package component

import (
	ical "github.com/arran4/golang-ical"
)

type Calendar ical.Calendar

func (c *Calendar) AddComponent(sub ical.Component) {
	c.Components = append(c.Components, sub)
}
