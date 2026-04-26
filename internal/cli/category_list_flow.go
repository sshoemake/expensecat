package cli

import (
	"expensecat/internal/storage"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	CategoryListView = iota
	CategoryListAdd
	CategoryListEdit
	CategoryListDeleteConfirm
	CategoryListBack
)

type CategoryListModel struct {
	store      storage.Storage
	categories []string
	cursor     int
	mode       int
	editIndex  int
	inputValue string
	errorMsg   string
	successMsg string
	quitting   bool
}

func NewCategoryListModel(store storage.Storage) CategoryListModel {
	m := CategoryListModel{
		store:     store,
		mode:      CategoryListView,
		editIndex: -1,
	}
	m.loadCategories()
	return m
}

func (m *CategoryListModel) loadCategories() {
	categories, err := m.store.GetCategories()
	if err != nil {
		m.errorMsg = fmt.Sprintf("Failed to load categories: %v", err)
		return
	}
	m.categories = categories
}

func (m CategoryListModel) Init() tea.Cmd {
	return nil
}

func (m CategoryListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "esc":
			switch m.mode {
			case CategoryListAdd, CategoryListEdit:
				m.inputValue = ""
				m.mode = CategoryListView
				m.editIndex = -1
			case CategoryListDeleteConfirm:
				m.mode = CategoryListView
			case CategoryListView:
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
			case CategoryListView:
				if m.cursor < len(m.categories) {
					m.cursor++
				}
			}
		case "enter":
			switch m.mode {
			case CategoryListView:
				if m.cursor < len(m.categories) {
					m.mode = CategoryListEdit
					m.editIndex = m.cursor
					m.inputValue = m.categories[m.cursor]
				} else if m.cursor == len(m.categories) {
					return m, func() tea.Msg {
						return NavigationMsg{Destination: ScreenSettings}
					}
				}
			case CategoryListAdd:
				if m.inputValue != "" {
					err := m.addCategory(m.inputValue)
					if err != nil {
						m.errorMsg = err.Error()
					} else {
						m.successMsg = "Category added successfully"
						m.loadCategories()
					}
					m.inputValue = ""
					m.mode = CategoryListView
				}
			case CategoryListEdit:
				if m.inputValue != "" && m.editIndex >= 0 {
					newCategories := make([]string, len(m.categories))
					copy(newCategories, m.categories)
					newCategories[m.editIndex] = m.inputValue
					err := m.store.UpdateCategories(newCategories)
					if err != nil {
						m.errorMsg = err.Error()
					} else {
						m.successMsg = "Category updated successfully"
						m.loadCategories()
					}
					m.inputValue = ""
					m.editIndex = -1
					m.mode = CategoryListView
				}
			case CategoryListDeleteConfirm:
				if m.cursor < len(m.categories) {
					m.removeCategory(m.categories[m.cursor])
					m.loadCategories()
					m.cursor = 0
					m.mode = CategoryListView
				} else {
					m.mode = CategoryListView
				}
			}
		case "y":
			if m.mode == CategoryListDeleteConfirm && m.cursor < len(m.categories) {
				m.removeCategory(m.categories[m.cursor])
				m.loadCategories()
				m.cursor = 0
				m.mode = CategoryListView
			}
		case "n":
			if m.mode == CategoryListDeleteConfirm {
				m.mode = CategoryListView
			}
		case "backspace":
			if (m.mode == CategoryListAdd || m.mode == CategoryListEdit) && m.inputValue != "" {
				m.inputValue = m.inputValue[:len(m.inputValue)-1]
			}
		case "left", "right":
			// Ignore arrow keys in category list input modes
		case " ":
			if (m.mode == CategoryListAdd || m.mode == CategoryListEdit) && msg.Type == tea.KeySpace {
				m.inputValue += " "
			}
		default:
			if (m.mode == CategoryListAdd || m.mode == CategoryListEdit) && msg.Type == tea.KeyRunes {
				m.inputValue += msg.String()
			}
			if m.mode == CategoryListView {
				switch msg.String() {
				case "a":
					m.mode = CategoryListAdd
					m.inputValue = ""
				case "e":
					if m.cursor < len(m.categories) {
						m.mode = CategoryListEdit
						m.editIndex = m.cursor
						m.inputValue = m.categories[m.cursor]
					}
				case "x":
					if m.cursor < len(m.categories) {
						m.mode = CategoryListDeleteConfirm
					}
				}
			}
		}
	}
	return m, nil
}

func (m *CategoryListModel) addCategory(category string) error {
	newCategories := make([]string, len(m.categories)+1)
	copy(newCategories, m.categories)
	newCategories[len(m.categories)] = category
	return m.store.UpdateCategories(newCategories)
}

func (m *CategoryListModel) removeCategory(category string) {
	newCategories := []string{}
	for _, c := range m.categories {
		if c != category {
			newCategories = append(newCategories, c)
		}
	}
	err := m.store.UpdateCategories(newCategories)
	if err != nil {
		m.errorMsg = fmt.Sprintf("Failed to remove category: %v", err)
	} else {
		m.successMsg = "Category removed successfully"
	}
}

func (m CategoryListModel) View() string {
	var s strings.Builder

	s.WriteString(settingsLogo + "\n\n")
	s.WriteString("Category List\n")
	s.WriteString("─────────────────────────────\n")

	if m.errorMsg != "" {
		fmt.Fprintf(&s, "Error: %s\n\n", m.errorMsg)
	}
	if m.successMsg != "" {
		fmt.Fprintf(&s, "✓ %s\n\n", m.successMsg)
	}

	switch m.mode {
	case CategoryListAdd:
		s.WriteString("Enter new category:\n")
		fmt.Fprintf(&s, "> %s_\n", m.inputValue)
		s.WriteString("\nPress Enter to add, Esc to cancel\n")
	case CategoryListEdit:
		s.WriteString("Edit category:\n")
		fmt.Fprintf(&s, "> %s_\n", m.inputValue)
		s.WriteString("\nPress Enter to save, Esc to cancel\n")
	case CategoryListDeleteConfirm:
		s.WriteString("Confirm delete")
		if m.cursor < len(m.categories) {
			fmt.Fprintf(&s, " '%s'", m.categories[m.cursor])
		}
		s.WriteString("?\n\n[Y]es / [N]o\n")
	default:
		fmt.Fprintf(&s, "Categories (%d):\n\n", len(m.categories))
		if len(m.categories) == 0 {
			s.WriteString("  (no categories defined)\n")
		} else {
			for i, cat := range m.categories {
				cursor := "  "
				if m.cursor == i {
					cursor = "> "
				}
				fmt.Fprintf(&s, "%s%s\n", cursor, cat)
			}
		}

		s.WriteString("\na Add     e Edit     x Delete     esc Back to Settings\n")
	}

	return s.String()
}

func (m CategoryListModel) IsQuitting() bool {
	return m.quitting
}
