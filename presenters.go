package main

import (
	"strconv"
	"time"
)

type EventApiResponse struct {
	Name    string `json:"name"`
	Creator string `json:"creator"`
	StartAt string `json:"start_at"`
	EndAt   string `json:"end_at"`
	Event   Event  `json:"-"`
}

type RoomApiResponse struct {
	Events         []EventApiResponse `json:"events"`
	Available      bool               `json:"available"`
	NextAvailable  string             `json:"next_available,omitempty"`
	AvailableUntil string             `json:"available_until,omitempty"`
	Room           Room               `json:"-"`
}

type RoomSetApiResponse struct {
	ResponseInfo map[string]string          `json:"_response_info"`
	Rooms        map[string]RoomApiResponse `json:"rooms"`
	RoomSet      RoomSet                    `json:"-"`
}

func (r RoomApiResponse) present() RoomApiResponse {
	response := RoomApiResponse{
		Events: []EventApiResponse{},
	}

	for _, event := range r.Room.Events {
		presentedEvent := EventApiResponse{
			Event: event,
		}

		response.Events = append(response.Events, presentedEvent.present())
	}

	response.Available = r.Room.Available()

	nextAvailable := r.Room.NextAvailable()
	if !r.Room.Available() && !nextAvailable.IsZero() {
		response.NextAvailable = nextAvailable.Format(time.RFC3339)
	}

	availableUntil := r.Room.AvailableUntil()
	if !availableUntil.IsZero() {
		response.AvailableUntil = availableUntil.Format(time.RFC3339)
	}

	return response
}

func (r EventApiResponse) present() EventApiResponse {
	r.Name = r.Event.Name()
	r.StartAt = r.Event.StartAt().Format(time.RFC3339)
	r.EndAt = r.Event.EndAt().Format(time.RFC3339)
	r.Creator = r.Event.Creator()

	return r
}

func (r RoomSetApiResponse) present(status string) RoomSetApiResponse {
	response := RoomSetApiResponse{
		ResponseInfo: make(map[string]string),
		Rooms:        make(map[string]RoomApiResponse),
	}
	response.ResponseInfo["status"] = status
	response.ResponseInfo["total_rooms"] = strconv.Itoa(r.RoomSet.TotalRooms)
	response.ResponseInfo["rooms_loaded"] = strconv.Itoa(r.RoomSet.RoomsLoaded)

	for _, room := range r.RoomSet.Rooms {
		presentedRoom := RoomApiResponse{
			Room: room,
		}
		response.Rooms[room.Name] = presentedRoom.present()
	}

	return response
}
