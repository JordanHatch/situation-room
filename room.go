package main

import (
	"time"

	calendar "code.google.com/p/google-api-go-client/calendar/v3"
)

// the amount of time (in minutes) in which the room should be free
// in order for the room to be 'next available'
//
const minRoomAvailabilityPeriod = 15

type Room struct {
	Name   string
	Events []Event
}

type RoomSet struct {
	Rooms       map[string]Room
	TotalRooms  int
	RoomsLoaded int
}

func (r Room) Available() bool {
	if len(r.Events) != 0 {
		firstEvent := r.Events[0]
		if firstEvent.StartAt().Before(time.Now()) && firstEvent.EndAt().After(time.Now()) {
			return false
		}
	}
	return true
}

func (r Room) NextAvailable() time.Time {
	if len(r.Events) == 0 {
		return time.Now()
	}

	var prevEvent = r.Events[0]
	var minTimeBeforeNextEvent time.Time

	for _, currentEvent := range r.Events {
		minTimeBeforeNextEvent = prevEvent.EndAt().Add(minRoomAvailabilityPeriod * time.Minute)

		if minTimeBeforeNextEvent.Before(currentEvent.StartAt()) {
			return prevEvent.EndAt()
		}

		prevEvent = currentEvent
	}

	return time.Time{}
}

func (r Room) AvailableUntil() time.Time {
	if len(r.Events) == 0 {
		return time.Time{}
	}

	firstEvent := r.Events[0]
	if firstEvent.StartAt().Before(time.Now()) {
		return time.Time{}
	}

	return firstEvent.StartAt()
}

func CreateRoomFromEvents(roomName string, calendarEvents []*calendar.Event) Room {
	room := Room{
		Name: roomName,
	}

	for _, calendarEvent := range calendarEvents {
		event := Event{
			Source: calendarEvent,
		}
		room.Events = append(room.Events, event)
	}
	return room
}
