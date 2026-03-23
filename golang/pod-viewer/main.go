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
	focused    int
}

type Main struct {
	width  int
	height int
	pods   []string
}

func initialModel() model {
	m := model{
		Dashboard{
			sidebar: Sidebar{
				namespaces: getNamespaces(),
			},
			main: Main{},
		},
	}

	m.dashboard.main.pods = getPods(m.dashboard.sidebar.namespaces[0])

	return m
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
		case "j", "down":
			if m.dashboard.sidebar.focused < len(m.dashboard.sidebar.namespaces)-1 {
				m.dashboard.sidebar.focused++
				newNs := m.dashboard.sidebar.namespaces[m.dashboard.sidebar.focused]
				m.dashboard.main.pods = getPods(newNs)
			}
		case "k", "up":
			if m.dashboard.sidebar.focused > 0 {
				m.dashboard.sidebar.focused--
				newNs := m.dashboard.sidebar.namespaces[m.dashboard.sidebar.focused]
				m.dashboard.main.pods = getPods(newNs)
			}
		}
	}

	return m, nil
}

func (s Sidebar) View(width, height int) string {
	renderedNamespaces := []string{}
	for i, ns := range s.namespaces {
		style := lipgloss.NewStyle().
			Width(width - 2).
			Align(lipgloss.Center)
		if i == s.focused {
			style = style.Background(lipgloss.Color("2"))
		}
		renderedNamespaces = append(renderedNamespaces, style.Render(ns))
	}

	return lipgloss.NewStyle().
		Width(width - 2).
		Height(height - 2).
		Border(lipgloss.RoundedBorder()).
		Render(lipgloss.JoinVertical(lipgloss.Top, renderedNamespaces...))
}

func (m Main) View(width, height int) string {
	renderedPods := []string{}
	for _, pod := range m.pods {
		renderedPods = append(renderedPods,
			lipgloss.NewStyle().
				Width(width-2).
				Align(lipgloss.Center).
				Render(pod),
		)
	}

	return lipgloss.NewStyle().
		Width(width - 2).
		Height(height - 2).
		Border(lipgloss.RoundedBorder()).
		Render(lipgloss.JoinVertical(lipgloss.Top, renderedPods...))
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
