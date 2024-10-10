package utils

import (
	"fmt"
	"github.com/TJN25/histgrep/hsdata"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
	"os"
	"strings"
)

// Main struct for the displayed data
type Model struct {
	Content        []string
	cursor         int
	viewportHeight int
	colorProfile   termenv.Profile
	terms          []string
	data           *hsdata.HsData
	line           hsdata.HsLine
	searchInput    textinput.Model
	searchMode     bool
}

func initialModel(content []string, data *hsdata.HsData, line hsdata.HsLine) Model {
	ti := textinput.New()
	ti.Placeholder = "Enter search terms..."
	ti.CharLimit = 100

	return Model{
		Content:      content,
		colorProfile: termenv.ColorProfile(),
		data:         data,
		line:         line,
		terms:        data.Terms,
		searchInput:  ti,
	}
}

// Set up the initial state. Currently unused.
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Process keyboard events, along with anything else and update the state
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.searchMode {
			switch msg.Type {
			case tea.KeyEnter:
				m.searchMode = false
				m.data.Terms = strings.Fields(m.searchInput.Value())
				m.Content, _ = LoopFile(m.data, SaveLine, m.line)
				print(m.Content[0])
				m.cursor = 0
				m.terms = m.data.Terms
				return m, tea.Batch(tea.ClearScreen, tea.EnterAltScreen)
			case tea.KeyEsc:
				m.searchMode = false
				m.searchInput.SetValue("")
				return m, nil
			}
			m.searchInput, cmd = m.searchInput.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.Content)-m.viewportHeight {
				m.cursor++
			}
		case "g":
			m.cursor = 0
		case "G":
			m.cursor = len(m.Content) - m.viewportHeight
			if m.cursor < 0 {
				m.cursor = 0
			}
		case "ctrl+u":
			m.cursor -= m.viewportHeight / 2
			if m.cursor < 0 {
				m.cursor = 0
			}
		case "ctrl+d":
			m.cursor += m.viewportHeight / 2
			if m.cursor > len(m.Content)-m.viewportHeight {
				m.cursor = len(m.Content) - m.viewportHeight
			}
			if m.cursor < 0 {
				m.cursor = 0
			}
		case "/":
			m.searchMode = true
			m.searchInput.Focus()
			return m, textinput.Blink
		}
	case tea.WindowSizeMsg:
		m.viewportHeight = msg.Height - 1 // Leave one line for status
	}

	return m, nil
}

// How the screen is rendered
func (m Model) View() string {
	if m.viewportHeight == 0 {
		return "Loading..."
	}

	start := m.cursor
	end := m.cursor + m.viewportHeight - 1 // -1 to leave room for status line
	if end > len(m.Content) {
		end = len(m.Content)
	}

	content := strings.Join(m.Content[start:end], "\n")

	// Create status line
	statusBg := m.colorProfile.Color("236") // Dark grey that works for both themes
	statusFg := m.colorProfile.Color("252") // Light grey that works for both themes
	boldFg := m.colorProfile.Color("255")   // Almost white, for bold text

	statusStyle := termenv.Style{}.Background(statusBg).Foreground(statusFg)
	boldStyle := termenv.Style{}.Background(statusBg).Foreground(boldFg).Bold()
	regularStyle := termenv.Style{}.Background(statusBg).Foreground(boldFg)
	var statusLine string
	if m.searchMode {
		statusLine = boldStyle.Styled("Search: ") + m.searchInput.View()
	} else {

		terms := boldStyle.Styled(strings.Join(m.terms, ", "))
		statusInfo := regularStyle.Styled(fmt.Sprintf(" line %d of %d (press q or C-c to quit)", m.cursor+1, len(m.Content)))

		statusLine = statusStyle.Styled(fmt.Sprintf("%s%s", terms, statusInfo))
	}

	// Ensure status line fills the entire width
	paddedStatusLine := fmt.Sprintf("%-*s", m.viewportHeight, statusLine)

	// Truncate if status line is too long
	if len(paddedStatusLine) > m.viewportHeight {
		paddedStatusLine = paddedStatusLine[:m.viewportHeight-3] + "..."
	}
	if len(m.Content) < m.viewportHeight {
		content += strings.Repeat("\n", m.viewportHeight-len(m.Content)-1)
	}

	return content + "\n" + statusLine
}

func ViewFileWithPager(content []string, data *hsdata.HsData, line hsdata.HsLine) {
	model := initialModel(content, data, line)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
