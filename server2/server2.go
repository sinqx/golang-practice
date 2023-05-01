package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	pb "tt/pkg/api"
	"tt/service"
)

type server struct {
	pb.UnimplementedApiServer
}

const (
	grpcPort string = ":8080"
	httpPort string = ":8090"
)

func runHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterApiHandlerFromEndpoint(ctx, mux, grpcPort, opts)
	if err != nil {
		log.Fatalf("Starting HTTP server Error: %v", err)
	}

	log.Printf("HTTP server: %s\n", httpPort)
	return http.ListenAndServe(httpPort, mux)
}

func runGRPCServer() error {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("Starting GRPC server Error: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterApiServer(s, &server{})

	log.Printf("gRPC server: %s\n", grpcPort)
	return s.Serve(lis)
}

func main() {
	go func() {
		if err := runGRPCServer(); err != nil {
			log.Fatalf("GRCP error: %v", err)
		}
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:8080",
		grpc.WithBlock(),
		grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	gwmux := runtime.NewServeMux()
	// Register User Service
	err = pb.RegisterApiHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", httpPort),
		Handler: gwmux,
	}
	ctx := context.Background()
	if err := runHTTPServer(ctx); err != nil {
		log.Fatalf("HTTP Error: %v", err)
	}
	log.Fatalln(gwServer.ListenAndServe())
}

func (s *server) CalculateSum(_ context.Context, req *pb.Request) (*pb.Response, error) {
	return &pb.Response{Sum: service.Calculate(req.Numbers)}, nil
}
