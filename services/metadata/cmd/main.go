package main

import (
	"log"
	"net/http"

	"sen1or/micromovie/services/metadata/internal/controller/metadata"
	httpHandler "sen1or/micromovie/services/metadata/internal/handler/http"
	"sen1or/micromovie/services/metadata/internal/repository"
)

func main() {
	repo := repository.NewRepository()
	ctrl := metadata.NewController(repo)
	handler := httpHandler.NewHandler(ctrl)

	http.HandleFunc("/metadata", handler.GetMetadata)
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Printf("error while running metadata server: %s\n", err)
	}
}
