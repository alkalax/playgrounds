package main

import (
	"fmt"
	"strings"
	"time"
)

type World struct {
	citizens []*Citizen
}

type Citizen struct {
	Id     int
	Name   string
	Job    string
	Status string
}

func NewWorld() *World {
	citizens := []*Citizen{
		{
			Id:     1,
			Name:   "Amy",
			Job:    "Laborer",
			Status: "working",
		},
		{
			Id:     2,
			Name:   "John",
			Job:    "Laborer",
			Status: "working",
		},
		{
			Id:     3,
			Name:   "Robert",
			Job:    "Laborer",
			Status: "idle",
		},
	}

	return &World{
		citizens: citizens,
	}
}

func (w *World) Display(time time.Time) {
	fmt.Print("\033[H\033[2J") // clear terminal
	timeStr := time.Format("02.01.2006 15:04:05")
	fmt.Println(timeStr)
	fmt.Println(strings.Repeat("=", len(timeStr)))

	for _, citizen := range w.citizens {
		fmt.Printf("%s %s: %s\n", citizen.Job, citizen.Name, citizen.Status)
	}
}

func main() {
	world := NewWorld()

	for {
		world.Display(time.Now())
		time.Sleep(time.Second)
	}
}
