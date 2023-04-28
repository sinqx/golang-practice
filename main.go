package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "tt/pkg/api"
)

type server struct {
	pb.UnimplementedApiServer
}

const (
	grcpPort string = "9090"
	httpPort string = "9091"
)

func runHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterApiHandlerFromEndpoint(ctx, mux, grcpPort, opts)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("HTTP server: %s\n", httpPort)
	return http.ListenAndServe(httpPort, mux)
}

func runGRPCServer() error {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterApiServer(s, &server{})

	log.Printf("gRPC server: %s\n", ":9090")
	return s.Serve(lis)
}

func main() {
	go func() {
		if err := runGRPCServer(); err != nil {
			log.Fatalf("gRPC Error: %v", err)
		}
	}()

	ctx := context.Background()
	if err := runHTTPServer(ctx); err != nil {
		log.Fatalf("HTTP Error: %v", err)
	}
}

func (s *server) CalculateSum(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	var sum int64
	for _, num := range req.Numbers {
		sum += num.Nums
	}

	return &pb.Response{Sum: sum}, nil
}
