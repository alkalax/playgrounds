package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	logs   []string
	stream chan LogMsg
}

type LogMsg struct {
	Line string
}

func startStream(ch chan<- LogMsg) tea.Cmd {
	return func() tea.Msg {
		resp, err := http.Get("http://127.0.0.1:8001/api/v1/namespaces/default/pods/logger/log?follow=true")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		reader := bufio.NewReader(resp.Body)
		for {
			line, _ := reader.ReadString('\n')
			ch <- LogMsg{Line: line}
			log.Println("Sent line")
		}
	}
}

func startSimpleStream(ch chan<- LogMsg) tea.Cmd {
	return func() tea.Msg {
		i := 0
		for {
			time.Sleep(1 * time.Second)
			ch <- LogMsg{Line: fmt.Sprintf("log line %d", i)}
			log.Println("Sent line")
			i++
		}
	}
}

func initialModel() model {
	return model{logs: []string{}, stream: make(chan LogMsg)}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(startStream(m.stream), waitForMsg(m.stream))
}

func waitForMsg(ch <-chan LogMsg) tea.Cmd {
	return func() tea.Msg {
		log.Println("waiting")
		return <-ch
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case LogMsg:
		log.Println("Read line")
		m.logs = append(m.logs, strings.Trim(msg.Line, "\n"))
		return m, waitForMsg(m.stream)
	}

	return m, nil
}

func (m model) View() string {
	log.Printf("Update: %d\n", len(m.logs))
	return lipgloss.JoinVertical(lipgloss.Top, m.logs...)
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
