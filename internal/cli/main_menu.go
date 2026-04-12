package cli

import (
	tea "github.com/charmbracelet/bubbletea"
)

const (
	MenuImport = iota
	MenuReport
	MenuQuit
)

type MainMenuModel struct {
	cursor   int
	choices  []string
	quitting bool
}

func NewMainMenuModel() MainMenuModel {
	return MainMenuModel{
		choices: []string{
			"Import Expenses",
			"View Monthly Totals",
			"Quit",
		},
	}
}

func (m MainMenuModel) Init() tea.Cmd {
	return nil
}

func (m MainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			if m.cursor == MenuImport {
				return m, func() tea.Msg {
					return NavigationMsg{Destination: ScreenImport}
				}
			}
			if m.cursor == MenuReport {
				return m, func() tea.Msg {
					return NavigationMsg{Destination: ScreenReport}
				}
			}
			if m.cursor == MenuQuit {
				m.quitting = true
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m MainMenuModel) View() string {
	s := "ExpenseCat\n"
	s += "─────────────────────────────\n"

	for i, choice := range m.choices {
		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}
		s += cursor + choice + "\n"
	}

	s += "\n↑↓ Navigate  Enter Select\n"
	s += "q Quit\n"

	return s
}

func (m MainMenuModel) IsQuitting() bool {
	return m.quitting
}
