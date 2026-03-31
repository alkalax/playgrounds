package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func worker(id int, ch chan string) {
	seconds := rand.Intn(5 + id)
	time.Sleep(time.Duration(seconds) * time.Second)
	ch <- fmt.Sprintf("worker %d done after %ds", id, seconds)
}

func workerWg(id int, jobs <-chan int, wg *sync.WaitGroup) {
	for job := range jobs {
		fmt.Printf("worker %d processing job %d\n", id, job)
		time.Sleep(time.Second * time.Duration(rand.Intn(5)))

		wg.Done()
	}
}

func main() {
	//ch := make(chan string)

	//n := 5
	//for i := 1; i <= n; i++ {
	//	go worker(i, ch)
	//}

	//for i := 1; i <= n; i++ {
	//	msg := <-ch
	//	fmt.Println(msg)
	//}

	ch := make(chan int)
	//jobs := []int{1, 2, 3, 4, 5, 6, 7}
	//for _, job := range jobs {
	//	go func(j int) {
	//		time.Sleep(time.Duration(j) * time.Second)
	//		ch <- j * 2
	//	}(job)
	//}
	wg := sync.WaitGroup{}
	numWorkers := 3
	for i := range numWorkers {
		go workerWg(i+1, ch, &wg)
	}

	numJobs := 30
	for job := range numJobs {
		wg.Add(1)
		ch <- job + 1
	}
	close(ch)
	//for i := 0; i < len(jobs); i++ {
	//	result := <-ch
	//	fmt.Println(result)
	//}

	wg.Wait()
	fmt.Println("All jobs are done")
}
