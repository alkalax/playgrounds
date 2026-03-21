package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	currentDir string
	focused    int
	entries    []os.DirEntry
	viewport   viewport.Model
}

func getDirContent(path string) []os.DirEntry {
	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	return entries
}

func initialModel() model {
	m := model{currentDir: "/"}
	m.entries = getDirContent(m.currentDir)
	m.viewport = viewport.New(0, 0)

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 4
		m.viewport.SetContent(m.renderEntries())
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "j", "down":
			if m.focused < len(m.entries)-1 {
				m.focused++
				m.viewport.SetContent(m.renderEntries())
			}
		case "k", "up":
			if m.focused > 0 {
				m.focused--
				m.viewport.SetContent(m.renderEntries())
			}
		case " ":
			focusedEntry := m.entries[m.focused]
			if focusedEntry.IsDir() {
				m.currentDir = filepath.Join(m.currentDir, focusedEntry.Name())
				m.entries = getDirContent(m.currentDir)
				m.focused = 0
				m.viewport.SetContent(m.renderEntries())
				m.viewport.GotoTop()
			}
		case "q":
			if m.currentDir != "/" {
				m.currentDir = filepath.Dir(m.currentDir)
				m.entries = getDirContent(m.currentDir)
				m.focused = 0
				m.viewport.SetContent(m.renderEntries())
				m.viewport.GotoTop()
			}
		}

		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
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

func (m model) renderEntries() string {
	var sb strings.Builder
	for i, entry := range m.entries {
		color := lipgloss.Color("255")
		if entry.IsDir() {
			color = lipgloss.Color("25")
		}
		entryStyle := lipgloss.NewStyle().Foreground(color)
		if m.focused == i {
			entryStyle = entryStyle.Background(lipgloss.Color("2"))
		}
		sb.WriteString(entryStyle.Render(entry.Name()))
		sb.WriteString("\n")
	}

	return sb.String()
}

func (m model) View() string {
	return fmt.Sprintf("\t%s\n\n%s", m.currentDir, m.viewport.View())
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
