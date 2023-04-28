package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "tt/pkg/api"
)

type server struct {
	pb.UnimplementedApiServer
}

const (
	grcpPort string = ":9090"
	httpPort string = ":9091"
)

func runHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterApiHandlerFromEndpoint(ctx, mux, grcpPort, opts)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("HTTP server: %s\n", httpPort)
	return http.ListenAndServe("localhost"+httpPort, mux)
}

func runGRPCServer() error {
	lis, err := net.Listen("tcp", grcpPort)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterApiServer(s, &server{})

	log.Printf("gRPC server: %s\n", grcpPort)
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

func (s *server) CalculateSum(_ context.Context, req *pb.Request) (*pb.Response, error) { // доделать
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
	return &pb.Response{Sum: int64(total)}, nil
}
