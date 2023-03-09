package routes

import "net/http"

type RouteController interface {
	Routes() http.Handler
}
