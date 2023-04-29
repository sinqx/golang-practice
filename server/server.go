// ///////////////////////////// client/server ///////////////////////////////
package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	grpcPort = ":8080"
	httpPort = ":8090"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Start a gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterApiServer(grpcServer, &server{})
	reflection.Register(grpcServer)

	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterApiServer(s, &server{})

	log.Printf("Port: %v", listen.Addr())
	go func() {
		log.Fatalln(s.Serve(listen))
	}()

	// Start a HTTP server
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = pb.RegisterApiHandlerFromEndpoint(ctx, mux, grpcPort, opts)
	if err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}

	log.Printf("Starting HTTP server on %s", httpPort)
	if err := http.ListenAndServe(httpPort, mux); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

func (s *server) CalculateSum(_ context.Context, req *pb.Request) (*pb.Response, error) { // доделать
	return &pb.Response{Sum: service.Calculate(req.Numbers)}, nil
}
