package objects

import "github.com/gmax79/otusgolang/microservices/internal/simple"

// Event - main structure to represent event in calendar
type Event struct {
	Alerttime simple.Date
	EventData
}

// EventData - information about event
type EventData struct {
	Information string
}

// SearchParameters - custom filters to search events
type SearchParameters struct {
	Year, Month, Week, Day int
}
