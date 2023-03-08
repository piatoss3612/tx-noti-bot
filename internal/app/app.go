package app

type App interface {
	Setup() App
	Open() (<-chan bool, error)
	Close() error
}
