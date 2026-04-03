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
	timeStr := time.Format("Mon 02 Jan 2006 - 15:03:05")
	fmt.Println(timeStr)

	tableFormatStr := "%-15s %-15s %-10s\n"
	tableHeading := fmt.Sprintf(tableFormatStr, "NAME", "JOB", "STATUS")

	fmt.Println(strings.Repeat("=", len(tableHeading)))
	fmt.Print(tableHeading)
	fmt.Println(strings.Repeat("=", len(tableHeading)))
	for _, citizen := range w.citizens {
		fmt.Printf(tableFormatStr, citizen.Name, citizen.Job, citizen.Status)
	}
}

func main() {
	world := NewWorld()

	for {
		world.Display(time.Now())
		time.Sleep(time.Second)
	}
}
