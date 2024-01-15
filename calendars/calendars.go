package calendars

import (
	"time"
)

type Event struct {
	ID          string
	CalendarID  string
	Title       string
	Description string
	Start       time.Time
	End         time.Time
}

type Calendar struct {
	ID     string
	Name   string
	Events []Event
}

var calendarsDb = []Calendar{
	{
		ID:     "tasks-calendar",
		Name:   "Tasks calendar",
		Events: []Event{},
	},
	{
		ID:     "meetings-calendar",
		Name:   "Meetings calendar",
		Events: []Event{},
	},
}

func AttatchDataPipeline() {
	// consume calendar events
	// consume event events
}

func GetCalendarById(id string) *Calendar {
	var cal *Calendar = nil
	for _, c := range calendarsDb {
		if c.ID == id {
			cal = &c
			break
		}
	}
	return cal
}

func GetEventById(calId, id string) *Event {
	var cal *Calendar = nil
	for _, c := range calendarsDb {
		if c.ID == id {
			cal = &c
			break
		}
	}

	if cal == nil {
		return nil
	}

	var event *Event = nil
	for _, e := range cal.Events {
		if e.ID == id {
			event = &e
			break
		}
	}
	return event
}
