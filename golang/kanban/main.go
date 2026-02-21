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
	colStyle := lipgloss.NewStyle().
		Height(height).
		Width(width).
		Align(lipgloss.Center).
		Border(lipgloss.NormalBorder())

	return colStyle.Render(c.title)
}

func (m model) View() string {
	renderedColumns := make([]string, len(m.columns))
	for i := range m.columns {
		renderedColumns[i] = m.columns[i].View(m.width/4, m.height*3/4)
	}

	row := lipgloss.JoinHorizontal(lipgloss.Center, renderedColumns...)

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
