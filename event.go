package main

import (
	"log"
	"time"

	calendar "code.google.com/p/google-api-go-client/calendar/v3"
)

const dateTimeFormat = "2006-01-02T15:04:05Z07:00"

type Event struct {
	Source *calendar.Event
}

func (e Event) Name() string {
	return e.Source.Summary
}

func (e Event) Creator() string {
	if e.Source.Creator != nil {
		return e.Source.Creator.DisplayName
	}
	return ""
}

func (e Event) Visibility() string {
	if e.Source.Visibility == "private" || e.Source.Summary == "" {
		return "private"
	}
	return "public"
}

func (e Event) StartAt() time.Time {
	if e.Source.Start != nil {
		if e.Source.Start.DateTime != "" {
			t, err := time.Parse(dateTimeFormat, e.Source.Start.DateTime)
			if err != nil {
				log.Fatal("error parsing time:", e.Source.End.DateTime, err)
			} else {
				return t
			}
		} else if e.Source.Start.Date != "" {
			t, err := time.Parse("2006-01-02", e.Source.Start.Date)
			if err != nil {
				log.Fatal("error parsing time:", err)
			} else {
				return t
			}
		}
	}
	return time.Time{}
}

func (e Event) EndAt() time.Time {
	if e.Source.End != nil {
		if e.Source.End.DateTime != "" {
			t, err := time.Parse(dateTimeFormat, e.Source.End.DateTime)
			if err != nil {
				log.Fatal("error parsing time:", e.Source.End.DateTime, err)
			} else {
				return t
			}
		} else if e.Source.End.Date != "" {
			t, err := time.Parse("2006-01-02", e.Source.End.Date)
			if err != nil {
				log.Fatal("error parsing time:", err)
			} else {
				return t
			}
		}
	}
	return time.Time{}
}
