package service

import (
	"runtime"
	"sync"
	pb "tt/pkg/api"
)

func Calculate(numbers []*pb.Numbers) int64 {
	var sum int64
	var wg sync.WaitGroup
	var mu sync.Mutex

	numGoroutines := runtime.NumCPU() //количество горутин

	chunkSize := len(numbers) / numGoroutines
	chunks := make([][]*pb.Numbers, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := start + chunkSize

		if i == numGoroutines-1 {
			end = len(numbers)
		}
		chunks[i] = numbers[start:end]
	}

	// Запустить горутины
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(nums []*pb.Numbers) {
			defer wg.Done()

			var localSum int64

			for _, num := range nums {
				localSum += num.Nums
			}

			mu.Lock()
			sum += localSum
			mu.Unlock()

		}(chunks[i])
	}

	wg.Wait()
	return sum
}

//func Calculate(numbers []*pb.Numbers) int64 {
//	var sum int64 = 0
//	for _, num := range numbers {
//		sum += num.Nums
//	}
//	return sum
//}
