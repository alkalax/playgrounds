package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	dashboard Dashboard
}

type Pane int

const (
	SidebarPane Pane = iota
	MainPane
)

type Dashboard struct {
	width        int
	height       int
	selectedPane Pane
	sidebar      Sidebar
	main         Main
}

type Sidebar struct {
	width      int
	height     int
	namespaces []string
	index      int
}

type Main struct {
	width       int
	height      int
	pods        []string
	index       int
	logViewport viewport.Model
	logLines    []string
	podView     bool
}

func initialModel() model {
	m := model{
		Dashboard{
			sidebar: Sidebar{
				namespaces: getNamespaces(),
			},
			main: Main{
				podView:  true,
				logLines: []string{},
			},
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
		case "ctrl+c":
			return m, tea.Quit
		case "j", "down":
			if m.dashboard.selectedPane == SidebarPane {
				if m.dashboard.sidebar.index < len(m.dashboard.sidebar.namespaces)-1 {
					m.dashboard.sidebar.index++
					newNs := m.dashboard.sidebar.namespaces[m.dashboard.sidebar.index]
					m.dashboard.main.pods = getPods(newNs)
				}
			} else {
				if m.dashboard.main.podView && m.dashboard.main.index < len(m.dashboard.main.pods)-1 {
					m.dashboard.main.index++
				}
			}
		case "k", "up":
			if m.dashboard.selectedPane == SidebarPane {
				if m.dashboard.sidebar.index > 0 {
					m.dashboard.sidebar.index--
					newNs := m.dashboard.sidebar.namespaces[m.dashboard.sidebar.index]
					m.dashboard.main.pods = getPods(newNs)
				}
			} else {
				if m.dashboard.main.podView && m.dashboard.main.index > 0 {
					m.dashboard.main.index--
				}
			}
		case " ":
			if m.dashboard.selectedPane == SidebarPane {
				m.dashboard.selectedPane = MainPane
				m.dashboard.main.index = 0
			} else if m.dashboard.main.podView {
				namespace := m.dashboard.sidebar.namespaces[m.dashboard.sidebar.index]
				pod := m.dashboard.main.pods[m.dashboard.main.index]
				m.dashboard.main.logLines = getLogs(namespace, pod)
				//m.dashboard.main.logViewport = viewport.New(m.dashboard.width-2, m.dashboard.height-2)
				//m.dashboard.main.logViewport.SetContent(strings.Join(m.dashboard.main.logLines, "\n"))
				m.dashboard.main.podView = false
			}
		case "q":
			if m.dashboard.selectedPane == MainPane {
				if m.dashboard.main.podView {
					m.dashboard.selectedPane = SidebarPane
				} else {
					m.dashboard.main.podView = true
					m.dashboard.main.logLines = []string{}
				}
			}
		}
	}

	return m, nil
}

func (s Sidebar) View(width, height int, focused bool) string {
	renderedNamespaces := []string{}
	for i, ns := range s.namespaces {
		nsStyle := lipgloss.NewStyle().
			Width(width - 2).
			Align(lipgloss.Center)
		if focused && i == s.index {
			nsStyle = nsStyle.Background(lipgloss.Color("2"))
		}
		renderedNamespaces = append(renderedNamespaces, nsStyle.Render(ns))
	}

	style := lipgloss.NewStyle().
		Width(width - 2).
		Height(height - 2).
		Border(lipgloss.RoundedBorder())
	if focused {
		style = style.BorderForeground(lipgloss.Color("10"))
	}

	return style.Render(lipgloss.JoinVertical(lipgloss.Top, renderedNamespaces...))
}

func (m Main) View(width, height int, focused bool) string {
	style := lipgloss.NewStyle().
		Width(width - 2).
		Height(height - 2).
		Border(lipgloss.RoundedBorder())
	if focused {
		style = style.BorderForeground(lipgloss.Color("10"))
	}

	if m.podView {
		renderedPods := []string{}
		for i, pod := range m.pods {
			podStyle := lipgloss.NewStyle().
				Width(width - 2).
				Align(lipgloss.Center)
			if focused && i == m.index {
				podStyle = podStyle.Background(lipgloss.Color("2"))
			}
			renderedPods = append(renderedPods, podStyle.Render(pod))
		}

		return style.Render(lipgloss.JoinVertical(lipgloss.Top, renderedPods...))
	} else {
		m.logViewport = viewport.New(width-2, height-2)
		m.logViewport.SetContent(strings.Join(m.logLines, "\n"))
		m.logViewport.GotoBottom()
		return style.Render(m.logViewport.View())
	}
}

func (d Dashboard) View() string {
	renderedSidebar := d.sidebar.View(d.width*1/5, d.height, d.selectedPane == SidebarPane)
	renderedMain := d.main.View(d.width*4/5, d.height, d.selectedPane == MainPane)

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
