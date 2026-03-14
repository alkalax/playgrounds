package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	progressBlocks []progressBlock
	selected       int
	width          int
	height         int
}

type progressBlock struct {
	label    string
	percent  float64
	progress progress.Model
	version  int
}

type tickMsg struct {
	index   int
	version int
}

func tick(index, version int) tea.Cmd {
	return tea.Tick(time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg{index: index, version: version}
	})
}

func initialModel(n int) model {
	progressBlocks := make([]progressBlock, n)
	for i := range progressBlocks {
		progressBlocks[i] = progressBlock{progress: progress.New(), label: fmt.Sprintf("Test %d", i)}
	}
	return model{
		progressBlocks: progressBlocks,
	}
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for i := range m.progressBlocks {
		cmds = append(cmds, tick(i, 0))
	}
	return tea.Batch(cmds...)
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
		case " ":
			label := m.progressBlocks[m.selected].label
			newVersion := m.progressBlocks[m.selected].version + 1
			m.progressBlocks[m.selected] = progressBlock{
				progress: progress.New(),
				label:    label,
				version:  newVersion,
			}

			return m, tick(m.selected, newVersion)
		case "j", "down":
			if m.selected < len(m.progressBlocks)-1 {
				m.selected++
			}
			return m, nil
		case "k", "up":
			if m.selected > 0 {
				m.selected--
			}
			return m, nil
		}
	case tickMsg:
		i := msg.index
		if m.progressBlocks[i].version != msg.version {
			return m, nil
		}

		m.progressBlocks[i].percent += 0.0001
		if m.progressBlocks[i].percent > 1.0 {
			m.progressBlocks[i].percent = 1.0
			return m, nil
		}
		return m, tick(i, msg.version)
	}

	return m, nil
}

func (pb progressBlock) View(width, height int, selected bool) string {
	borderColor := lipgloss.Color("135")
	if selected {
		borderColor = lipgloss.Color("250")
	}
	blockWidth := width - 2
	blockHeight := height - 2

	renderedLabel := lipgloss.NewStyle().
		Width(blockWidth - 1).
		Height(blockHeight/2 - 1).
		MarginBottom(1).
		Render(pb.label)

	renderedBar := lipgloss.NewStyle().
		Width(blockWidth - 1).
		Height(blockHeight/2 - 1).
		Render(pb.progress.ViewAs(pb.percent))

	return lipgloss.NewStyle().
		Width(blockWidth).
		Height(blockHeight).
		Align(lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Foreground(lipgloss.Color("135")).
		Padding(1, 1).
		Render(lipgloss.JoinVertical(lipgloss.Top, renderedLabel, renderedBar))
}

func (m model) View() string {
	renderedProgressBlocks := make([]string, len(m.progressBlocks))
	for i := range m.progressBlocks {
		renderedProgressBlocks[i] = m.progressBlocks[i].View(min(50, m.width), min(5, m.height/len(m.progressBlocks)), m.selected == i)
	}
	return lipgloss.JoinVertical(lipgloss.Top, renderedProgressBlocks...)
}

func main() {
	p := tea.NewProgram(initialModel(4), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
