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
		Width(width-2).
		Height(height-2).
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.NormalBorder()).
		Render(h.content)
}

func (m mainDiv) View(width, height int) string {
	sidebarWidth := width / 5
	sidebar := lipgloss.NewStyle().
		Width(sidebarWidth-2).
		Height(height-2).
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.NormalBorder()).
		Render("sidebar")

	mainWidth := width - sidebarWidth
	main := lipgloss.NewStyle().
		Width(mainWidth-2).
		Height(height-2).
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.NormalBorder()).
		Render(m.content)

	return lipgloss.JoinHorizontal(lipgloss.Center, sidebar, main)
}

func (m model) View() string {
	header := m.header.View(m.width, m.height/8)
	main := m.main.View(m.width, m.height*7/8)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Top,
		lipgloss.JoinVertical(lipgloss.Top, header, main),
	)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
