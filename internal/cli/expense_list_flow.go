package cli

import (
	"expensecat/internal/storage"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	ExpenseListBrowse = iota
	ExpenseListSearch
	ExpenseListEdit
	ExpenseListDeleteConfirm
)

var expenseListLogo = ` /\_/\  
( $.$ )  
(  >💳 )  ExpenseCat`

const expensesPerPage = 10

type ExpenseListModel struct {
	storage            storage.Storage
	expenses           []storage.Expense
	filteredExpenses   []storage.Expense
	cursor             int
	page               int
	categories         []string
	searchQuery        string
	categoryFilter     string
	startDateFilter    string
	endDateFilter      string
	mode               int
	editingExpense     storage.Expense
	editField          int
	editValue          string
	showCategoryPicker bool
	categoryPickerIdx  int
	errorMsg           string
	successMsg         string
	quitting           bool
}

func NewExpenseListModel(store storage.Storage) ExpenseListModel {
	m := ExpenseListModel{
		storage: store,
		mode:    ExpenseListBrowse,
	}
	m.loadData()
	return m
}

func (m *ExpenseListModel) loadData() {
	categories, err := m.storage.GetCategories()
	if err == nil {
		m.categories = categories
	}

	expenses, err := m.storage.GetAllExpenses(nil, nil)
	if err != nil {
		m.errorMsg = fmt.Sprintf("Failed to load expenses: %v", err)
		m.expenses = []storage.Expense{}
	} else {
		m.expenses = expenses
	}
	m.applyFilters()
}

func (m *ExpenseListModel) applyFilters() {
	m.filteredExpenses = []storage.Expense{}

	for _, exp := range m.expenses {
		if m.searchQuery != "" {
			if !strings.Contains(strings.ToLower(exp.Name), strings.ToLower(m.searchQuery)) {
				continue
			}
		}

		if m.categoryFilter != "" && exp.Category != m.categoryFilter {
			continue
		}

		if m.startDateFilter != "" {
			startDate, err := parseDate(m.startDateFilter)
			if err == nil && exp.Date.Before(startDate) {
				continue
			}
		}

		if m.endDateFilter != "" {
			endDate, err := parseDate(m.endDateFilter)
			if err == nil && exp.Date.After(endDate) {
				continue
			}
		}

		m.filteredExpenses = append(m.filteredExpenses, exp)
	}

	if m.page*m.expensesPerPage() >= len(m.filteredExpenses) {
		m.page = 0
	}
	if m.cursor >= len(m.filteredExpenses) {
		m.cursor = len(m.filteredExpenses) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

func parseDate(s string) (time.Time, error) {
	formats := []string{"2006-01-02", "01/02/2006", "01-02-2006", "2006/01/02"}
	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date format")
}

func (m *ExpenseListModel) expensesPerPage() int {
	return expensesPerPage
}

func (m *ExpenseListModel) totalPages() int {
	if len(m.filteredExpenses) == 0 {
		return 1
	}
	return (len(m.filteredExpenses) + expensesPerPage - 1) / expensesPerPage
}

func (m *ExpenseListModel) currentPageExpenses() []storage.Expense {
	start := m.page * m.expensesPerPage()
	end := start + m.expensesPerPage()
	if end > len(m.filteredExpenses) {
		end = len(m.filteredExpenses)
	}
	if start >= len(m.filteredExpenses) {
		return []storage.Expense{}
	}
	return m.filteredExpenses[start:end]
}

func (m ExpenseListModel) Init() tea.Cmd {
	return nil
}

func (m ExpenseListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.mode {
		case ExpenseListBrowse:
			return m.updateBrowseMode(msg)
		case ExpenseListSearch:
			return m.updateSearchMode(msg)
		case ExpenseListEdit:
			return m.updateEditMode(msg)
		case ExpenseListDeleteConfirm:
			return m.updateDeleteConfirmMode(msg)
		}
	}
	return m, nil
}

func (m *ExpenseListModel) updateBrowseMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.quitting = true
		return m, tea.Quit
	case "esc":
		return m, func() tea.Msg {
			return NavigationMsg{Destination: ScreenMainMenu}
		}
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
			if m.cursor < m.page*m.expensesPerPage() {
				m.page--
			}
		}
	case "down", "j":
		if m.cursor < len(m.filteredExpenses)-1 {
			m.cursor++
			if m.cursor >= (m.page+1)*m.expensesPerPage() {
				m.page++
			}
		}
	case "/":
		m.mode = ExpenseListSearch
		m.searchQuery = ""
		m.categoryFilter = ""
		m.startDateFilter = ""
		m.endDateFilter = ""
	case "c":
		m.mode = ExpenseListSearch
		m.categoryFilter = ""
	case "d":
		m.mode = ExpenseListSearch
		m.startDateFilter = ""
		m.endDateFilter = ""
	case "enter", "e":
		if len(m.filteredExpenses) > 0 {
			m.editingExpense = m.filteredExpenses[m.cursor]
			m.mode = ExpenseListEdit
			m.editField = 0
			m.editValue = m.filteredExpenses[m.cursor].Name
		}
	case "x", "D":
		if len(m.filteredExpenses) > 0 {
			m.editingExpense = m.filteredExpenses[m.cursor]
			m.mode = ExpenseListDeleteConfirm
		}
	case "left", "h":
		if m.page > 0 {
			m.cursor -= expensesPerPage
			if m.cursor < 0 {
				m.cursor = 0
			}
			m.page--
		}
	case "right", "l":
		if m.page < m.totalPages()-1 {
			m.cursor += expensesPerPage
			if m.cursor >= len(m.filteredExpenses) {
				m.cursor = len(m.filteredExpenses) - 1
			}
			m.page++
		}
	}
	return m, nil
}

func (m *ExpenseListModel) updateSearchMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		m.quitting = true
		return m, tea.Quit
	case "esc":
		m.mode = ExpenseListBrowse
		m.applyFilters()
	case "enter":
		m.mode = ExpenseListBrowse
		m.applyFilters()
	case "backspace":
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
		}
	default:
		if len(msg.Runes) > 0 {
			m.searchQuery += msg.String()
		}
	}
	return m, nil
}

func (m *ExpenseListModel) updateEditMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.showCategoryPicker {
		switch msg.String() {
		case "up", "k":
			if m.categoryPickerIdx > 0 {
				m.categoryPickerIdx--
			} else {
				m.categoryPickerIdx = len(m.categories) - 1
			}
		case "down", "j":
			if m.categoryPickerIdx < len(m.categories)-1 {
				m.categoryPickerIdx++
			} else {
				m.categoryPickerIdx = 0
			}
		case "enter":
			m.editValue = m.categories[m.categoryPickerIdx]
			m.showCategoryPicker = false
		case "esc", "tab":
			m.showCategoryPicker = false
		}
		return m, nil
	}

	switch msg.String() {
	case "ctrl+c":
		m.quitting = true
		return m, tea.Quit
	case "esc":
		m.mode = ExpenseListBrowse
		m.errorMsg = ""
	case "enter":
		return m.handleEditSave()
	case "tab":
		if m.editField == 1 && len(m.categories) > 0 {
			m.showCategoryPicker = true
			currentCat := m.editingExpense.Category
			for i, cat := range m.categories {
				if cat == currentCat {
					m.categoryPickerIdx = i
					break
				}
			}
		} else {
			m.editField = (m.editField + 1) % 5
			m.editValue = m.getEditFieldValue()
		}
	case "up", "k":
		m.editField = (m.editField - 1 + 5) % 5
		m.editValue = m.getEditFieldValue()
	case "down", "j":
		m.editField = (m.editField + 1) % 5
		m.editValue = m.getEditFieldValue()
	case "backspace":
		if len(m.editValue) > 0 {
			m.editValue = m.editValue[:len(m.editValue)-1]
		}
	default:
		if len(msg.Runes) > 0 {
			m.editValue += msg.String()
		}
	}
	return m, nil
}

func (m *ExpenseListModel) handleEditSave() (tea.Model, tea.Cmd) {
	exp := m.editingExpense

	switch m.editField {
	case 0:
		exp.Name = storage.SanitizeString(m.editValue)
	case 1:
		exp.Category = storage.SanitizeString(m.editValue)
	case 2:
		var amount float64
		fmt.Sscanf(m.editValue, "%f", &amount)
		exp.Amount = amount
	case 3:
		exp.Currency = strings.ToLower(storage.SanitizeString(m.editValue))
	case 4:
		if t, err := parseDate(m.editValue); err == nil {
			exp.Date = t
		}
	}

	if err := exp.Validate(); err != nil {
		m.errorMsg = err.Error()
		return m, nil
	}

	if err := m.storage.UpdateExpense(exp.ID, exp); err != nil {
		m.errorMsg = fmt.Sprintf("Failed to update expense: %v", err)
		return m, nil
	}

	m.loadData()
	m.mode = ExpenseListBrowse
	m.successMsg = "Expense updated successfully"
	return m, nil
}

func (m *ExpenseListModel) getEditFieldValue() string {
	exp := m.editingExpense
	switch m.editField {
	case 0:
		return exp.Name
	case 1:
		return exp.Category
	case 2:
		return fmt.Sprintf("%.2f", exp.Amount)
	case 3:
		return exp.Currency
	case 4:
		return exp.Date.Format("2006-01-02")
	}
	return ""
}

func (m ExpenseListModel) getDisplayValue(fieldIndex int, originalValue, fallbackValue string) string {
	if m.editField == fieldIndex {
		if m.editValue != "" {
			return m.editValue
		}
		return originalValue
	}
	if originalValue != "" {
		return originalValue
	}
	return fallbackValue
}

func (m *ExpenseListModel) updateDeleteConfirmMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		m.quitting = true
		return m, tea.Quit
	case "esc":
		m.mode = ExpenseListBrowse
	case "y", "Y", "enter":
		if err := m.storage.RemoveExpense(m.editingExpense.ID); err != nil {
			m.errorMsg = fmt.Sprintf("Failed to delete expense: %v", err)
			m.mode = ExpenseListBrowse
			return m, nil
		}
		m.loadData()
		m.mode = ExpenseListBrowse
		m.successMsg = "Expense deleted successfully"
	case "n", "N":
		m.mode = ExpenseListBrowse
	}
	return m, nil
}

func (m ExpenseListModel) View() string {
	switch m.mode {
	case ExpenseListBrowse:
		return m.viewBrowse()
	case ExpenseListSearch:
		return m.viewSearch()
	case ExpenseListEdit:
		return m.viewEdit()
	case ExpenseListDeleteConfirm:
		return m.viewDeleteConfirm()
	}
	return ""
}

func (m ExpenseListModel) viewBrowse() string {
	var s strings.Builder

	s.WriteString(expenseListLogo + "\n\n")
	s.WriteString("Manage Expenses\n")
	s.WriteString("─────────────────────────────\n\n")

	if m.successMsg != "" {
		s.WriteString("\x1b[32m") // green
		s.WriteString(m.successMsg)
		s.WriteString("\x1b[0m\n\n")
	}

	if m.errorMsg != "" {
		s.WriteString("\x1b[31m") // red
		s.WriteString(m.errorMsg)
		s.WriteString("\x1b[0m\n\n")
	}

	s.WriteString("Filters: ")
	if m.searchQuery != "" {
		fmt.Fprintf(&s, "search='%s' ", m.searchQuery)
	}
	if m.categoryFilter != "" {
		fmt.Fprintf(&s, "category='%s' ", m.categoryFilter)
	}
	if m.startDateFilter != "" {
		fmt.Fprintf(&s, "from=%s ", m.startDateFilter)
	}
	if m.endDateFilter != "" {
		fmt.Fprintf(&s, "to=%s ", m.endDateFilter)
	}
	s.WriteString("\n")

	fmt.Fprintf(&s, "Showing %d of %d expenses | Page %d/%d\n",
		len(m.filteredExpenses), len(m.expenses), m.page+1, m.totalPages())
	s.WriteString("─────────────────────────────\n")

	pageExpenses := m.currentPageExpenses()
	for i, exp := range pageExpenses {
		globalIdx := m.page*m.expensesPerPage() + i
		cursor := "  "
		if globalIdx == m.cursor {
			cursor = "> "
		}
		dateStr := exp.Date.Format("2006-01-02")
		fmt.Fprintf(&s, "%s%-20s %10.2f %-12s %s\n", cursor, trunc(exp.Name, 20), exp.Amount, exp.Category, dateStr)
	}

	if len(m.filteredExpenses) == 0 {
		s.WriteString("\nNo expenses found.\n")
		s.WriteString("Import some expenses or adjust filters.\n")
	}

	s.WriteString("\n─────────────────────────────\n")
	s.WriteString("/ Search   c Category   d Dates\n")
	s.WriteString("e Edit     x Delete     ←→ Page\n")
	s.WriteString("esc Back   q Quit\n")

	return s.String()
}

func (m ExpenseListModel) viewSearch() string {
	var s strings.Builder

	s.WriteString(expenseListLogo + "\n\n")
	s.WriteString("Search Expenses\n")
	s.WriteString("─────────────────────────────\n\n")

	s.WriteString("Search: ")
	s.WriteString(m.searchQuery)
	s.WriteString("_\n\n")

	s.WriteString("Press Enter to search or Esc to cancel\n")

	return s.String()
}

func (m ExpenseListModel) viewEdit() string {
	var s strings.Builder

	s.WriteString(expenseListLogo + "\n\n")
	s.WriteString("Edit Expense\n")
	s.WriteString("─────────────────────────────\n\n")

	exp := m.editingExpense

	fields := []struct {
		label string
		value string
		edit  bool
	}{
		{"Name", m.getDisplayValue(0, exp.Name, ""), m.editField == 0},
		{"Category", m.getDisplayValue(1, exp.Category, ""), m.editField == 1},
		{"Amount", m.getDisplayValue(2, "", fmt.Sprintf("%.2f", exp.Amount)), m.editField == 2},
		{"Currency", m.getDisplayValue(3, exp.Currency, ""), m.editField == 3},
		{"Date", m.getDisplayValue(4, exp.Date.Format("2006-01-02"), ""), m.editField == 4},
	}

	for i, f := range fields {
		cursor := "  "
		if m.editField == i {
			cursor = "> "
		}
		fmt.Fprintf(&s, "%s%-10s: %s", cursor, f.label, f.value)
		if f.edit {
			s.WriteString(" _")
		}
		s.WriteString("\n")
	}

	if m.errorMsg != "" {
		s.WriteString("\n\x1b[31m") // red
		s.WriteString(m.errorMsg)
		s.WriteString("\x1b[0m")
	}

	if m.showCategoryPicker {
		s.WriteString(m.viewCategoryPicker())
	}

	s.WriteString("\n─────────────────────────────\n")
	if m.editField == 1 {
		s.WriteString("Tab Select category\n")
	} else {
		s.WriteString("Tab Next field  Type to edit\n")
	}
	s.WriteString("↑↓ Select field  Enter Save\n")
	s.WriteString("Esc Cancel\n")

	return s.String()
}

func (m ExpenseListModel) viewCategoryPicker() string {
	var s strings.Builder

	s.WriteString("\n\n")
	s.WriteString("\x1b[7m") // inverse
	s.WriteString(" Select Category ")
	s.WriteString("\x1b[0m") // reset
	s.WriteString("\n")

	for i, cat := range m.categories {
		cursor := "  "
		if i == m.categoryPickerIdx {
			cursor = "> "
		}
		fmt.Fprintf(&s, "%s%s\n", cursor, cat)
	}

	s.WriteString("\n\x1b[90m") // dim
	s.WriteString("↑↓ Select  Enter Choose  Esc Cancel\x1b[0m")

	return s.String()
}

func (m ExpenseListModel) viewDeleteConfirm() string {
	var s strings.Builder

	s.WriteString(expenseListLogo + "\n\n")
	s.WriteString("Delete Expense?\n")
	s.WriteString("─────────────────────────────\n\n")

	exp := m.editingExpense
	fmt.Fprintf(&s, "Name:     %s\n", exp.Name)
	fmt.Fprintf(&s, "Amount:   %.2f %s\n", exp.Amount, exp.Currency)
	fmt.Fprintf(&s, "Category: %s\n", exp.Category)
	fmt.Fprintf(&s, "Date:     %s\n", exp.Date.Format("2006-01-02"))

	s.WriteString("\n─────────────────────────────\n")
	s.WriteString("y Yes, delete   n No, cancel\n")

	return s.String()
}

func trunc(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-2] + ".."
}

func (m ExpenseListModel) IsQuitting() bool {
	return m.quitting
}
