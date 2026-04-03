package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var sample string = "This is a sample string. This is also a different sentence altogether."

type Model struct {
	content string
	width   int
	height  int
	padding int
}

func initialModel() tea.Model {
	return Model{
		content: "test",
		padding: 1,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	return lipgloss.NewStyle().
		Width(m.width-2).
		Height(m.height-2).
		Padding(m.padding, m.padding).
		Border(lipgloss.NormalBorder()).
		Render(m.content)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
