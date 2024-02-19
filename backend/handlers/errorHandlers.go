package handlers

import "net/http"

func InternalServerErrorHandler(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.WriteHeader(http.StatusInternalServerError)

	if _, err := writer.Write([]byte("500 internal server error")); err != nil {
		return
	}
}

func BadRequestErrorHandler(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.WriteHeader(http.StatusBadRequest)

	if _, err := writer.Write([]byte("400 bad request")); err != nil {
		return
	}
}

func NotFoundErrorHandler(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.WriteHeader(http.StatusNotFound)

	if _, err := writer.Write([]byte("404 not found")); err != nil {
		return
	}
}
