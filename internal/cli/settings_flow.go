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
	SettingsCategoryList
	SettingsRecurringList
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
			if m.cursor < 4 {
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
			if m.cursor == SettingsCategoryList {
				return m, func() tea.Msg {
					return NavigationMsg{Destination: ScreenCategoryList}
				}
			}
			if m.cursor == SettingsRecurringList {
				return m, func() tea.Msg {
					return NavigationMsg{Destination: ScreenRecurringList}
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

	choices := []string{"Exclusion List", "Import Path", "Category List", "Recurring Expenses", "Back"}
	for i, choice := range choices {
		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}
		fmt.Fprintf(&s, "%s%s\n", cursor, choice)
	}

	s.WriteString("\n↑↓ Navigate  Enter Select\n")
	s.WriteString("esc Back to Main Menu\n")

	return s.String()
}

func (m SettingsModel) IsQuitting() bool {
	return m.quitting
}

const (
	ExclusionListView = iota
	ExclusionListAdd
	ExclusionListEdit
	ExclusionListDeleteConfirm
	ExclusionListBack
)

type ExclusionListModel struct {
	store      storage.Storage
	patterns   []string
	cursor     int
	mode       int
	editIndex  int
	inputValue string
	errorMsg   string
	successMsg string
	quitting   bool
}

func NewExclusionListModel(store storage.Storage) ExclusionListModel {
	m := ExclusionListModel{
		store:     store,
		mode:      ExclusionListView,
		editIndex: -1,
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
			switch m.mode {
			case ExclusionListAdd, ExclusionListEdit:
				m.inputValue = ""
				m.mode = ExclusionListView
				m.editIndex = -1
			case ExclusionListDeleteConfirm:
				m.mode = ExclusionListView
			case ExclusionListView:
				return m, func() tea.Msg {
					return NavigationMsg{Destination: ScreenSettings}
				}
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			switch m.mode {
			case ExclusionListView:
				if m.cursor < len(m.patterns) {
					m.cursor++
				}
			}
		case "enter":
			switch m.mode {
			case ExclusionListView:
				if m.cursor < len(m.patterns) {
					m.mode = ExclusionListEdit
					m.editIndex = m.cursor
					m.inputValue = m.patterns[m.cursor]
				} else if m.cursor == len(m.patterns) {
					return m, func() tea.Msg {
						return NavigationMsg{Destination: ScreenSettings}
					}
				}
			case ExclusionListAdd:
				if m.inputValue != "" {
					err := m.store.AddExclusion(m.inputValue)
					if err != nil {
						m.errorMsg = err.Error()
					} else {
						m.successMsg = "Pattern added successfully"
						m.loadPatterns()
					}
					m.inputValue = ""
					m.mode = ExclusionListView
				}
			case ExclusionListEdit:
				if m.inputValue != "" && m.editIndex >= 0 {
					oldPattern := m.patterns[m.editIndex]
					err := m.store.RemoveExclusion(oldPattern)
					if err != nil {
						m.errorMsg = err.Error()
					} else {
						err = m.store.AddExclusion(m.inputValue)
						if err != nil {
							m.errorMsg = err.Error()
							m.store.AddExclusion(oldPattern)
						} else {
							m.successMsg = "Pattern updated successfully"
							m.loadPatterns()
						}
					}
					m.inputValue = ""
					m.editIndex = -1
					m.mode = ExclusionListView
				}
			case ExclusionListDeleteConfirm:
				if m.cursor < len(m.patterns) {
					m.removePattern(m.patterns[m.cursor])
					m.loadPatterns()
					m.cursor = 0
					m.mode = ExclusionListView
				} else {
					m.mode = ExclusionListView
				}
			}
		case "y":
			if m.mode == ExclusionListDeleteConfirm && m.cursor < len(m.patterns) {
				m.removePattern(m.patterns[m.cursor])
				m.loadPatterns()
				m.cursor = 0
				m.mode = ExclusionListView
			}
		case "n":
			if m.mode == ExclusionListDeleteConfirm {
				m.mode = ExclusionListView
			}
		case "backspace":
			if (m.mode == ExclusionListAdd || m.mode == ExclusionListEdit) && m.inputValue != "" {
				m.inputValue = m.inputValue[:len(m.inputValue)-1]
			}
		case "left", "right":
			// Ignore arrow keys in exclusion list input modes
		case " ":
			if (m.mode == ExclusionListAdd || m.mode == ExclusionListEdit) && msg.Type == tea.KeySpace {
				m.inputValue += " "
			}
		default:
			if (m.mode == ExclusionListAdd || m.mode == ExclusionListEdit) && msg.Type == tea.KeyRunes {
				m.inputValue += msg.String()
			}
			if m.mode == ExclusionListView {
				switch msg.String() {
				case "a":
					m.mode = ExclusionListAdd
					m.inputValue = ""
				case "e":
					if m.cursor < len(m.patterns) {
						m.mode = ExclusionListEdit
						m.editIndex = m.cursor
						m.inputValue = m.patterns[m.cursor]
					}
				case "x":
					if m.cursor < len(m.patterns) {
						m.mode = ExclusionListDeleteConfirm
					}
				}
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
		fmt.Fprintf(&s, "Error: %s\n\n", m.errorMsg)
	}
	if m.successMsg != "" {
		fmt.Fprintf(&s, "✓ %s\n\n", m.successMsg)
	}

	switch m.mode {
	case ExclusionListAdd:
		s.WriteString("Enter new pattern:\n")
		fmt.Fprintf(&s, "> %s_\n", m.inputValue)
		s.WriteString("\nPress Enter to add, Esc to cancel\n")
	case ExclusionListEdit:
		s.WriteString("Edit pattern:\n")
		fmt.Fprintf(&s, "> %s_\n", m.inputValue)
		s.WriteString("\nPress Enter to save, Esc to cancel\n")
	case ExclusionListDeleteConfirm:
		s.WriteString("Confirm delete")
		if m.cursor < len(m.patterns) {
			fmt.Fprintf(&s, " '%s'", m.patterns[m.cursor])
		}
		s.WriteString("?\n\n[Y]es / [N]o\n")
	default:
		fmt.Fprintf(&s, "Patterns (%d):\n\n", len(m.patterns))
		if len(m.patterns) == 0 {
			s.WriteString("  (no patterns defined)\n")
		} else {
			for i, pattern := range m.patterns {
				cursor := "  "
				if m.cursor == i {
					cursor = "> "
				}
				fmt.Fprintf(&s, "%s%s\n", cursor, pattern)
			}
		}

		s.WriteString("\na Add     e Edit     x Delete     esc Back to Settings\n")
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
			switch m.mode {
			case ImportPathEdit:
				m.inputValue = ""
				m.cursorPos = 0
				m.mode = ImportPathMenu
				m.cursor = 0
			case ImportPathMenu:
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
			switch m.mode {
			case ImportPathMenu:
				switch m.cursor {
				case 0:
					m.mode = ImportPathEdit
					m.inputValue = m.path
					m.cursorPos = len(m.path)
				case 1:
					return m, func() tea.Msg {
						return NavigationMsg{Destination: ScreenSettings}
					}
				}
			case ImportPathEdit:
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
		fmt.Fprintf(&s, "Error: %s\n\n", m.errorMsg)
	}
	if m.successMsg != "" {
		fmt.Fprintf(&s, "✓ %s\n\n", m.successMsg)
	}

	if m.mode == ImportPathEdit {
		s.WriteString("Enter new path:\n")
		before := m.inputValue[:m.cursorPos]
		after := m.inputValue[m.cursorPos:]
		fmt.Fprintf(&s, "> %s|%s\n", before, after)
		s.WriteString("\n←→ Move cursor  Enter to save, Esc to cancel\n")
	} else {
		fmt.Fprintf(&s, "Current path:\n\n  %s\n\n", m.path)

		choices := []string{"Edit Path", "Back"}
		for i, choice := range choices {
			cursor := "  "
			if m.cursor == i {
				cursor = "> "
			}
			fmt.Fprintf(&s, "%s%s\n", cursor, choice)
		}

		s.WriteString("\n↑↓ Navigate  Enter Select\n")
		s.WriteString("esc Back to Settings\n")
	}

	return s.String()
}

func (m ImportPathModel) IsQuitting() bool {
	return m.quitting
}
