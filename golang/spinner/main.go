package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	items []item
}

type item struct {
	name    string
	ready   bool
	spinner spinner.Model
}

func initialModel() model {
	items := make([]item, 5)

	for i := range items {
		items[i].name = fmt.Sprintf("items%03d", i)

		items[i].spinner = spinner.New()
		items[i].spinner.Spinner = spinner.Points
	}

	return model{items: items}
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for i := range m.items {
		cmds = append(cmds, m.items[i].spinner.Tick)
	}

	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	default:
		cmds := make([]tea.Cmd, len(m.items))
		for i := range m.items {
			m.items[i].spinner, cmds[i] = m.items[i].spinner.Update(msg)
		}
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (i item) View() string {
	return fmt.Sprintf("\t%s\t%s", i.name, i.spinner.View())
}

func (m model) View() string {
	var sb strings.Builder
	for i := range m.items {
		sb.WriteString(m.items[i].View())
		sb.WriteString("\n")
	}

	return sb.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
