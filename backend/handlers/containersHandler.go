package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"backend/models"
)

type ContainersHandler struct {
	store *models.ContStore
}

func NewContainersHandler(store *models.ContStore) *ContainersHandler {
	return &ContainersHandler{store}
}

func (c ContainersHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		BadRequestErrorHandler(writer, request)
		return
	}

	var containers []models.Container

	if err := json.NewDecoder(request.Body).Decode(&containers); err != nil {
		InternalServerErrorHandler(writer, request)
		return
	}

	success := 0
	var incorrectPos []string

	for _, cont := range containers {
		if err := c.store.Add(cont.Id, cont); err != nil {
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
