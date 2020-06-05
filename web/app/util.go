package web

import (
	"net/http"
)

// 301 redirect a user to a given URL
func Redirect(url string, writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, url, 301)
}

func HttpError(writer http.ResponseWriter, statusCode int) {
	http.Error(writer, http.StatusText(statusCode), statusCode)
}
