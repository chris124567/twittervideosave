package web

import (
	"net/http"
)

func GenericErrorHandler(writer http.ResponseWriter, request *http.Request, statusCode int) {
	writer.WriteHeader(statusCode)
	err := WriteTemplate(writer, request, "./web/template/error.html.tmpl")
	if err != nil {
		HttpError(writer, templateError)
		return
	}
}
