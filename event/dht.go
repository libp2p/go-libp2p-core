package event

// TODO This is simply for being able to show DHT events on the UI for now.
// DhtEvent is an event related to the DHT.
type DhtEvent struct {
	// EventType is the type of the event that has occured.
	EventType string
	// EventJson is the JSON representation of the event payload.
	EventJson string
}