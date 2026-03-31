package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	for job := range jobs {
		fmt.Printf("worker %d working on job %d\n", id, job)
		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)

		wg.Done()
	}
}

func main() {
	numWorkers := 3
	numJobs := 30

	jobs := make(chan int)
	wg := sync.WaitGroup{}

	for i := range numWorkers {
		go worker(i+1, jobs, &wg)
	}

	for j := range numJobs {
		wg.Add(1)
		jobs <- j + 1
	}

	wg.Wait()
	fmt.Println("All jobs are done.")
}
