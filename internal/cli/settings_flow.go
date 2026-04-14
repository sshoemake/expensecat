package cli

import (
	"expensecat/internal/storage"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var settingsLogo = ` /\_/\  
( -.- )  
(  >⚙️ )  Settings`

const (
	SettingsExclusionList = iota
	SettingsImportPath
	SettingsBack
)

type SettingsModel struct {
	store    storage.Storage
	cursor   int
	quitting bool
}

func NewSettingsModel(store storage.Storage) SettingsModel {
	return SettingsModel{
		store:  store,
		cursor: 0,
	}
}

func (m SettingsModel) Init() tea.Cmd {
	return nil
}

func (m SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "esc":
			return m, func() tea.Msg {
				return NavigationMsg{Destination: ScreenMainMenu}
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < 2 {
				m.cursor++
			}
		case "enter":
			if m.cursor == SettingsExclusionList {
				return m, func() tea.Msg {
					return NavigationMsg{Destination: ScreenExclusionList}
				}
			}
			if m.cursor == SettingsImportPath {
				return m, func() tea.Msg {
					return NavigationMsg{Destination: ScreenImportPath}
				}
			}
			if m.cursor == SettingsBack {
				return m, func() tea.Msg {
					return NavigationMsg{Destination: ScreenMainMenu}
				}
			}
		}
	}
	return m, nil
}

func (m SettingsModel) View() string {
	var s strings.Builder

	s.WriteString(settingsLogo + "\n\n")
	s.WriteString("Settings\n")
	s.WriteString("─────────────────────────────\n")

	choices := []string{"Exclusion List", "Import Path", "Back"}
	for i, choice := range choices {
		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}
		s.WriteString(fmt.Sprintf("%s%s\n", cursor, choice))
	}

	s.WriteString("\n↑↓ Navigate  Enter Select\n")
	s.WriteString("esc Back to Main Menu\n")

	return s.String()
}

func (m SettingsModel) IsQuitting() bool {
	return m.quitting
}

const (
	ExclusionListMenu = iota
	ExclusionListAdd
	ExclusionListRemove
	ExclusionListBack
)

type ExclusionListModel struct {
	store      storage.Storage
	patterns   []string
	cursor     int
	mode       int
	inputValue string
	errorMsg   string
	successMsg string
	quitting   bool
}

func NewExclusionListModel(store storage.Storage) ExclusionListModel {
	m := ExclusionListModel{
		store: store,
		mode:  ExclusionListMenu,
	}
	m.loadPatterns()
	return m
}

func (m *ExclusionListModel) loadPatterns() {
	patterns, err := m.store.GetExclusionList()
	if err != nil {
		m.errorMsg = fmt.Sprintf("Failed to load patterns: %v", err)
		return
	}
	m.patterns = patterns
}

func (m ExclusionListModel) Init() tea.Cmd {
	return nil
}

func (m ExclusionListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "esc":
			if m.mode == ExclusionListAdd {
				m.inputValue = ""
				m.mode = ExclusionListMenu
				m.cursor = 0
			} else if m.mode == ExclusionListRemove {
				m.mode = ExclusionListMenu
				m.cursor = 0
			} else if m.mode == ExclusionListMenu {
				return m, func() tea.Msg {
					return NavigationMsg{Destination: ScreenSettings}
				}
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.mode == ExclusionListMenu {
				if m.cursor < 2 {
					m.cursor++
				}
			} else if m.mode == ExclusionListRemove {
				maxChoices := len(m.patterns) + 1
				if m.cursor < maxChoices-1 {
					m.cursor++
				}
			}
		case "enter":
			if m.mode == ExclusionListMenu {
				if m.cursor == 0 {
					m.mode = ExclusionListAdd
					m.inputValue = ""
				} else if m.cursor == 1 {
					m.mode = ExclusionListRemove
					m.cursor = 0
				} else if m.cursor == 2 {
					return m, func() tea.Msg {
						return NavigationMsg{Destination: ScreenSettings}
					}
				}
			} else if m.mode == ExclusionListAdd {
				if m.inputValue != "" {
					err := m.store.AddExclusion(m.inputValue)
					if err != nil {
						m.errorMsg = err.Error()
					} else {
						m.successMsg = "Pattern added successfully"
						m.loadPatterns()
					}
					m.inputValue = ""
					m.mode = ExclusionListMenu
					m.cursor = 0
				}
			} else if m.mode == ExclusionListRemove {
				if m.cursor < len(m.patterns) {
					m.removePattern(m.patterns[m.cursor])
					m.loadPatterns()
					m.cursor = 0
					m.mode = ExclusionListMenu
				} else {
					m.mode = ExclusionListMenu
					m.cursor = 0
				}
			}
		case "backspace":
			if m.mode == ExclusionListAdd && m.inputValue != "" {
				m.inputValue = m.inputValue[:len(m.inputValue)-1]
			}
		default:
			if m.mode == ExclusionListAdd && msg.String() != "" {
				m.inputValue += msg.String()
			}
		}
	}
	return m, nil
}

func (m *ExclusionListModel) removePattern(pattern string) {
	err := m.store.RemoveExclusion(pattern)
	if err != nil {
		m.errorMsg = err.Error()
	} else {
		m.successMsg = "Pattern removed successfully"
	}
}

func (m ExclusionListModel) View() string {
	var s strings.Builder

	s.WriteString(settingsLogo + "\n\n")
	s.WriteString("Exclusion List\n")
	s.WriteString("─────────────────────────────\n")

	if m.errorMsg != "" {
		s.WriteString(fmt.Sprintf("Error: %s\n\n", m.errorMsg))
		m.errorMsg = ""
	}
	if m.successMsg != "" {
		s.WriteString(fmt.Sprintf("✓ %s\n\n", m.successMsg))
		m.successMsg = ""
	}

	if m.mode == ExclusionListAdd {
		s.WriteString("Enter new pattern:\n")
		s.WriteString(fmt.Sprintf("> %s_\n", m.inputValue))
		s.WriteString("\nPress Enter to add, Esc to cancel\n")
	} else if m.mode == ExclusionListRemove {
		s.WriteString("Select pattern to remove:\n\n")
		for i, pattern := range m.patterns {
			cursor := "  "
			if m.cursor == i {
				cursor = "> "
			}
			s.WriteString(fmt.Sprintf("%s%s\n", cursor, pattern))
		}
		cursor := "  "
		if m.cursor == len(m.patterns) {
			cursor = "> "
		}
		s.WriteString(fmt.Sprintf("%s[Back to Settings]\n", cursor))
		s.WriteString("\n↑↓ Navigate  Enter to remove\n")
		s.WriteString("Esc to cancel\n")
	} else {
		s.WriteString(fmt.Sprintf("Current patterns (%d):\n\n", len(m.patterns)))
		if len(m.patterns) == 0 {
			s.WriteString("  (no patterns defined)\n")
		} else {
			for _, pattern := range m.patterns {
				s.WriteString(fmt.Sprintf("  • %s\n", pattern))
			}
		}

		s.WriteString("\n")
		choices := []string{"Add Pattern", "Remove Pattern", "Back"}
		for i, choice := range choices {
			cursor := "  "
			if m.cursor == i {
				cursor = "> "
			}
			s.WriteString(fmt.Sprintf("%s%s\n", cursor, choice))
		}

		s.WriteString("\n↑↓ Navigate  Enter Select\n")
		s.WriteString("esc Back to Settings\n")
	}

	return s.String()
}

func (m ExclusionListModel) IsQuitting() bool {
	return m.quitting
}

const (
	ImportPathMenu = iota
	ImportPathEdit
	ImportPathBack
)

type ImportPathModel struct {
	store      storage.Storage
	path       string
	cursor     int
	mode       int
	inputValue string
	cursorPos  int
	errorMsg   string
	successMsg string
	quitting   bool
}

func NewImportPathModel(store storage.Storage) ImportPathModel {
	m := ImportPathModel{
		store: store,
		mode:  ImportPathMenu,
	}
	m.loadPath()
	return m
}

func (m *ImportPathModel) loadPath() {
	path, err := m.store.GetImportPath()
	if err != nil {
		m.errorMsg = fmt.Sprintf("Failed to load path: %v", err)
		return
	}
	m.path = path
}

func (m ImportPathModel) Init() tea.Cmd {
	return nil
}

func (m ImportPathModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "esc":
			if m.mode == ImportPathEdit {
				m.inputValue = ""
				m.cursorPos = 0
				m.mode = ImportPathMenu
				m.cursor = 0
			} else if m.mode == ImportPathMenu {
				return m, func() tea.Msg {
					return NavigationMsg{Destination: ScreenSettings}
				}
			}
		case "left":
			if m.mode == ImportPathEdit && m.cursorPos > 0 {
				m.cursorPos--
			}
		case "right":
			if m.mode == ImportPathEdit && m.cursorPos < len(m.inputValue) {
				m.cursorPos++
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.mode == ImportPathMenu {
				if m.cursor < 1 {
					m.cursor++
				}
			}
		case "enter":
			if m.mode == ImportPathMenu {
				if m.cursor == 0 {
					m.mode = ImportPathEdit
					m.inputValue = m.path
					m.cursorPos = len(m.path)
				} else if m.cursor == 1 {
					return m, func() tea.Msg {
						return NavigationMsg{Destination: ScreenSettings}
					}
				}
			} else if m.mode == ImportPathEdit {
				if m.inputValue != "" {
					err := m.store.UpdateImportPath(m.inputValue)
					if err != nil {
						m.errorMsg = err.Error()
					} else {
						m.successMsg = "Path updated successfully"
						m.loadPath()
					}
					m.inputValue = ""
					m.cursorPos = 0
					m.mode = ImportPathMenu
					m.cursor = 0
				}
			}
		case "backspace":
			if m.mode == ImportPathEdit && m.inputValue != "" {
				if m.cursorPos > 0 {
					m.inputValue = m.inputValue[:m.cursorPos-1] + m.inputValue[m.cursorPos:]
					m.cursorPos--
				}
			}
		default:
			if m.mode == ImportPathEdit {
				key := msg.String()
				if key == "left" || key == "right" || key == "up" || key == "down" || key == "ctrl+c" || key == "ctrl+u" || key == "ctrl+w" {
					return m, nil
				}
				if len(key) == 1 {
					m.inputValue = m.inputValue[:m.cursorPos] + key + m.inputValue[m.cursorPos:]
					m.cursorPos++
				}
			}
		}
	}
	return m, nil
}

func (m ImportPathModel) View() string {
	var s strings.Builder

	s.WriteString(settingsLogo + "\n\n")
	s.WriteString("Import Path\n")
	s.WriteString("─────────────────────────────\n")

	if m.errorMsg != "" {
		s.WriteString(fmt.Sprintf("Error: %s\n\n", m.errorMsg))
		m.errorMsg = ""
	}
	if m.successMsg != "" {
		s.WriteString(fmt.Sprintf("✓ %s\n\n", m.successMsg))
		m.successMsg = ""
	}

	if m.mode == ImportPathEdit {
		s.WriteString("Enter new path:\n")
		before := m.inputValue[:m.cursorPos]
		after := m.inputValue[m.cursorPos:]
		s.WriteString(fmt.Sprintf("> %s|%s\n", before, after))
		s.WriteString("\n←→ Move cursor  Enter to save, Esc to cancel\n")
	} else {
		s.WriteString(fmt.Sprintf("Current path:\n\n  %s\n\n", m.path))

		choices := []string{"Edit Path", "Back"}
		for i, choice := range choices {
			cursor := "  "
			if m.cursor == i {
				cursor = "> "
			}
			s.WriteString(fmt.Sprintf("%s%s\n", cursor, choice))
		}

		s.WriteString("\n↑↓ Navigate  Enter Select\n")
		s.WriteString("esc Back to Settings\n")
	}

	return s.String()
}

func (m ImportPathModel) IsQuitting() bool {
	return m.quitting
}
