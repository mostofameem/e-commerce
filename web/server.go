package web

import (
	"ecommerce/config"
	"ecommerce/web/middlewares"
	"ecommerce/web/swagger"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
)

func StartServer(wg *sync.WaitGroup) {
	manager := middlewares.NewManager()
	mux := http.NewServeMux()

	InitRouts(mux, manager)

	handler := middlewares.EnableCors(mux)

	wg.Add(1)
	swagger.SetupSwagger(mux, manager)

	go func() {
		defer wg.Done()

		conf := config.GetConfig()

		addr := fmt.Sprintf(":%d", conf.HttpPort)
		slog.Info(fmt.Sprintf("Listening at %s", addr))

		if err := http.ListenAndServe(addr, handler); err != nil {
			slog.Error(err.Error())
		}
	}()

}
