package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	. "testTask/pkg/api"
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

func main() {
	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	s := grpc.NewServer()
	RegisterApiServer(s, &server{})

	log.Printf("Port: %v", listen.Addr())
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
