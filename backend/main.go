package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"backend/container"
)

var (
	store = container.NewContainerStore()
)

type ImportResponse struct {
	Success      int      "json:\"success\""
	IncorrectPos []string "json:\"incorrectPositions\""
}

func NewImportResponse(success int, incorrectPos []string) *ImportResponse {
	return &ImportResponse{success, incorrectPos}
}

type containersHandler struct{}

func (c containersHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		BadRequestErrorHandler(writer, request)
		return
	}

	var containers []container.Container

	if err := json.NewDecoder(request.Body).Decode(&containers); err != nil {
		InternalServerErrorHandler(writer, request)
		return
	}

	success := 0
	var incorrectPos []string

	for _, cont := range containers {
		if err := store.Add(cont.Id, cont); err != nil {
			if errors.Is(err, container.IncorrectPosErr) {
				incorrectPos = append(incorrectPos, cont.Id)
			} else {
				InternalServerErrorHandler(writer, request)
				return
			}
		} else {
			success++
		}
	}

	response := NewImportResponse(success, incorrectPos)
	jsonBytes, err := json.Marshal(response)

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

type containersImportHandler struct{}

func (c containersImportHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
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
	var containers []container.Container
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
		cont := container.NewContainer(id, blockId, bayNum, stackNum, tierNum, arrivedAt)
		containers = append(containers, *cont)

		if err := store.Add(id, *cont); err != nil {
			if errors.Is(err, container.IncorrectPosErr) {
				incorrectPos = append(incorrectPos, cont.Id)
			} else {
				InternalServerErrorHandler(writer, request)
				return
			}
		} else {
			success++
		}
	}

	response := NewImportResponse(success, incorrectPos)
	jsonBytes, err := json.Marshal(response)

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

type containersSearchHandle struct{}

func (c containersSearchHandle) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
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
	var contS []container.Container
	var cont container.Container

	switch len(slicedUrl) {
	case 5:
		blockId, _ := strconv.ParseInt(slicedUrl[4], 10, 32)
		contS, err = store.GetBlock(blockId)
		break

	case 6:
		blockId, _ := strconv.ParseInt(slicedUrl[4], 10, 32)
		bayNum, _ := strconv.ParseInt(slicedUrl[5], 10, 32)
		contS, err = store.GetBay(blockId, bayNum)
		break

	case 7:
		blockId, _ := strconv.ParseInt(slicedUrl[4], 10, 32)
		bayNum, _ := strconv.ParseInt(slicedUrl[5], 10, 32)
		stackNum, _ := strconv.ParseInt(slicedUrl[6], 10, 32)
		contS, err = store.GetStack(blockId, bayNum, stackNum)
		break

	case 8:
		blockId, _ := strconv.ParseInt(slicedUrl[4], 10, 32)
		bayNum, _ := strconv.ParseInt(slicedUrl[5], 10, 32)
		stackNum, _ := strconv.ParseInt(slicedUrl[6], 10, 32)
		tierNum, _ := strconv.ParseInt(slicedUrl[7], 10, 32)
		cont, err = store.GetContainer(blockId, bayNum, stackNum, tierNum)
		contS = append(contS, cont)
		break
	}

	if err != nil {
		if errors.Is(err, container.NotFoundErr) {
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

type StatResponse struct {
	BlockId           int64   "json:\"blockId\""
	Capacity          float64 "json:\"capacity\""
	AverageAge        float64 "json:\"averageAge\""
	OldestContainerId string  "json:\"oldestContainerId\""
	NewestContainerId string  "json:\"newestContainerId\""
	EmptyPositions    int     "json:\"emptyPositions\""
	EmptyBays         int     "json:\"emptyBays\""
	EmptyStacks       int     "json:\"emptyStacks\""
}

func NewStatResponse(blockId int64, capacity float64, averageAge float64, oldestContainerId string,
	newestContainerId string, emptyPositions int, emptyBays int, emptyStacks int) *StatResponse {
	return &StatResponse{
		blockId,
		capacity,
		averageAge,
		oldestContainerId,
		newestContainerId,
		emptyPositions,
		emptyBays,
		emptyStacks,
	}
}

type blocksStatHandler struct{}

func (b blocksStatHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var responses []StatResponse

	for blockId := int64(1); blockId <= 4; blockId++ {
		currentBlock, err := store.GetBlock(blockId)

		if err != nil {
			continue
		}

		capacity, err := strconv.ParseFloat(fmt.Sprintf("%.2f",
			float64(125-len(currentBlock))/float64(125)*100), 64)

		if err != nil {
			InternalServerErrorHandler(writer, request)
			return
		}

		times := make(map[string]time.Time)
		ages := make(map[string]float64)
		agesSum := .0
		oldestId := currentBlock[0].Id
		newestId := currentBlock[0].Id
		count := 0
		nonEmptyBays := make(map[int64]int)
		nonEmptyStacks := make(map[int64]int)

		for _, c := range currentBlock {
			// the error here determines if the given timestamp is unix or "normal"
			// this works because the unix time stamp can be converted to a simple integer,
			// while the "normal" time stamp will fail
			unixTimeMsec, err := strconv.ParseInt(c.ArrivedAt, 10, 64)

			if err == nil {
				// Unix time stamp
				times[c.Id] = time.Unix(unixTimeMsec/1000, 0)
			} else {
				// "Normal" time stamp
				times[c.Id], err = time.Parse("2006-01-02T15:04:05.000Z", c.ArrivedAt)

				if err != nil {
					InternalServerErrorHandler(writer, request)
					return
				}
			}

			ages[c.Id] = (time.Now().Sub(times[c.Id]).Hours()) / 24
			agesSum += ages[c.Id]

			if ages[c.Id] > ages[oldestId] {
				oldestId = c.Id
			} else if ages[c.Id] < ages[newestId] {
				newestId = c.Id
			}

			count++
			nonEmptyBays[c.BayNum] = 1
			nonEmptyStacks[c.StackNum] = 1
		}

		averageAge, err := strconv.ParseFloat(fmt.Sprintf("%.2f", agesSum/float64(len(ages))), 64)

		if err != nil {
			InternalServerErrorHandler(writer, request)
			return
		}

		emptyPositions := 125 - count
		emptyBays := 5 - len(nonEmptyBays)
		emptyStacks := 5 - len(nonEmptyStacks)

		responses = append(responses, *NewStatResponse(
			blockId,
			capacity,
			averageAge,
			oldestId,
			newestId,
			emptyPositions,
			emptyBays,
			emptyStacks,
		))
	}

	jsonBytes, err := json.Marshal(responses)

	if err != nil {
		InternalServerErrorHandler(writer, request)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	if _, err = writer.Write(jsonBytes); err != nil {
		InternalServerErrorHandler(writer, request)
		return
	}
}

func InternalServerErrorHandler(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusInternalServerError)

	if _, err := writer.Write([]byte("500 internal server error")); err != nil {
		return
	}
}

func BadRequestErrorHandler(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusBadRequest)

	if _, err := writer.Write([]byte("400 bad request")); err != nil {
		return
	}
}

func NotFoundErrorHandler(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusNotFound)

	if _, err := writer.Write([]byte("404 not found")); err != nil {
		return
	}
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("/api/containers", &containersHandler{})
	mux.Handle("/api/containers/import", &containersImportHandler{})
	mux.Handle("/api/containers/search/", &containersSearchHandle{})
	mux.Handle("/api/blocks/stat", &blocksStatHandler{})

	if err := http.ListenAndServe(":3001", mux); err != nil {
		panic(err)
	}
}
