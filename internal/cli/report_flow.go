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

var reportLogo = ` /\_/\  
( $.$ )  
(  >💳 )  ExpenseCat`

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
	now := time.Now()
	month := int(now.Month()) - 1
	year := now.Year()
	if month < 1 {
		month = 12
		year--
	}
	return ReportModel{
		storage:    store,
		month:    month,
		year:     year,
		showResults: true,
	}
}

func (m ReportModel) Init() tea.Cmd {
	m.executeReport()
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
		case "left", "right":
			if m.showResults {
				delta := 1
				if msg.String() == "left" {
					delta = -1
				}
				m.navigateMonth(delta)
			}
		case "up", "k":
			if m.showResults {
				if m.cursor > 0 {
					m.cursor--
				} else {
					m.cursor = 0
				}
			}
		case "down", "j":
			if m.showResults {
				if m.cursor < len(m.results)-1 {
					m.cursor++
				} else {
					m.cursor = 0
				}
			}
		case "enter":
			if m.showResults {
				return m, func() tea.Msg {
					return NavigationMsg{Destination: ScreenMainMenu}
				}
			}
		}
	}
	return m, nil
}

func (m *ReportModel) executeReport() {
	m.errorMsg = ""

	results, err := storage.GetExpensesByMonthGroupedByCategory(m.storage, m.year, m.month)
	if err != nil {
		m.errorMsg = fmt.Sprintf("Failed to get expenses: %v", err)
		return
	}

	m.results = results
	m.total = 0
	for _, r := range results {
		m.total += r.Total
	}
	m.showResults = true
}

func (m *ReportModel) navigateMonth(delta int) {
	m.month += delta
	if m.month > 12 {
		m.month = 1
		m.year++
	} else if m.month < 1 {
		m.month = 12
		m.year--
	}
	m.inputMonth = strconv.Itoa(m.month)
	m.inputYear = strconv.Itoa(m.year)
	m.executeReport()
}

func (m ReportModel) View() string {
	var s strings.Builder

	s.WriteString(reportLogo + "\n\n")

	if m.showResults {
		if m.errorMsg != "" {
			fmt.Fprintf(&s, "Error: %s\n", m.errorMsg)
			return s.String()
		}

		monthName := time.Month(m.month).String()
		fmt.Fprintf(&s, "%s %d Totals\n", monthName, m.year)
		s.WriteString("─────────────────────────────\n")

		for _, r := range m.results {
			fmt.Fprintf(&s, "%-15s $%9.2f\n", r.Category, r.Total)
		}

		s.WriteString("─────────────────────────────\n")
		fmt.Fprintf(&s, "%-15s $%9.2f\n", "TOTAL", m.total)
		s.WriteString("\n←→ Previous/Next Month  Esc to Menu\n")
		return s.String()
	}

	return ""
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
