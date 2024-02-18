package models

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
