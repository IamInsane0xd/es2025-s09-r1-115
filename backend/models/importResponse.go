package models

type ImportResponse struct {
	Success      int      "json:\"success\""
	IncorrectPos []string "json:\"incorrectPositions\""
}

func NewImportResponse(success int, incorrectPos []string) *ImportResponse {
	return &ImportResponse{success, incorrectPos}
}
