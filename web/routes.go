package web

import (
	"ecommerce/web/handlers"
	"ecommerce/web/middlewares"
	"net/http"
)

func InitRouts(mux *http.ServeMux, manager *middlewares.Manager) {
	mux.Handle(
		"POST /register",
		manager.With(
			http.HandlerFunc(handlers.Register),
		),
	)

	mux.Handle(
		"POST /login",
		manager.With(
			http.HandlerFunc(handlers.Login),
		),
	)

	mux.Handle(
		"POST /register/{id}",
		manager.With(
			http.HandlerFunc(handlers.Verification),
		),
	)

	mux.Handle(
		"POST /products",
		manager.With(
			http.HandlerFunc(handlers.BuyProduct),
			middlewares.AuthenticateJWT,
		),
	)

	mux.Handle(
		"GET /cart",
		manager.With(
			http.HandlerFunc(handlers.ShowCart),
			middlewares.AuthenticateJWT,
		),
	)

}
