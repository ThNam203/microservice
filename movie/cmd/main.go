package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sen1or/micromovie/movie/internal/controller/movie"
	metadatagateway "sen1or/micromovie/movie/internal/gateway/metadata/http"
	ratinggateway "sen1or/micromovie/movie/internal/gateway/rating/http"
	httphandler "sen1or/micromovie/movie/internal/handler/http"
	"sen1or/micromovie/pkg/discovery"
	"time"
)

const serviceName = "movie"

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "API handler port")
	flag.Parse()
	log.Printf("Starting the movie service on port %d", port)

	registry, err := discovery.NewConsulRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)

	if err := registry.Register(ctx, serviceName, instanceID, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.ReportHealthyState(serviceName, instanceID); err != nil {
				log.Println("failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, serviceName, instanceID)

	metadataGateway := metadatagateway.NewGateway(registry)
	ratingGateway := ratinggateway.NewGateway(registry)

	ctrl := movie.New(ratingGateway, metadataGateway)
	handler := httphandler.New(*ctrl)
	http.HandleFunc("/movie", handler.GetMovieDetails)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
