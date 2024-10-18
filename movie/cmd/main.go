package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sen1or/micromovie/gen"
	"sen1or/micromovie/movie/internal/controller/movie"
	metadatagateway "sen1or/micromovie/movie/internal/gateway/metadata/grpc"
	ratinggateway "sen1or/micromovie/movie/internal/gateway/rating/grpc"
	grpchandler "sen1or/micromovie/movie/internal/handler/grpc"
	"sen1or/micromovie/pkg/discovery"
	"time"

	"google.golang.org/grpc"
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

	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)
	ctrl := movie.New(ratingGateway, metadataGateway)
	handler := grpchandler.New(*ctrl)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	sv := grpc.NewServer()
	gen.RegisterMovieServiceServer(sv, handler)
	sv.Serve(lis)
}
