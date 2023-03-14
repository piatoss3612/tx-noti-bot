package mapper

import (
	"fmt"
	"reflect"

	"github.com/piatoss3612/tx-notification/internal/event"
)

type dynamicEventMapper struct {
	eventMap map[string]reflect.Type
}

func New() event.Mapper {
	return &dynamicEventMapper{
		eventMap: make(map[string]reflect.Type),
	}
}

func (m *dynamicEventMapper) RegisterEvent(eventType reflect.Type) error {
	iface := reflect.New(eventType).Interface()

	event, ok := iface.(event.Event)
	if !ok {
		return fmt.Errorf("type %T does not implement the Event interface", iface)
	}

	m.eventMap[event.EventName()] = eventType
	return nil
}

func (m *dynamicEventMapper) MapEvent(eventName string, data any) (event.Event, error) {
	target, ok := m.eventMap[eventName]
	if !ok {
		return nil, fmt.Errorf("event %s is not registered", eventName)
	}

	iface := reflect.New(target).Interface()

	event, ok := iface.(event.Event)
	if !ok {
		delete(m.eventMap, eventName)
		return nil, fmt.Errorf("type %T does not implement the Event interface", iface)
	}

	return event.Decode(data)
}

func (m *dynamicEventMapper) RegisteredEvents() []string {
	keys := make([]string, 0, len(m.eventMap))

	for k := range m.eventMap {
		keys = append(keys, k)
	}

	return keys
}
