package cli

import (
	"expensecat/internal/storage"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	RecurringView = iota
	RecurringAdd
	RecurringEdit
	RecurringDeleteConfirm
)

type RecurringField int

const (
	FieldName RecurringField = iota
	FieldCategory
	FieldScheduleType
	FieldMonths
	FieldEnabled
)

type RecurringListModel struct {
	store          storage.Storage
	payments       []storage.RecurringExpense
	categories     []string
	cursor         int
	mode           int
	editIndex      int
	editing        storage.RecurringExpense
	field          RecurringField
	inputValue     string
	categoryCursor int
	errorMsg       string
	successMsg     string
	quitting       bool
}

func NewRecurringListModel(store storage.Storage) RecurringListModel {
	m := RecurringListModel{
		store: store,
		mode:  RecurringView,
	}
	m.loadPayments()
	m.loadCategories()
	return m
}

func (m *RecurringListModel) loadPayments() {
	payments, err := m.store.GetRecurringExpenses()
	if err != nil {
		m.errorMsg = fmt.Sprintf("Failed to load payments: %v", err)
		return
	}
	m.payments = payments
}

func (m *RecurringListModel) loadCategories() {
	categories, err := m.store.GetCategories()
	if err != nil {
		m.categories = []string{"Uncategorized", "Food", "Groceries", "Travel", "Rent", "Utilities", "Entertainment", "Healthcare", "Shopping", "Miscellaneous"}
		return
	}
	m.categories = categories
}

func (m RecurringListModel) Init() tea.Cmd {
	return nil
}

func (m RecurringListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "esc":
			switch m.mode {
			case RecurringAdd, RecurringEdit:
				m.inputValue = ""
				m.field = FieldName
				m.mode = RecurringView
				m.editIndex = -1
				m.editing = storage.RecurringExpense{}
			case RecurringDeleteConfirm:
				m.mode = RecurringView
			case RecurringView:
				return m, func() tea.Msg {
					return NavigationMsg{Destination: ScreenSettings}
				}
			}
		case "up", "k":
			switch m.mode {
			case RecurringView:
				if m.cursor > 0 {
					m.cursor--
				}
			case RecurringAdd, RecurringEdit:
				if m.field == FieldCategory && m.categoryCursor > 0 {
					m.categoryCursor--
				}
			}
		case "down", "j":
			switch m.mode {
			case RecurringView:
				if m.cursor < len(m.payments) {
					m.cursor++
				}
			case RecurringAdd, RecurringEdit:
				if m.field == FieldCategory && m.categoryCursor < len(m.categories)-1 {
					m.categoryCursor++
				}
			}
		case "enter":
			switch m.mode {
			case RecurringView:
				if m.cursor < len(m.payments) {
					m.mode = RecurringEdit
					m.editIndex = m.cursor
					m.editing = m.payments[m.cursor]
					m.field = FieldName
					m.inputValue = m.editing.Name
					m.setCategoryCursor()
				} else if m.cursor == len(m.payments) {
					return m, func() tea.Msg {
						return NavigationMsg{Destination: ScreenSettings}
					}
				}
			case RecurringAdd, RecurringEdit:
				m.handleFieldInput()
			case RecurringDeleteConfirm:
				if m.cursor < len(m.payments) {
					m.removePayment(m.payments[m.cursor].ID)
					m.loadPayments()
					m.cursor = 0
					m.mode = RecurringView
				} else {
					m.mode = RecurringView
				}
			}
		case "y":
			if m.mode == RecurringDeleteConfirm && m.cursor < len(m.payments) {
				m.removePayment(m.payments[m.cursor].ID)
				m.loadPayments()
				m.cursor = 0
				m.mode = RecurringView
			} else if (m.mode == RecurringAdd || m.mode == RecurringEdit) && m.field != FieldCategory {
				m.inputValue += msg.String()
			}
		case "n":
			if m.mode == RecurringDeleteConfirm {
				m.mode = RecurringView
			} else if (m.mode == RecurringAdd || m.mode == RecurringEdit) && m.field != FieldCategory {
				m.inputValue += msg.String()
			}
		case " ":
			if (m.mode == RecurringAdd || m.mode == RecurringEdit) && m.field != FieldCategory {
				m.inputValue += " "
			}
		case "backspace":
			if (m.mode == RecurringAdd || m.mode == RecurringEdit) && m.field != FieldCategory && m.inputValue != "" {
				m.inputValue = m.inputValue[:len(m.inputValue)-1]
			}
		case "left", "right":
			// Ignore arrow keys in all input modes
		default:
			switch m.mode {
			case RecurringView:
				switch msg.String() {
				case "a":
					m.mode = RecurringAdd
					m.editing = storage.RecurringExpense{
						StartDate:   time.Now(),
						Interval:    "monthly",
						Occurrences: 12,
					}
					m.inputValue = ""
					m.field = FieldName
					m.categoryCursor = 0
				case "e":
					if m.cursor < len(m.payments) {
						m.mode = RecurringEdit
						m.editIndex = m.cursor
						m.editing = m.payments[m.cursor]
						m.field = FieldName
						m.inputValue = m.editing.Name
						m.setCategoryCursor()
					}
				case "x":
					if m.cursor < len(m.payments) {
						m.mode = RecurringDeleteConfirm
					}
				}
			case RecurringAdd, RecurringEdit:
				if m.field != FieldCategory && (msg.Type == tea.KeyRunes || msg.Type == tea.KeySpace) {
					m.inputValue += msg.String()
				}
			}
		}
	}
	return m, nil
}

func (m *RecurringListModel) setCategoryCursor() {
	for i, cat := range m.categories {
		if cat == m.editing.Category {
			m.categoryCursor = i
			return
		}
	}
	m.categoryCursor = 0
}

func (m *RecurringListModel) handleFieldInput() {
	switch m.field {
	case FieldName:
		m.editing.Name = m.inputValue
		if m.inputValue != "" {
			m.field = FieldCategory
			m.inputValue = ""
		}
	case FieldCategory:
		m.editing.Category = m.categories[m.categoryCursor]
		m.field = FieldScheduleType
		m.inputValue = m.editing.Interval
	case FieldScheduleType:
		if m.inputValue == "daily" || m.inputValue == "weekly" || m.inputValue == "monthly" || m.inputValue == "yearly" {
			m.editing.Interval = m.inputValue
			m.field = FieldEnabled
			if m.editing.Occurrences > 0 {
				m.inputValue = "y"
			} else {
				m.inputValue = "n"
			}
		}
	case FieldEnabled:
		if m.inputValue == "y" || m.inputValue == "Y" {
			m.editing.Occurrences = 12
		} else {
			m.editing.Occurrences = 0
		}
		m.savePayment()
	}
}

func (m *RecurringListModel) savePayment() {
	var err error
	switch m.mode {
	case RecurringAdd:
		err = m.store.AddRecurringExpense(m.editing)
		if err == nil {
			m.successMsg = "Payment added successfully"
		}
	case RecurringEdit:
		err = m.store.UpdateRecurringExpense(m.editing.ID, m.editing, false)
		if err == nil {
			m.successMsg = "Payment updated successfully"
		}
	}

	if err != nil {
		m.errorMsg = err.Error()
	} else {
		m.loadPayments()
	}

	m.inputValue = ""
	m.field = FieldName
	m.editIndex = -1
	m.mode = RecurringView
	m.editing = storage.RecurringExpense{}
}

func (m *RecurringListModel) removePayment(id string) {
	err := m.store.RemoveRecurringExpense(id, true)
	if err != nil {
		m.errorMsg = fmt.Sprintf("Failed to remove payment: %v", err)
	} else {
		m.successMsg = "Payment removed successfully"
	}
}

func monthsToString(months []int) string {
	var strs []string
	for _, m := range months {
		strs = append(strs, monthToString(m))
	}
	return strings.Join(strs, ",")
}

func stringToMonths(s string) []int {
	var months []int
	parts := strings.Split(s, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		m := monthToNumber(p)
		if m >= 1 && m <= 12 {
			months = append(months, m)
		}
	}
	return months
}

func monthToString(m int) string {
	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	if m >= 1 && m <= 12 {
		return months[m-1]
	}
	return ""
}

func monthToNumber(s string) int {
	s = strings.ToLower(s)
	months := map[string]int{
		"jan": 1, "january": 1,
		"feb": 2, "february": 2,
		"mar": 3, "march": 3,
		"apr": 4, "april": 4,
		"may": 5,
		"jun": 6, "june": 6,
		"jul": 7, "july": 7,
		"aug": 8, "august": 8,
		"sep": 9, "september": 9,
		"oct": 10, "october": 10,
		"nov": 11, "november": 11,
		"dec": 12, "december": 12,
	}
	return months[s]
}

func (m RecurringListModel) View() string {
	var s strings.Builder

	s.WriteString(settingsLogo + "\n\n")
	s.WriteString("Recurring Expenses\n")
	s.WriteString("─────────────────────────────\n")

	if m.errorMsg != "" {
		fmt.Fprintf(&s, "Error: %s\n\n", m.errorMsg)
		m.errorMsg = ""
	}
	if m.successMsg != "" {
		fmt.Fprintf(&s, "✓ %s\n\n", m.successMsg)
		m.successMsg = ""
	}

	switch m.mode {
	case RecurringAdd:
		m.renderForm(&s, "Add")
	case RecurringEdit:
		m.renderForm(&s, "Edit")
	case RecurringDeleteConfirm:
		s.WriteString("Confirm delete")
		if m.cursor < len(m.payments) {
			fmt.Fprintf(&s, " '%s'", m.payments[m.cursor].Name)
		}
		s.WriteString("?\n\n[Y]es / [N]o\n")
	default:
		m.renderList(&s)
	}

	return s.String()
}

func (m RecurringListModel) renderForm(s *strings.Builder, action string) {
	fmt.Fprintf(s, "%s expense:\n\n", action)

	if m.field == FieldName {
		fmt.Fprintf(s, "> Name: %s_\n\n", m.inputValue)
	} else {
		fmt.Fprintf(s, "  Name: %s\n", m.editing.Name)
	}

	if m.field == FieldCategory {
		fmt.Fprintf(s, "> Category:\n")
		for i, cat := range m.categories {
			prefix := "    "
			if m.categoryCursor == i {
				prefix = "  > "
			}
			fmt.Fprintf(s, "%s%s\n", prefix, cat)
		}
		fmt.Fprintf(s, "\n  (Use ↑↓ to select, Enter to confirm)\n")
	} else {
		fmt.Fprintf(s, "  Category: %s\n", m.editing.Category)
	}

	if m.field == FieldScheduleType {
		fmt.Fprintf(s, "\n> Interval (daily/weekly/monthly/yearly): %s_\n", m.inputValue)
	} else {
		fmt.Fprintf(s, "\n  Interval: %s\n", m.editing.Interval)
	}

	if m.field == FieldEnabled {
		fmt.Fprintf(s, "\n> Enable (y/n): %s_\n", m.inputValue)
	} else if m.field != FieldName && m.field != FieldCategory && m.field != FieldScheduleType {
		fmt.Fprintf(s, "\n  Occurrences: %d\n", m.editing.Occurrences)
	}

	s.WriteString("\nEnter Continue     Esc Cancel\n")
}

func (m RecurringListModel) renderList(s *strings.Builder) {
	fmt.Fprintf(s, "Payments (%d):\n\n", len(m.payments))

	if len(m.payments) == 0 {
		s.WriteString("  (no recurring payments defined)\n")
	} else {
		s.WriteString("  Name                   Category        Interval    Occurrences\n")
		s.WriteString("  ───────────────────────────────────────────────────────────────\n")

		for i, p := range m.payments {
			cursor := "  "
			if m.cursor == i {
				cursor = "> "
			}
			occurrences := p.Occurrences
			if occurrences == 0 {
				occurrences = 0
			}
			fmt.Fprintf(s, "%s%-22s %-15s %-12s %d\n",
				cursor, truncate(p.Name, 22), truncate(p.Category, 15), truncate(p.Interval, 12), occurrences)
		}
	}

	s.WriteString("\na Add     e Edit     x Delete     esc Back to Settings\n")
}

func truncate(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-2] + ".."
	}
	return s
}

func yesNo(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

func (m RecurringListModel) IsQuitting() bool {
	return m.quitting
}
