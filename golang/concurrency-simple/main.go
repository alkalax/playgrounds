package main

import (
	"fmt"
	"math/rand"
	"time"
)

func worker(id int, ch chan string) {
	seconds := rand.Intn(5 + id)
	time.Sleep(time.Duration(seconds) * time.Second)
	ch <- fmt.Sprintf("worker %d done after %ds", id, seconds)
}

func main() {
	ch := make(chan string)

	n := 5
	for i := 1; i <= n; i++ {
		go worker(i, ch)
	}

	for i := 1; i <= n; i++ {
		msg := <-ch
		fmt.Println(msg)
	}
}
