package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	dashboard Dashboard
}

type Dashboard struct {
	width   int
	height  int
	sidebar Sidebar
	main    Main
}

type Sidebar struct {
	width      int
	height     int
	namespaces []string
}

type Main struct {
	width  int
	height int
}

func initialModel() model {
	return model{
		Dashboard{
			sidebar: Sidebar{
				namespaces: getNamespaces(),
			},
			main: Main{},
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.dashboard.width = msg.Width
		m.dashboard.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (s Sidebar) View(width, height int) string {
	renderedNamespaces := []string{}
	for _, ns := range s.namespaces {
		renderedNamespaces = append(renderedNamespaces,
			lipgloss.NewStyle().
				Width(width-2).
				Align(lipgloss.Center).
				Render(ns),
		)
	}

	return lipgloss.NewStyle().
		Width(width - 2).
		Height(height - 2).
		Border(lipgloss.RoundedBorder()).
		Render(lipgloss.JoinVertical(lipgloss.Top, renderedNamespaces...))
}

func (m Main) View(width, height int) string {
	return lipgloss.NewStyle().
		Width(width - 2).
		Height(height - 2).
		Border(lipgloss.RoundedBorder()).
		Render("")
}

func (d Dashboard) View() string {
	renderedSidebar := d.sidebar.View(d.width*1/5, d.height)
	renderedMain := d.main.View(d.width*4/5, d.height)

	return lipgloss.JoinHorizontal(lipgloss.Left, renderedSidebar, renderedMain)
}

func (m model) View() string {
	return m.dashboard.View()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
