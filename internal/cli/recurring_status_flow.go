package cli

import (
	"expensecat/internal/storage"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var recurringStatusLogo = ` /\_/\  
( $.$ )  
(  >💳 )  ExpenseCat`

type RecurringStatusModel struct {
	storage     storage.Storage
	cursor      int
	showResults bool
	month       int
	year        int
	results     []RecurringStatusItem
	errorMsg    string
	quitting    bool
}

type RecurringStatusItem struct {
	Name   string
	Amount float64
	Paid   bool
}

func NewRecurringStatusModel(store storage.Storage) RecurringStatusModel {
	now := time.Now()
	month := int(now.Month()) - 1
	year := now.Year()
	if month < 1 {
		month = 12
		year--
	}
	return RecurringStatusModel{
		storage: store,
		month:   month,
		year:    year,
	}
}

func (m RecurringStatusModel) Init() tea.Cmd {
	return nil
}

func (m RecurringStatusModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "left", "h":
			m.navigateMonth(-1)
		case "right", "l":
			m.navigateMonth(1)
		case "up", "k":
			if m.showResults && m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.showResults && m.cursor < len(m.results) {
				m.cursor++
			}
		case "enter":
			if m.showResults {
				return m, func() tea.Msg {
					return NavigationMsg{Destination: ScreenMainMenu}
				}
			} else {
				m.executeQuery()
			}
		}
	}
	return m, nil
}

func (m *RecurringStatusModel) navigateMonth(delta int) {
	m.month += delta
	if m.month > 12 {
		m.month = 1
		m.year++
	} else if m.month < 1 {
		m.month = 12
		m.year--
	}
	if m.showResults {
		m.executeQuery()
	}
}

func (m *RecurringStatusModel) executeQuery() {
	m.errorMsg = ""
	m.showResults = false
	m.results = nil

	recurringExpenses, err := m.storage.GetRecurringExpenses()
	if err != nil {
		m.errorMsg = fmt.Sprintf("Failed to get recurring expenses: %v", err)
		return
	}

	startDate := time.Date(m.year, time.Month(m.month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	expenses, err := m.storage.GetAllExpenses(&startDate, &endDate)
	if err != nil {
		m.errorMsg = fmt.Sprintf("Failed to get expenses: %v", err)
		return
	}

	for _, rec := range recurringExpenses {
		item := RecurringStatusItem{
			Name:   rec.Name,
			Amount: rec.Amount,
			Paid:   false,
		}

		// Check if there's an expense with a name containing the recurring expense name
		for _, exp := range expenses {
			if strings.Contains(strings.ToLower(exp.Name), strings.ToLower(rec.Name)) {
				item.Paid = true
				item.Amount = exp.Amount // Use actual paid amount
				break
			}
		}

		m.results = append(m.results, item)
	}

	m.showResults = true
}

func (m RecurringStatusModel) View() string {
	var s strings.Builder

	s.WriteString(recurringStatusLogo + "\n\n")

	if m.showResults {
		monthName := time.Month(m.month).String()
		fmt.Fprintf(&s, "%s %d Recurring Expenses\n", monthName, m.year)
		s.WriteString("───────────────────────────────────────\n")

		if len(m.results) == 0 {
			s.WriteString("  (no recurring expenses defined)\n")
		} else {
			for i, r := range m.results {
				cursor := "  "
				if m.cursor == i {
					cursor = "> "
				}
				status := "unpaid"
				if r.Paid {
					status = "paid"
				}
				fmt.Fprintf(&s, "%s%-15s %7.2f %s\n", cursor, truncateStatus(r.Name, 15), r.Amount, status)
			}
		}

		s.WriteString("\n←→ Previous/Next Month\n")
		s.WriteString("Enter Return to Menu\n")
		return s.String()
	}

	m.executeQuery()

	monthName := time.Month(m.month).String()
	fmt.Fprintf(&s, "%s %d Recurring Expenses\n", monthName, m.year)
	s.WriteString("───────────────────────────────────────\n")

	if m.errorMsg != "" {
		fmt.Fprintf(&s, "Error: %s\n", m.errorMsg)
		return s.String()
	}

	if len(m.results) == 0 {
		s.WriteString("  (no recurring expenses defined)\n")
	} else {
		for i, r := range m.results {
			cursor := "  "
			if m.cursor == i {
				cursor = "> "
			}
			status := "unpaid"
			if r.Paid {
				status = "paid"
			}
			fmt.Fprintf(&s, "%s%-15s %7.2f %s\n", cursor, truncateStatus(r.Name, 15), r.Amount, status)
		}
	}

	s.WriteString("\n←→ Navigate Months  Esc to Menu\n")

	return s.String()
}

func truncateStatus(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-2] + ".."
	}
	return s
}

func (m RecurringStatusModel) IsQuitting() bool {
	return m.quitting
}
