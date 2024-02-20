package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"backend/handlers"
	"backend/models"
)

var (
	store = models.NewContainerStore()
)

func main() {
	res, err := http.Get("http://localhost:3000/containers")

	if err != nil {
		fmt.Printf("error: failed to make request to json server (no container data pre-loaded): %s\n", err)
	} else {
		var containers []models.Container

		if err = json.NewDecoder(res.Body).Decode(&containers); err != nil {
			panic(err)
		}

		for _, val := range containers {
			// error ignored, because the only error this can cause is invalid position or container id but that cannot
			// happen, or is otherwise not important.
			_ = store.Add(val.Id, val)
		}
	}

	fmt.Println("Containers imported, creating mux!")

	mux := http.NewServeMux()

	mux.Handle("/api/containers", handlers.NewContainersHandler(store))
	mux.Handle("/api/containers/import", handlers.NewContainersImportHandler(store))
	mux.Handle("/api/containers/search", handlers.NewContainersSearchHandler(store))
	mux.Handle("/api/blocks/stat", handlers.NewBlocksStatHandler(store))

	fmt.Println("Mux created, starting server!")

	if err = http.ListenAndServe(":3001", mux); err != nil {
		panic(err)
	}
}
