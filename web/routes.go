package web

import (
	"ecommerce/web/handlers"
	"ecommerce/web/middlewares"
	"net/http"
)

func InitRouts(mux *http.ServeMux, manager *middlewares.Manager) {
	mux.Handle(
		"POST /users",
		manager.With(
			http.HandlerFunc(handlers.Register),
		),
	)

	mux.Handle(
		"GET /users",
		manager.With(
			http.HandlerFunc(handlers.Login),
		),
	)
}
