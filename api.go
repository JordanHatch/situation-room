package main

import (
	auth "github.com/abbot/go-http-auth"

	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var port = os.Getenv("PORT")

var calendarConfig = os.Getenv("MEETING_ROOM_CALENDARS")
var calendars map[string]string = make(map[string]string)

var client ApiClient
var googleApiKey = os.Getenv("MEETING_ROOM_API_KEY")
var googleClientId = os.Getenv("MEETING_ROOM_CLIENT_ID")

var authUsername = os.Getenv("MEETING_ROOM_AUTH_USER")
var authPassword = os.Getenv("MEETING_ROOM_AUTH_PASS")

var rooms map[string]Room = make(map[string]Room)

func main() {
	client = ApiClient{
		ClientId:   googleClientId,
		EncodedKey: googleApiKey,
	}
	calendars = parseCalendarConfig(calendarConfig)

	startTicker()
	go loadEvents()

	log.Println("API is starting up on :" + port)
	log.Println("Use Ctrl+C to stop")

	authenticator := auth.NewBasicAuthenticator("situation-room", Authenticate)
	http.HandleFunc("/", authenticator.Wrap(eventsHandler))
	http.ListenAndServe(":"+port, nil)
}

func Authenticate(user, realm string) string {
	if user == authUsername {
		d := sha1.New()
		d.Write([]byte(authPassword))
		e := base64.StdEncoding.EncodeToString(d.Sum(nil))

		return "{SHA}" + e
	}
	return ""
}

func eventsHandler(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	w.Header().Set("Content-Type", "application/json")

	roomSet := RoomSet{
		Rooms:       rooms,
		TotalRooms:  len(calendars),
		RoomsLoaded: len(rooms),
	}
	apiResponse := RoomSetApiResponse{
		RoomSet:      roomSet,
		ResponseInfo: make(map[string]string),
	}

	status := "ok"
	if roomSet.TotalRooms > roomSet.RoomsLoaded {
		status = "incomplete"
	}

	b, err := json.Marshal(apiResponse.present(status))
	if err != nil {
		log.Fatal("Error preparing JSON: ", err)
	}
	response := string(b)
	fmt.Fprintf(w, response)
}

func loadEvents() {
	log.Print("Loading events...")

	client.Token = client.GetToken()

	for calendarName, calendarId := range calendars {
		go loadEventsForRoom(calendarName, calendarId)
	}
}

func loadEventsForRoom(calendarName string, calendarId string) {
	log.Printf("Loading %v", calendarName)
	events, err := client.Api().Events.List(calendarId).
		TimeMin(time.Now().Format(time.RFC3339)).
		TimeMax(time.Now().Add(24 * time.Hour).Format(time.RFC3339)).
		SingleEvents(true).
		OrderBy("startTime").Do()
	if err != nil {
		log.Fatal("Error making the calendar api request: ", err)
	}

	rooms[calendarName] = CreateRoomFromEvents(calendarName, events.Items)
	log.Printf("Finished loading %v events for %v", len(rooms[calendarName].Events), calendarName)
}

func parseCalendarConfig(config string) map[string]string {
	calendarMap := map[string]string{}
	lines := strings.Split(config, ";")

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		name := parts[0]
		url := parts[1]

		calendarMap[name] = url
	}

	return calendarMap
}

func startTicker() {
	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				loadEvents()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
