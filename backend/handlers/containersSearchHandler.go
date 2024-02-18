package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"backend/models"
)

type ContainersSearchHandle struct {
	store *models.ContStore
}

func NewContainersSearchHandler(store *models.ContStore) *ContainersSearchHandle {
	return &ContainersSearchHandle{store}
}

func (c ContainersSearchHandle) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	re, err := regexp.Compile("^/api/containers/search/[1-4](/[1-5]){0,3}$")
	url := request.URL.String()

	switch {
	case err != nil:
		InternalServerErrorHandler(writer, request)
		return

	case !re.MatchString(url):
		BadRequestErrorHandler(writer, request)
		return
	}

	slicedUrl := strings.Split(url, "/")
	var contS []models.Container
	var cont models.Container

	switch len(slicedUrl) {
	case 5:
		blockId, _ := strconv.ParseInt(slicedUrl[4], 10, 32)
		contS, err = c.store.GetBlock(blockId)
		break

	case 6:
		blockId, _ := strconv.ParseInt(slicedUrl[4], 10, 32)
		bayNum, _ := strconv.ParseInt(slicedUrl[5], 10, 32)
		contS, err = c.store.GetBay(blockId, bayNum)
		break

	case 7:
		blockId, _ := strconv.ParseInt(slicedUrl[4], 10, 32)
		bayNum, _ := strconv.ParseInt(slicedUrl[5], 10, 32)
		stackNum, _ := strconv.ParseInt(slicedUrl[6], 10, 32)
		contS, err = c.store.GetStack(blockId, bayNum, stackNum)
		break

	case 8:
		blockId, _ := strconv.ParseInt(slicedUrl[4], 10, 32)
		bayNum, _ := strconv.ParseInt(slicedUrl[5], 10, 32)
		stackNum, _ := strconv.ParseInt(slicedUrl[6], 10, 32)
		tierNum, _ := strconv.ParseInt(slicedUrl[7], 10, 32)
		cont, err = c.store.GetContainer(blockId, bayNum, stackNum, tierNum)
		contS = append(contS, cont)
		break
	}

	if err != nil {
		if errors.Is(err, models.NotFoundErr) {
			NotFoundErrorHandler(writer, request)
		} else {
			InternalServerErrorHandler(writer, request)
		}

		return
	}

	jsonBytes, err := json.Marshal(contS)

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
