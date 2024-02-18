package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"backend/models"
)

type BlocksStatHandler struct {
	store *models.ContStore
}

func NewBlocksStatHandler(store *models.ContStore) *BlocksStatHandler {
	return &BlocksStatHandler{store}
}

func (b BlocksStatHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var responses []models.StatResponse

	for blockId := int64(1); blockId <= 4; blockId++ {
		currentBlock, err := b.store.GetBlock(blockId)

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

		responses = append(responses, *models.NewStatResponse(
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
