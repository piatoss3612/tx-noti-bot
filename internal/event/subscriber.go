package event

type Subscriber interface {
	Subscribe(topics ...string) (<-chan Event, <-chan error, error)
	Mapper() Mapper
	Close() error
}
