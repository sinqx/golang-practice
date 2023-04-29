package service

import (
	"sync"
	pb "tt/pkg/api"
)

func Calculate(numbers []*pb.Numbers) int64 {
	var wg sync.WaitGroup

	input := make(chan int)
	results := make(chan int)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, num := range numbers {
			input <- int(num.Nums)
		}
		close(input)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var sum int
		for n := range input {
			sum += n
		}
		results <- sum
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	return int64(<-results)
}
