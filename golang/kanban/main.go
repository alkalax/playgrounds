package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	colorNormal  = lipgloss.Color("23")
	colorFocused = lipgloss.Color("25")
)

type model struct {
	width         int
	height        int
	focusedColumn int
	columns       []Column
}

type Column struct {
	title       string
	tasks       []string
	focusedTask int
}

func (c Column) changeTaskPriority(increase bool) (changed bool) {
	changed = false
	if !increase && c.focusedTask < len(c.tasks)-1 {
		tmp := c.tasks[c.focusedTask+1]
		c.tasks[c.focusedTask+1] = c.tasks[c.focusedTask]
		c.tasks[c.focusedTask] = tmp

		changed = true

	} else if increase && c.focusedTask > 0 {
		tmp := c.tasks[c.focusedTask-1]
		c.tasks[c.focusedTask-1] = c.tasks[c.focusedTask]
		c.tasks[c.focusedTask] = tmp

		changed = true
	}

	return changed
}

func (c *Column) moveTask(index int, dest *Column) {
	dest.tasks = append(dest.tasks, c.tasks[index])
	c.tasks = append(c.tasks[:c.focusedTask], c.tasks[c.focusedTask+1:]...)
	if c.focusedTask > 0 {
		c.focusedTask--
	}
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
		case "h", "left":
			if m.focusedColumn > 0 {
				m.focusedColumn--
			}
			return m, nil
		case "l", "right":
			if m.focusedColumn < len(m.columns)-1 {
				m.focusedColumn++
			}
			return m, nil
		case "j", "down":
			focusedColumn := m.columns[m.focusedColumn]
			if focusedColumn.focusedTask < len(focusedColumn.tasks)-1 {
				m.columns[m.focusedColumn].focusedTask++
			}
			return m, nil
		case "k", "up":
			focusedColumn := m.columns[m.focusedColumn]
			if focusedColumn.focusedTask > 0 {
				m.columns[m.focusedColumn].focusedTask--
			}
			return m, nil
		case "ctrl+j", "ctrl+down":
			if m.columns[m.focusedColumn].changeTaskPriority(false) {
				m.columns[m.focusedColumn].focusedTask++
			}
			return m, nil
		case "ctrl+k", "ctrl+up":
			if m.columns[m.focusedColumn].changeTaskPriority(true) {
				m.columns[m.focusedColumn].focusedTask--
			}
			return m, nil
		case "ctrl+l", "ctrl+right":
			if m.focusedColumn < len(m.columns)-1 && len(m.columns[m.focusedColumn].tasks) > 0 {
				currColumn := &m.columns[m.focusedColumn]
				nextColumn := &m.columns[m.focusedColumn+1]
				currColumn.moveTask(currColumn.focusedTask, nextColumn)
				nextColumn.focusedTask = len(nextColumn.tasks) - 1
				m.focusedColumn++
			}
			return m, nil
		case "ctrl+h", "ctrl+left":
			if m.focusedColumn > 0 && len(m.columns[m.focusedColumn].tasks) > 0 {
				currColumn := &m.columns[m.focusedColumn]
				prevColumn := &m.columns[m.focusedColumn-1]
				currColumn.moveTask(currColumn.focusedTask, prevColumn)
				prevColumn.focusedTask = len(prevColumn.tasks) - 1
				m.focusedColumn--
			}
			return m, nil
		}
	}

	return m, nil
}

func (c Column) View(width, height int, focused bool) string {
	color := colorNormal
	if focused {
		color = colorFocused
	}
	renderedTasks := make([]string, len(c.tasks)+1) // +1 for column title
	renderedTasks[0] = lipgloss.NewStyle().
		Width(width * 3 / 4).
		Align(lipgloss.Center).
		MarginBottom(1).
		Border(lipgloss.HiddenBorder()).
		BorderForeground(color).
		Foreground(color).
		Bold(true).
		Render(c.title)

	for i, t := range c.tasks {
		taskColor := colorNormal
		if focused && i == c.focusedTask {
			taskColor = colorFocused
		}
		renderedTasks[i+1] = lipgloss.NewStyle().
			Width(width*3/4).
			Align(lipgloss.Center).
			Padding(0, 1).
			Border(lipgloss.NormalBorder()).
			BorderForeground(taskColor).
			Render(t)
	}

	return lipgloss.NewStyle().
		Height(height).
		Width(width).
		Align(lipgloss.Center, lipgloss.Top).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(color).
		Render(lipgloss.JoinVertical(lipgloss.Top, renderedTasks...))
}

func (m model) View() string {
	colWidth := m.width / 4
	colHeight := m.height * 3 / 4

	style := lipgloss.NewStyle()
	row := lipgloss.JoinHorizontal(
		lipgloss.Center,
		style.MarginRight(3).Render(m.columns[0].View(colWidth, colHeight, m.focusedColumn == 0)),
		style.MarginRight(3).Render(m.columns[1].View(colWidth, colHeight, m.focusedColumn == 1)),
		style.Render(m.columns[2].View(colWidth, colHeight, m.focusedColumn == 2)),
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
