package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"sen1or/micromovie/metadata/internal/controller/metadata"
	httpHandler "sen1or/micromovie/metadata/internal/handler/http"
	"sen1or/micromovie/metadata/internal/repository"
	"sen1or/micromovie/pkg/discovery"
)

const serviceName = "metadata"

func main() {
	var port int
	flag.IntVar(&port, "port", 8081, "API Serve Port")
	flag.Parse()

	log.Printf("starting the metadata service on port %d", port)

	registry, err := discovery.NewConsulRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, serviceName, instanceID, fmt.Sprintf("%s:%d", "localhost", port)); err != nil {
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

	repo := repository.NewRepository()
	ctrl := metadata.NewController(repo)
	handler := httpHandler.NewHandler(ctrl)

	http.HandleFunc("/metadata", handler.GetMetadata)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Printf("error while running metadata server: %s\n", err)
	}
}
