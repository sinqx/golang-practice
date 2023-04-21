package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"tt/pkg/api"
)

func main() {

	conn, err := grpc.Dial("my_server:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer conn.Close()

	c := api.NewApiClient(conn)

	req := &api.Request{}
	for i := 1; i <= 100000; i++ {
		req.Numbers = append(req.Numbers, &api.Numbers{Nums: int64(i)})
	}

	res, err := c.CalculateSum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Sum: %d\n", res.Sum)
}
