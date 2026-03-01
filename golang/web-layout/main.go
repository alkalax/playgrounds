package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	width  int
	height int
	header header
	main   mainDiv
}

type header struct {
	content string
	width   int
	height  int
}

type mainDiv struct {
	content string
	width   int
	height  int
}

func initialModel() model {
	return model{
		header: header{content: "Title Here"},
		main:   mainDiv{content: "Some text here"},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (h header) View(width, height int) string {
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.NormalBorder()).
		Render(h.content)
}

func (m model) View() string {
	header := m.header.View(m.width-2, m.height/8)
	main := m.main.content

	return lipgloss.JoinVertical(lipgloss.Top, header, main)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
