package handler

type Handler interface {
	Inject(target any) error
	Cleanup() error
}
