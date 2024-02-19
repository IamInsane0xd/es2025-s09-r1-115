package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"backend/models"
)

type ContainersImportHandler struct {
	store *models.ContStore
}

func NewContainersImportHandler(store *models.ContStore) *ContainersImportHandler {
	return &ContainersImportHandler{store}
}

func (c ContainersImportHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		BadRequestErrorHandler(writer, request)
		return
	}

	bodyBytes, err := io.ReadAll(request.Body)

	if err != nil {
		InternalServerErrorHandler(writer, request)
		return
	}

	lines := strings.Split(string(bodyBytes), "\n")
	var containers []models.Container
	success := 0
	var incorrectPos []string

	for i, line := range lines {
		if i == 0 {
			continue
		}

		data := strings.Split(line, ",")
		id := data[0]
		blockId, err := strconv.ParseInt(data[1], 10, 32)

		if err != nil {
			InternalServerErrorHandler(writer, request)
			return
		}

		bayNum, err := strconv.ParseInt(data[2], 10, 32)

		if err != nil {
			InternalServerErrorHandler(writer, request)
			return
		}

		stackNum, err := strconv.ParseInt(data[3], 10, 32)

		if err != nil {
			InternalServerErrorHandler(writer, request)
			return
		}

		tierNum, err := strconv.ParseInt(data[4], 10, 32)

		if err != nil {
			InternalServerErrorHandler(writer, request)
			return
		}

		arrivedAt := data[5]
		cont := models.NewContainer(id, blockId, bayNum, stackNum, tierNum, arrivedAt)
		containers = append(containers, *cont)

		if err := c.store.Add(id, *cont); err != nil {
			if errors.Is(err, models.IncorrectPosErr) {
				incorrectPos = append(incorrectPos, cont.Id)
			} else {
				InternalServerErrorHandler(writer, request)
				return
			}
		} else {
			success++
		}
	}

	response := models.NewImportResponse(success, incorrectPos)
	jsonBytes, err := json.Marshal(response)

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
