package events

type Event int

const (
	EVENT_NULL Event = iota
	EVENT_CHUNK_UPDATE
)
