package main

import (
	"net/http"
	"sen1or/micromovie/services/movie/internal/controller/movie"
	metadatagateway "sen1or/micromovie/services/movie/internal/gateway/metadata/http"
	ratinggateway "sen1or/micromovie/services/movie/internal/gateway/rating/http"
	httphandler "sen1or/micromovie/services/movie/internal/handler/http"
)

func main() {
	metadataGateway := metadatagateway.NewGateway("localhost:8081")
	ratingGateway := ratinggateway.NewGateway("localhost:8082")

	ctrl := movie.New(ratingGateway, metadataGateway)
	handler := httphandler.New(*ctrl)
	http.HandleFunc("/movie", handler.GetMovieDetails)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
