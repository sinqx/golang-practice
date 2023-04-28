// ///////////////////////////// deprecated ///////////////////////////////
package main

import "sync"

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	. "testTask/pkg/api"
	"tt/pkg/api"
)

type server struct {
	UnimplementedApiServer
}

func (s *server) CalculateSum(ctx context.Context, req *Request) (*Response, error) {
	var wg sync.WaitGroup
	parts := 10                          // кол-во потоков
	partSize := len(req.Numbers) / parts // кол-во чисел в потоке
	results := make(chan int, parts)     // канал для получения результата
	for i := 0; i < parts; i++ {
		wg.Add(1)
		start := i * partSize   // начало слайса
		end := start + partSize // конец слайса
		go func() {
			defer wg.Done()
			sum := 0
			for _, num := range req.Numbers[start:end] {
				sum += int(num.Nums)
			}
			results <- sum
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	total := 0
	for res := range results {
		total += res
	}
	return &Response{Sum: int64(total)}, nil
}

const (
	grpcPort = ":50051"
	httpPort = ":8080"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Start a gRPC server
	grpcServer := grpc.NewServer()
	RegisterApiServer(grpcServer, &server{})
	reflection.Register(grpcServer)

	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	s := grpc.NewServer()
	RegisterApiServer(s, &server{})

	log.Printf("Port: %v", listen.Addr())
	go func() {
		log.Fatalln(s.Serve(listen))
	}()

	// Start a HTTP server
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = api.RegisterApiHandlerFromEndpoint(ctx, mux, grpcPort, opts)
	if err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}

	log.Printf("Starting HTTP server on %s", httpPort)
	if err := http.ListenAndServe(httpPort, mux); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
