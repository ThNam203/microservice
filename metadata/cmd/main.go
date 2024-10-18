package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"sen1or/micromovie/gen"
	"sen1or/micromovie/metadata/internal/controller/metadata"
	grpcHandler "sen1or/micromovie/metadata/internal/handler/grpc"
	"sen1or/micromovie/metadata/internal/repository"
	"sen1or/micromovie/pkg/discovery"

	"google.golang.org/grpc"
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
	handler := grpcHandler.New(ctrl)

	sv := grpc.NewServer()
	gen.RegisterMetadataServiceServer(sv, handler)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		panic(err)
	}
	sv.Serve(lis)
}
