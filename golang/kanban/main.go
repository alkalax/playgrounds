package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	width   int
	height  int
	columns []Column
}

type Column struct {
	title string
	tasks []string
}

func initialModel() tea.Model {
	return model{
		columns: []Column{
			{
				title: "TODO",
				tasks: []string{
					"Task 1",
					"Task 2",
					"Task 3",
				},
			},
			{
				title: "IN PROGRESS",
				tasks: []string{
					"Task 4",
					"Task 5",
				},
			},
			{
				title: "DONE",
				tasks: []string{
					"Task 6",
				},
			},
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (c Column) View(width, height int) string {
	renderedTasks := make([]string, len(c.tasks)+1) // +1 for column title
	renderedTasks[0] = lipgloss.NewStyle().
		Width(width * 3 / 4).
		Align(lipgloss.Center).
		MarginBottom(1).
		Border(lipgloss.HiddenBorder()).
		BorderForeground(lipgloss.Color("23")).
		Foreground(lipgloss.Color("23")).
		Bold(true).
		Render(c.title)

	for i, t := range c.tasks {
		renderedTasks[i+1] = lipgloss.NewStyle().
			Width(width*3/4).
			Align(lipgloss.Center).
			Padding(0, 1).
			Border(lipgloss.NormalBorder()).
			Foreground(lipgloss.Color("23")).
			BorderForeground(lipgloss.Color("23")).
			Render(t)
	}

	return lipgloss.NewStyle().
		Height(height).
		Width(width).
		Align(lipgloss.Center, lipgloss.Top).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("23")).
		Render(lipgloss.JoinVertical(lipgloss.Top, renderedTasks...))
}

func (m model) View() string {
	colWidth := m.width / 4
	colHeight := m.height * 3 / 4

	style := lipgloss.NewStyle()
	row := lipgloss.JoinHorizontal(
		lipgloss.Center,
		style.MarginRight(3).Render(m.columns[0].View(colWidth, colHeight)),
		style.MarginRight(3).Render(m.columns[1].View(colWidth, colHeight)),
		style.Render(m.columns[2].View(colWidth, colHeight)),
	)

	return lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(row)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
