package cli

import (
	"expensecat/internal/api"
	"expensecat/internal/storage"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var importLogo = ` /\_/\  
( $.$ )  
(  >💳 )  ExpenseCat`

type ImportModel struct {
	storage         storage.Storage
	basePath        string
	selectedFile    string
	files           []string
	cursor          int
	fileListVisible bool
	fileCursor      int
	showResults     bool
	result          *api.ImportResult
	errorMsg        string
	dirError        string
	quitting        bool
}

func NewImportModel(store storage.Storage, basePath string) ImportModel {
	m := ImportModel{
		storage:  store,
		basePath: basePath,
	}
	m.loadFiles()
	return m
}

func (m *ImportModel) loadFiles() {
	expandedPath := expandPath(m.basePath)

	entries, err := os.ReadDir(expandedPath)
	if err != nil {
		m.dirError = fmt.Sprintf("Cannot read directory: %v", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			m.files = append(m.files, entry.Name())
		}
	}
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

func (m ImportModel) Init() tea.Cmd {
	return nil
}

func (m ImportModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.fileListVisible {
				if m.fileCursor > 0 {
					m.fileCursor--
				}
			} else {
				if m.cursor > 0 {
					m.cursor--
				}
			}
		case "down", "j":
			if m.fileListVisible {
				maxCursor := len(m.files)
				if m.fileCursor < maxCursor-1 {
					m.fileCursor++
				}
			} else {
				if m.cursor < 2 {
					m.cursor++
				}
			}
		case "enter":
			if m.showResults {
				return m, func() tea.Msg {
					return NavigationMsg{Destination: ScreenMainMenu}
				}
			}

			if m.fileListVisible {
				if m.fileCursor < len(m.files) {
					m.selectedFile = m.files[m.fileCursor]
					m.fileListVisible = false
				} else {
					m.fileListVisible = false
				}
				return m, nil
			}

			if m.cursor == 0 {
				m.fileListVisible = true
				m.fileCursor = 0
			} else if m.cursor == 1 {
				m.executeImport()
			} else if m.cursor == 2 {
				return m, func() tea.Msg {
					return NavigationMsg{Destination: ScreenMainMenu}
				}
			}
		case "backspace":
			if m.fileListVisible && len(m.selectedFile) > 0 {
				m.selectedFile = m.selectedFile[:len(m.selectedFile)-1]
			}
		}
	}
	return m, nil
}

func (m *ImportModel) executeImport() {
	m.errorMsg = ""

	if m.selectedFile == "" {
		m.errorMsg = "No file selected"
		return
	}

	fullPath := filepath.Join(expandPath(m.basePath), m.selectedFile)

	result, err := api.ImportCSVFile(m.storage, fullPath)
	if err != nil {
		m.errorMsg = fmt.Sprintf("Import failed: %v", err)
		return
	}

	m.result = result
	m.showResults = true
	m.fileListVisible = false
}

func (m ImportModel) View() string {
	var s strings.Builder

	s.WriteString(importLogo + "\n\n")

	if m.showResults {
		s.WriteString("Import Results\n")
		s.WriteString("─────────────────────────────\n")
		if m.result != nil {
			s.WriteString(fmt.Sprintf("Imported:    %d expenses\n", m.result.Imported))
			s.WriteString(fmt.Sprintf("Skipped:     %d records\n", m.result.Skipped))
			s.WriteString(fmt.Sprintf("Duplicates:  %d records\n", m.result.Duplicates))
		}
		s.WriteString("\nPress any key to return\n")
		return s.String()
	}

	s.WriteString("Import Expenses\n")
	s.WriteString("─────────────────────────────\n")
	s.WriteString(fmt.Sprintf("Base path: %s\n", m.basePath))

	if m.dirError != "" {
		s.WriteString(fmt.Sprintf("\nError: %s\n", m.dirError))
	}

	s.WriteString("\n")

	if m.fileListVisible {
		if len(m.files) == 0 {
			s.WriteString("No files found in directory\n\n")
		} else {
			s.WriteString("Select a file to import:\n")
			for i, file := range m.files {
				cursor := "  "
				if m.fileCursor == i {
					cursor = "> "
				}
				s.WriteString(fmt.Sprintf("%s%s\n", cursor, file))
			}
		}

		s.WriteString("\n")
		cursor := "  "
		if m.fileCursor == len(m.files) {
			cursor = "> "
		}
		s.WriteString(fmt.Sprintf("%s[Back to menu]\n", cursor))

		s.WriteString("\nUp/Down to navigate, Enter to select, Esc to go back\n")
	} else {
		s.WriteString(fmt.Sprintf("Selected file: %s\n", m.selectedFile))
		if m.errorMsg != "" {
			s.WriteString(fmt.Sprintf("\nError: %s\n", m.errorMsg))
		}

		s.WriteString("\n")
		cursor := "  "
		if m.cursor == 0 {
			cursor = "> "
		}
		s.WriteString(fmt.Sprintf("%s[Select File]\n", cursor))

		cursor = "  "
		if m.cursor == 1 {
			cursor = "> "
		}
		s.WriteString(fmt.Sprintf("%s[Import]\n", cursor))

		cursor = "  "
		if m.cursor == 2 {
			cursor = "> "
		}
		s.WriteString(fmt.Sprintf("%s[Back]\n", cursor))

		s.WriteString("\nUp/Down to navigate, Enter to select option, Esc to go back\n")
	}

	return s.String()
}

func (m ImportModel) IsQuitting() bool {
	return m.quitting
}
