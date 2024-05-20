package web

import (
	"ecommerce/web/handlers"
	"ecommerce/web/middlewares"
	"net/http"
)

func InitRouts(mux *http.ServeMux, manager *middlewares.Manager) {
	mux.Handle(
		"POST /xx",
		manager.With(
			http.HandlerFunc(handlers.Register),
		),
	)

	mux.Handle(
		"GET /xx",
		manager.With(
			http.HandlerFunc(handlers.Login),
		),
	)

	mux.Handle(
		"GET /up",
		manager.With(
			http.HandlerFunc(handlers.Update),
		),
	)
}
