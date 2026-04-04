package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var sample string = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

type Styles struct {
	focusedToken lipgloss.Style
}

type Model struct {
	tokenField TokenField
	width      int
	height     int
	index      int
}

type TokenField struct {
	tokens            []Token
	width             int
	height            int
	horizontalPadding int
	verticalPadding   int
}

type Token struct {
	sep   bool
	word  string
	start int
	end   int
	line  int
	index int
}

func defaultStyles() Styles {
	return Styles{
		focusedToken: lipgloss.NewStyle().Bold(true).Underline(true).Foreground(lipgloss.Color("10")),
	}
}

func tokenize(text string) []Token {
	tokens := []Token{}
	for word := range strings.SplitSeq(text, " ") {
		tokens = append(tokens, Token{word: word})
	}

	return tokens
}

func initialModel() *Model {
	return &Model{
		tokenField: TokenField{
			tokens:            tokenize(sample),
			horizontalPadding: 1,
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
		case "l", "right":
			if m.index < len(m.tokenField.tokens)-1 {
				m.index++
			}
		case "h", "left":
			if m.index > 0 {
				m.index--
			}
		case "j", "down":
			m.index = m.tokenField.switchFocusVertically(m.index, false)
		case "k", "up":
			m.index = m.tokenField.switchFocusVertically(m.index, true)
		}
	}

	return m, nil
}

func (tf *TokenField) switchFocusVertically(currentIndex int, up bool) int {
	currToken := 0
	for i, token := range tf.tokens {
		if token.index == currentIndex {
			currToken = i
			break
		}
	}

	focusedToken := tf.tokens[currToken]
	lastLine := tf.tokens[len(tf.tokens)-1].line
	if (up && focusedToken.line == 0) || (!up && focusedToken.line == lastLine) {
		return currentIndex
	}

	newLine := focusedToken.line
	if up {
		newLine--
	} else {
		newLine++
	}

	candidate := 0
	for i, token := range tf.tokens {
		if token.line == newLine {
			candidate = i
			break
		}
	}

	for {
		if candidate >= len(tf.tokens) {
			return len(tf.tokens) - 1
		}

		candidateToken := tf.tokens[candidate]
		if candidateToken.line == focusedToken.line {
			return candidate - 1
		}

		if candidateToken.end >= focusedToken.start {
			return candidate
		}

		candidate++
	}
}

func (tf *TokenField) renderTokens(focusedToken int) string {
	var netLineLength int = tf.width - 2*tf.horizontalPadding
	var sbTokenField strings.Builder

	line := 0
	index := 0
	renderedIndex := 0
	var sbLinePlain strings.Builder // Tracks plain text for layout decisions
	var sbLine strings.Builder      // Tracks actual rendered output
	for i, token := range tf.tokens {
		log.Println("========================================")
		log.Println("Word:", token.word)
		lineWithWord := sbLinePlain.String() + token.word
		if index > 0 {
			// Accounting for a space if not first word in line
			lineWithWord += " "
		}
		log.Printf("Index %d, lineww: %s\n", index, lineWithWord)

		if len(lineWithWord) > netLineLength {
			log.Printf("%d > %d, resetting line buffer\n", len(lineWithWord), netLineLength)
			sbTokenField.WriteString(sbLine.String())
			sbTokenField.WriteRune('\n')
			sbLine.Reset()
			sbLinePlain.Reset()
			line++
			index = 0
		}

		tf.tokens[i].start = index
		if index > 0 {
			tf.tokens[i].start++
			sbLine.WriteRune(' ')
			sbLinePlain.WriteRune(' ')
		}
		tf.tokens[i].end = tf.tokens[i].start + len(token.word)
		tf.tokens[i].line = line
		tf.tokens[i].index = renderedIndex
		log.Println(tf.tokens[i])

		renderedWord := token.word
		if focusedToken == i {
			renderedWord = defaultStyles().focusedToken.Render(renderedWord)
		}

		sbLine.WriteString(renderedWord)
		sbLinePlain.WriteString(token.word)
		index = tf.tokens[i].end
		renderedIndex++
		log.Println("========================================")
	}

	if sbLine.Len() > 0 {
		sbTokenField.WriteString(sbLine.String())
	}

	return sbTokenField.String()
}

func (tf *TokenField) View(width, height, focusedToken int) string {
	tf.width = width - 2
	tf.height = height - 2

	return lipgloss.NewStyle().
		Width(tf.width).
		Height(tf.height).
		Padding(tf.verticalPadding, tf.horizontalPadding).
		Border(lipgloss.NormalBorder()).
		Render(tf.renderTokens(focusedToken))
}

func (m *Model) View() string {
	return m.tokenField.View(m.width/2, m.height, m.index)
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
