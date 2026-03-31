package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func worker(ctx context.Context, id int, jobs <-chan int, wg *sync.WaitGroup) {
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				fmt.Printf("worker %d: channel is closed, exiting\n", id)
				return
			}

			fmt.Printf("worker %d working on job %d\n", id, job)
			time.Sleep(time.Duration(rand.Intn(5)) * time.Second)

			wg.Done()

		case <-ctx.Done():
			fmt.Printf("worker %d: received cancellation\n", id)
			return
		}
	}
}

func main() {
	numWorkers := 3
	numJobs := 30

	jobs := make(chan int, numJobs)
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	for i := range numWorkers {
		go worker(ctx, i+1, jobs, &wg)
	}

	go func() {
		for j := range numJobs {
			wg.Add(1)
			jobs <- j + 1
		}
	}()
	// allow time for jobs to start before main exiting
	time.Sleep(time.Second)

	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("===== CANCELLING =====")
		cancel()
	}()

	wg.Wait()
	fmt.Println("All jobs are done.")
}
