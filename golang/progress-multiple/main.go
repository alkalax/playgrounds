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
}

type progressBlock struct {
	percent  float64
	progress progress.Model
}

type tickMsg struct {
	index int
}

func tick(index int) tea.Cmd {
	return tea.Tick(time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg{index: index}
	})
}

func initialModel(n int) model {
	progressBlocks := make([]progressBlock, n)
	for i := range progressBlocks {
		progressBlocks[i] = progressBlock{progress: progress.New()}
	}
	return model{
		progressBlocks: progressBlocks,
	}
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for i := range m.progressBlocks {
		cmds = append(cmds, tick(i))
	}
	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
			//case " ":
			//	m.progressBlock.percent = 0.0
			//	return m, tick()
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
		m.progressBlocks[i].percent += 0.0001
		if m.progressBlocks[i].percent > 1.0 {
			m.progressBlocks[i].percent = 1.0
			return m, nil
		}
		return m, tick(i)
	}

	return m, nil
}

func (pb progressBlock) View(selected bool) string {
	borderColor := lipgloss.Color("135")
	if selected {
		borderColor = lipgloss.Color("250")
	}
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Foreground(lipgloss.Color("135")).
		Padding(1, 1).
		Render(pb.progress.ViewAs(pb.percent))
}

func (m model) View() string {
	renderedProgressBlocks := make([]string, len(m.progressBlocks))
	for i := range m.progressBlocks {
		renderedProgressBlocks[i] = m.progressBlocks[i].View(m.selected == i)
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
