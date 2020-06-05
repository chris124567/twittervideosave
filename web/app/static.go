package web

import (
	"net/http"
	"os"
)

func StaticHandler(writer http.ResponseWriter, request *http.Request) {
	// Basically just reads file of a given path, i.e. static/main.css
	path := "./web" + request.URL.Path
	file, err := os.Stat(path)
	if err == nil && !file.IsDir() {
		writer.Header().Set("Vary", "Accept-Encoding")
		writer.Header().Set("Cache-Control", "public, max-age=7776000")
		http.ServeFile(writer, request, path)
		return
	}

	http.NotFound(writer, request)
}
