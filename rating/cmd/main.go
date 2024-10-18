package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sen1or/micromovie/gen"
	"sen1or/micromovie/pkg/discovery"
	"sen1or/micromovie/rating/internal/controller/rating"
	grpcHandler "sen1or/micromovie/rating/internal/handler/grpc"
	"sen1or/micromovie/rating/internal/repository/memory"
	"time"

	"google.golang.org/grpc"
)

const serviceName = "rating"

func main() {
	var port int
	flag.IntVar(&port, "p", 8082, "API Serve Port")
	flag.Parse()

	log.Printf("starting the rating service on port %d", port)

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

	repo := memory.NewRepository()
	ctrl := rating.New(repo)
	handler := grpcHandler.New(ctrl)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		panic(err)
	}

	sv := grpc.NewServer()
	gen.RegisterRatingServiceServer(sv, handler)
	sv.Serve(lis)
}
