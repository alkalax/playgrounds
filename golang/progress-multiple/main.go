package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	progressBlock progressBlock
}

type progressBlock struct {
	percent  float64
	progress progress.Model
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialModel() model {
	return model{
		progressBlock: progressBlock{
			progress: progress.New(),
		},
	}
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case " ":
			m.progressBlock.percent = 0.0
			return m, tick()
		}
	case tickMsg:
		m.progressBlock.percent += 0.0001
		if m.progressBlock.percent > 1.0 {
			m.progressBlock.percent = 1.0
			return m, nil
		}
		return m, tick()
	}

	return m, nil
}

func (pb progressBlock) View() string {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("135")).
		Foreground(lipgloss.Color("135")).
		Padding(1, 1).
		Render(pb.progress.ViewAs(pb.percent))
}

func (m model) View() string {
	return m.progressBlock.View()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
