// ///////////////////////////// deprecated ///////////////////////////////
package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"
	"tt/pkg/api"
)

func main() {

	conn, err := grpc.DialContext(context.Background(), "localhost:8080",
		grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer conn.Close()

	c := api.NewApiClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &api.Request{}
	for i := 1; i <= 100000; i++ {
		req.Numbers = append(req.Numbers, &api.Numbers{Nums: int64(i)})
	}

	res, err := c.CalculateSum(ctx, req)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Sum: %d\n", res.Sum)
}
