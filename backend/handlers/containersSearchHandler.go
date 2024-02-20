package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"backend/models"
)

type ContainersSearchHandle struct {
	store *models.ContStore
}

func NewContainersSearchHandler(store *models.ContStore) *ContainersSearchHandle {
	return &ContainersSearchHandle{store}
}

func (c ContainersSearchHandle) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	re1, _ := regexp.Compile(`^/api/containers/search\?[a-zA-Z]+=[a-zA-Z0-9]+(&[a-zA-Z]+=[a-zA-Z0-9]+)*$`)
	re2, _ := regexp.Compile(`^/api/containers/search$`)

	if !re1.MatchString(request.URL.String()) && !re2.MatchString(request.URL.String()) {
		BadRequestErrorHandler(writer, request)
		return
	}

	query := request.URL.Query()
	containers, err := c.store.GetByQuery(query)

	if err != nil {
		if errors.Is(err, models.NotFoundErr) {
			NotFoundErrorHandler(writer, request)
			return

		} else if errors.Is(err, models.ParamError) {
			BadRequestErrorHandler(writer, request)
			return
		}

		InternalServerErrorHandler(writer, request)
		return
	}

	jsonBytes, err := json.Marshal(containers)

	if err != nil {
		InternalServerErrorHandler(writer, request)
		return
	}

	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	if _, err := writer.Write(jsonBytes); err != nil {
		InternalServerErrorHandler(writer, request)
		return
	}
}
