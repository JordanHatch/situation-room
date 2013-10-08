package main

import (
	"time"

	calendar "code.google.com/p/google-api-go-client/calendar/v3"
)

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
