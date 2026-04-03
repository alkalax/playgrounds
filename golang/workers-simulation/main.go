package main

import (
	"fmt"
	"time"
)

type World struct {
	workers []int
}

func NewWorld(workers []int) *World {
	return &World{
		workers: workers,
	}
}

func (w *World) Display() {
	fmt.Print("\033[H\033[2J") // clear terminal
	for _, worker := range w.workers {
		fmt.Printf("worker %d: idle\n", worker)
	}
}

func main() {
	world := NewWorld([]int{1, 2, 3})

	for {
		world.Display()
		time.Sleep(time.Second)
	}
}
