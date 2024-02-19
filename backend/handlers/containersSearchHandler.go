package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
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
	re, err := regexp.Compile(`^/api/containers/search\?[a-zA-Z]+=[a-zA-Z0-9]+(&[a-zA-Z]+=[a-zA-Z0-9]+)*$`)

	if err != nil {
		InternalServerErrorHandler(writer, request)
		return
	}

	if !re.MatchString(request.URL.String()) {
		fmt.Println("1")
		BadRequestErrorHandler(writer, request)
		return
	}

	// query := utils.GetParamsFromURL(request.URL)
	query := request.URL.Query()

	if !query.Has("blockId") && !query.Has("bayNum") &&
		!query.Has("stackNum") && !query.Has("tierNum") &&
		!query.Has("id") {
		fmt.Println("2")
		BadRequestErrorHandler(writer, request)
		return
	}

	containers, err := c.store.GetByQuery(query)

	if err != nil {
		if errors.Is(err, models.NotFoundErr) {
			NotFoundErrorHandler(writer, request)
			return

		} else if errors.Is(err, models.ParamError) {
			fmt.Println("3")
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

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	if _, err := writer.Write(jsonBytes); err != nil {
		InternalServerErrorHandler(writer, request)
		return
	}
}
