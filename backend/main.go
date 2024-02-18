package main

import (
	"net/http"

	"backend/handlers"
	"backend/models"
)

var (
	store = models.NewContainerStore()
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/api/containers", handlers.NewContainersHandler(store))
	mux.Handle("/api/containers/import", handlers.NewContainersImportHandler(store))
	mux.Handle("/api/containers/search/", handlers.NewContainersSearchHandler(store))
	mux.Handle("/api/blocks/stat", handlers.NewBlocksStatHandler(store))

	if err := http.ListenAndServe(":3001", mux); err != nil {
		panic(err)
	}
}
