package web

import (
	"ecommerce/config"
	"ecommerce/web/middlewares"
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

	go func() {
		defer wg.Done()

		conf := config.GetConfig()

		addr := fmt.Sprintf(":%s", conf.Http_port)
		slog.Info(fmt.Sprintf("Listening at %s", addr))

		if err := http.ListenAndServe(addr, handler); err != nil {
			slog.Error(err.Error())
		}
	}()

}
