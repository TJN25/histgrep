package utils

import (
	"fmt"
	"github.com/TJN25/histgrep/hsdata"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
	"os"
	"strconv"
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
	searchExcludes bool
	commandMode    bool
	commandInput   textinput.Model
}

func initialModel(content []string, data *hsdata.HsData, line hsdata.HsLine) Model {
	ti := textinput.New()
	ti.Placeholder = "Enter terms..."
	ti.CharLimit = 100

	ci := textinput.New()
	ci.Placeholder = ""
	ci.CharLimit = 100

	return Model{
		Content:      content,
		colorProfile: termenv.ColorProfile(),
		data:         data,
		line:         line,
		terms:        data.Terms,
		searchInput:  ti,
		commandInput: ci,
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
		if m.commandMode {
			switch msg.Type {
			case tea.KeyEnter:
				m.commandMode = false
				lineNum, err := strconv.Atoi(m.commandInput.Value())
				if err == nil && lineNum > 0 && lineNum <= len(m.Content) {
					m.cursor = lineNum - 1
					if m.cursor > len(m.Content)-m.viewportHeight {
						m.cursor = len(m.Content) - m.viewportHeight
					}
				}
				m.commandInput.SetValue("")
				return m, nil
			case tea.KeyEsc:
				m.commandMode = false
				m.commandInput.SetValue("")
				return m, nil
			}
			m.commandInput, cmd = m.commandInput.Update(msg)
			return m, cmd
		}
		if m.searchMode {
			switch msg.Type {
			case tea.KeyEnter:
				m.searchMode = false
				terms := strings.Fields(m.searchInput.Value())
				if m.searchExcludes {
					m.data.ExcludeTerms = terms
				} else {
					m.data.Terms = terms
				}
				if bufferedInput, ok := m.data.Reader.(*BufferedInput); ok {
					bufferedInput.Reset()
				}
				m.Content, _ = LoopFile(m.data, SaveLine, m.line)
				print(m.Content[0])
				m.cursor = 0
				return m, tea.Batch(tea.ClearScreen, tea.EnterAltScreen)
			case tea.KeySpace:
				terms := strings.Fields(m.searchInput.Value())
				if m.searchExcludes {
					m.data.ExcludeTerms = terms
				} else {
					m.data.Terms = terms
				}
				if bufferedInput, ok := m.data.Reader.(*BufferedInput); ok {
					bufferedInput.Reset()
				}
				m.Content, _ = LoopFile(m.data, SaveLine, m.line)
				print(m.Content[0])
				m.cursor = 0
				terms_string := strings.Join(terms, " ") + " "
				m.searchInput.SetValue(terms_string)
				m.searchInput.CursorEnd()
				return m, tea.Batch(tea.ClearScreen, tea.EnterAltScreen)
			case tea.KeyBackspace:
				terms_string := m.searchInput.Value()
				if len(terms_string) > 1 {
					if terms_string[len(terms_string)-1] == ' ' {
						terms := strings.Fields(m.searchInput.Value())
						if m.searchExcludes {
							m.data.ExcludeTerms = terms
						} else {
							m.data.Terms = terms
						}
						if bufferedInput, ok := m.data.Reader.(*BufferedInput); ok {
							bufferedInput.Reset()
						}
						m.Content, _ = LoopFile(m.data, SaveLine, m.line)
						print(m.Content[0])
						m.cursor = 0

					}
					terms_string = terms_string[:len(terms_string)-1]
					m.searchInput.SetValue(terms_string)
					m.searchInput.CursorEnd()
				} else {

					terms := strings.Fields(m.searchInput.Value())
					if m.searchExcludes {
						m.data.ExcludeTerms = terms
					} else {
						m.data.Terms = terms
					}
					if bufferedInput, ok := m.data.Reader.(*BufferedInput); ok {
						bufferedInput.Reset()
					}
					m.Content, _ = LoopFile(m.data, SaveLine, m.line)
					print(m.Content[0])
					m.cursor = 0
					terms_string = ""
					m.searchInput.SetValue(terms_string)
					m.searchInput.CursorEnd()
				}
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
		case " ":
			m.cursor += m.viewportHeight
			if m.cursor > len(m.Content)-m.viewportHeight {
				m.cursor = len(m.Content) - m.viewportHeight
			}
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
			m.searchExcludes = false
			m.searchInput.SetValue(strings.Join(m.data.Terms, " "))
			m.searchInput.Focus()
			m.searchInput.CursorEnd()
			return m, textinput.Blink
		case "?":
			m.searchMode = true
			m.searchExcludes = true
			m.searchInput.SetValue(strings.Join(m.data.ExcludeTerms, " "))
			m.searchInput.Focus()
			m.searchInput.CursorEnd()
			return m, textinput.Blink
		case ":":
			m.commandMode = true
			m.commandInput.Focus()
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
	if m.commandMode {
		statusLine = boldStyle.Styled(":") + m.commandInput.View()
	} else if m.searchMode {
		if m.searchExcludes {
			statusLine = boldStyle.Styled("Exclude: ") + m.searchInput.View()
		} else {
			statusLine = boldStyle.Styled("Search: ") + m.searchInput.View()
		}
	} else {

		terms := boldStyle.Styled(strings.Join(m.terms, ", "))
		statusInfo := regularStyle.Styled(fmt.Sprintf(" line %d of %d | Searching for terms: %s | Excluding terms: %s | (use '/' to search and '?' to exclude or press q to quit)", m.cursor+1, len(m.Content), strings.Join(m.terms, ", "), strings.Join(m.data.ExcludeTerms, ", ")))

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
