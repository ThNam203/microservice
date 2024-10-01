package main

import (
	"net/http"
	"sen1or/micromovie/services/rating/internal/controller/rating"
	httphandler "sen1or/micromovie/services/rating/internal/handler/http"
	"sen1or/micromovie/services/rating/internal/repository/memory"
)

func main() {
	repo := memory.NewRepository()
	ctrl := rating.New(repo)
	handler := httphandler.New(ctrl)
	http.HandleFunc("/rating", handler.Handle)
	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}
}
