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

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	s := grpc.NewServer()

	pb.RegisterApiServer(s, &server{})
	log.Printf("Serving gRPC on %s:%s", "localhost", "8080")
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:8080",
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	err = pb.RegisterApiHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", "8090"),
		Handler: gwmux,
	}

	log.Println(fmt.Sprintf("Serving gRPC-Gateway on %s:%s", "localhost", "8090"))
	log.Fatalln(gwServer.ListenAndServe())
}

func (s *server) CalculateSum(_ context.Context, req *pb.Request) (*pb.Response, error) {
	return &pb.Response{Sum: service.Calculate(req.Numbers)}, nil
}
