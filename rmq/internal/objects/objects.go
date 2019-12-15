package objects

// Event - information about event
type Event string

// SearchParameters - custom filters to search events
type SearchParameters struct {
	Year, Month, Week, Day int
}
