package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var sample string = "This is a sample string. This is also a different sentence altogether."

type Model struct {
	tokenField TokenField
	width      int
	height     int
}

type TokenField struct {
	tokens  []Token
	width   int
	height  int
	padding int
}

type Token struct {
	id    int
	word  string
	start int
	end   int
	line  int
}

func initialModel() *Model {
	tokens := []Token{}
	for i, word := range strings.Split(sample, " ") {
		tokens = append(tokens, Token{
			id:   i,
			word: word,
		})
	}
	return &Model{
		tokenField: TokenField{
			tokens:  tokens,
			padding: 1,
		},
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (tf *TokenField) View(width, height int) string {
	tokenFieldWidth := width - 2*tf.padding - 2
	tokenFieldHeight := height - 2*tf.padding - 2

	var sb strings.Builder
	for i, token := range tf.tokens {
		sb.WriteString(token.word)
		if i < len(tf.tokens) {
			sb.WriteRune(' ')
		}
	}

	return lipgloss.NewStyle().
		Width(tokenFieldWidth).
		Height(tokenFieldHeight).
		Padding(tf.padding, tf.padding).
		Border(lipgloss.NormalBorder()).
		Render(sb.String())
}

func (m *Model) View() string {
	return m.tokenField.View(m.width, m.height)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
