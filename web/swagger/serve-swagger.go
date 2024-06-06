package swagger

import (
	"ecommerce/web/middlewares"
	"ecommerce/web/utils"
	"embed"
	"fmt"
	"mime"
	"net/http"
	"path"
	"strings"
)

//go:embed dist/*
var distFS embed.FS

//go:embed swagger.json
var swaggerFS embed.FS

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	filePath := r.PathValue("path")

	// if file path not specified serve index file
	if filePath == "" || filePath == "/" {
		filePath = "index.html"
	}

	// for swagger json file
	if strings.HasSuffix(filePath, "swagger.json") {
		data, err := swaggerFS.ReadFile("swagger.json")
		if err != nil {
			utils.SendError(w, http.StatusNotFound, err)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	}

	// for static dist files
	data, err := distFS.ReadFile(path.Join("dist", filePath))
	if err != nil {
		utils.SendError(w, http.StatusNotFound, fmt.Errorf("file not found"))
		return
	}
	ext := path.Ext(filePath)
	mime := mime.TypeByExtension(ext)
	w.Header().Add("Content-Type", mime)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func SetupSwagger(mux *http.ServeMux, manager *middlewares.Manager) {
	mux.Handle("GET /swagger/{path...}",
		manager.With(
			http.HandlerFunc(serveSwagger),
		),
	)
}
