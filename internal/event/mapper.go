package event

import "reflect"

type Mapper interface {
	RegisterEvent(eventType reflect.Type) error
	MapEvent(eventName string, data any) (Event, error)
	RegisteredEvents() []string
}
