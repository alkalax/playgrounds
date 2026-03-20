package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	currentDir string
	contents   dirContent
	focused    int
}

type dirContent struct {
	entries []os.DirEntry
}

func initialModel() model {
	entries, err := os.ReadDir("/")
	if err != nil {
		panic(err)
	}
	return model{
		currentDir: "/",
		contents:   dirContent{entries},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func getEntries(path string) []string {
	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	renderedEntries := []string{}
	for _, entry := range entries {
		color := lipgloss.Color("255")
		if entry.IsDir() {
			color = lipgloss.Color("25")
		}
		renderedEntries = append(renderedEntries, lipgloss.NewStyle().Foreground(color).Render(entry.Name()))
	}

	return renderedEntries
}

func (dc *dirContent) View(focused int) string {
	var sb strings.Builder
	for i, entry := range dc.entries {
		color := lipgloss.Color("255")
		if entry.IsDir() {
			color = lipgloss.Color("25")
		}
		entryStyle := lipgloss.NewStyle().Foreground(color)
		if focused == i {
			entryStyle = entryStyle.Background(lipgloss.Color("1"))
		}
		sb.WriteString(entryStyle.Render(entry.Name()))
		sb.WriteString("\n")
	}

	return sb.String()
}

func (m model) View() string {
	//var sb strings.Builder
	//sb.WriteString(m.currentDir)
	//sb.WriteString("\n\n")
	//for _, entry := range getEntries(m.currentDir) {
	//	sb.WriteString(entry)
	//	sb.WriteString("\n")
	//}

	//return sb.String()
	return fmt.Sprintf("\t%s\n\n%s", m.currentDir, m.contents.View(m.focused))
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
