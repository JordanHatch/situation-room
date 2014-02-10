# Situation Room

Fetches a collection of meeting room calendars from the Google Calendar API and exposes a list of the upcoming events and whether the room is available.

You can use the [situation-room-dashboard](https://github.com/JordanHatch/situation-room-dashboard) to present the data in this API as a dashboard.

## Configuration

* `MEETING_ROOM_CLIENT_ID` - the email address of a Google API service account to which the calendars are shared
* `MEETING_ROOM_API_KEY` - a base64 string (strict-encoded) of the `.pem` key for your Google API service account
* `MEETING_ROOM_CALENDARS` - a list of calendar IDs and names to use, using the format `<name>,<calendarId>;<name>,<calendarId>`
* `MEETING_ROOM_AUTH_USER` - the username required to access the api
* `MEETING_ROOM_AUTH_PASS` - the password required to access the api
* `PORT` - the port on which the API will run

## Note

I'm new to Go and so there's likely to be a lot here which isn't quite right, or breaking convention. Pointers and pull requests are welcome :)
