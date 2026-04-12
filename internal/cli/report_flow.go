package cli

import (
	"expensecat/internal/storage"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type ReportModel struct {
	storage      storage.Storage
	cursor       int
	inputMonth   string
	inputYear    string
	inputtingFor int
	showResults  bool
	results      []storage.CategoryTotal
	total        float64
	month        int
	year         int
	errorMsg     string
	quitting     bool
}

func NewReportModel(store storage.Storage) ReportModel {
	currentYear := time.Now().Year()
	return ReportModel{
		storage:   store,
		inputYear: strconv.Itoa(currentYear),
	}
}

func (m ReportModel) Init() tea.Cmd {
	return nil
}

func (m ReportModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if !m.showResults {
				if m.cursor > 0 {
					m.cursor--
				}
			}
		case "down", "j":
			if !m.showResults {
				if m.cursor < 3 {
					m.cursor++
				}
			}
		case "enter":
			if m.showResults {
				return m, func() tea.Msg {
					return NavigationMsg{Destination: ScreenMainMenu}
				}
			}
			if m.cursor == 0 {
				m.inputtingFor = 1
			} else if m.cursor == 1 {
				m.inputtingFor = 2
			} else if m.cursor == 2 {
				m.executeReport()
			}
		case "backspace":
			if m.inputtingFor == 1 && len(m.inputMonth) > 0 {
				m.inputMonth = m.inputMonth[:len(m.inputMonth)-1]
			} else if m.inputtingFor == 2 && len(m.inputYear) > 0 {
				m.inputYear = m.inputYear[:len(m.inputYear)-1]
			}
		default:
			if !m.showResults {
				if m.inputtingFor == 1 {
					m.inputMonth += msg.String()
				} else if m.inputtingFor == 2 {
					m.inputYear += msg.String()
				}
			}
		}
	}
	return m, nil
}

func (m *ReportModel) executeReport() {
	m.errorMsg = ""

	month, err := strconv.Atoi(strings.TrimSpace(m.inputMonth))
	if err != nil || month < 1 || month > 12 {
		m.errorMsg = "Invalid month (must be 1-12)"
		return
	}

	year, err := strconv.Atoi(strings.TrimSpace(m.inputYear))
	if err != nil || year < 1 {
		m.errorMsg = "Invalid year"
		return
	}

	results, err := storage.GetExpensesByMonthGroupedByCategory(m.storage, year, month)
	if err != nil {
		m.errorMsg = fmt.Sprintf("Failed to get expenses: %v", err)
		return
	}

	m.results = results
	m.month = month
	m.year = year
	m.total = 0
	for _, r := range results {
		m.total += r.Total
	}
	m.showResults = true
}

func (m ReportModel) View() string {
	var s strings.Builder

	if m.showResults {
		monthName := time.Month(m.month).String()
		s.WriteString(fmt.Sprintf("%s %d Totals\n", monthName, m.year))
		s.WriteString("─────────────────────────────\n")

		for _, r := range m.results {
			s.WriteString(fmt.Sprintf("%-15s $%9.2f\n", r.Category, r.Total))
		}

		s.WriteString("─────────────────────────────\n")
		s.WriteString(fmt.Sprintf("%-15s $%9.2f\n", "TOTAL", m.total))
		s.WriteString("\nPress any key to return\n")
		return s.String()
	}

	s.WriteString("Monthly Report\n")
	s.WriteString("─────────────────────────────\n")

	monthDisplay := m.inputMonth
	if monthDisplay == "" {
		monthDisplay = "  "
	}
	s.WriteString(fmt.Sprintf("Enter month (1-12): %s\n", monthDisplay))

	yearDisplay := m.inputYear
	if yearDisplay == "" {
		yearDisplay = "    "
	}
	s.WriteString(fmt.Sprintf("Enter year (e.g. 2024): %s\n", yearDisplay))

	if m.errorMsg != "" {
		s.WriteString(fmt.Sprintf("\nError: %s\n", m.errorMsg))
	}

	s.WriteString("\n")

	cursor := "  "
	if m.cursor == 0 {
		cursor = "> "
	}
	s.WriteString(fmt.Sprintf("%s[Month]\n", cursor))

	cursor = "  "
	if m.cursor == 1 {
		cursor = "> "
	}
	s.WriteString(fmt.Sprintf("%s[Year]\n", cursor))

	cursor = "  "
	if m.cursor == 2 {
		cursor = "> "
	}
	s.WriteString(fmt.Sprintf("%s[View Report]\n", cursor))

	cursor = "  "
	if m.cursor == 3 {
		cursor = "> "
	}
	s.WriteString(fmt.Sprintf("%s[Back]\n", cursor))

	s.WriteString("\nArrow keys navigate, Enter to select, Esc to go back\n")

	return s.String()
}

func (m ReportModel) IsQuitting() bool {
	return m.quitting
}

func SortCategoryTotals(totals []storage.CategoryTotal) []storage.CategoryTotal {
	sorted := make([]storage.CategoryTotal, len(totals))
	copy(sorted, totals)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Category < sorted[j].Category
	})
	return sorted
}
