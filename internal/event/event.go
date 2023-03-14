package event

type Event interface {
	EventName() string
	Decode(data any) (Event, error)
}
