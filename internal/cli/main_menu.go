package cli

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var catFrames = []string{
	` /\_/\  
( $.$ )  
(  >💳 )  ExpenseCat`,

	` /\_/\  
( -.- )  
(  >💳 )  ExpenseCat`,

	` /\_/\  
( $.$ )  
(  >💰 )  ExpenseCat`,

	` /\_/\  
( o.o )  
(  >💳 )  ExpenseCat`,
}

type tickMsg time.Time

const animationTickRate = 200 * time.Millisecond

const (
	MenuImport = iota
	MenuExpenseList
	MenuReport
	MenuSettings
	MenuQuit
)

type MainMenuModel struct {
	cursor     int
	choices    []string
	quitting   bool
	frameIndex int
	tickCount  int
}

func NewMainMenuModel() MainMenuModel {
	return MainMenuModel{
		choices: []string{
			"Import Expenses",
			"Manage Expenses",
			"View Monthly Totals",
			"Settings",
			"Quit",
		},
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(animationTickRate, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m MainMenuModel) Init() tea.Cmd {
	return tickCmd()
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
			if m.cursor == MenuExpenseList {
				return m, func() tea.Msg {
					return NavigationMsg{Destination: ScreenExpenseList}
				}
			}
			if m.cursor == MenuReport {
				return m, func() tea.Msg {
					return NavigationMsg{Destination: ScreenReport}
				}
			}
			if m.cursor == MenuSettings {
				return m, func() tea.Msg {
					return NavigationMsg{Destination: ScreenSettings}
				}
			}
			if m.cursor == MenuQuit {
				m.quitting = true
				return m, tea.Quit
			}
		}
	case tickMsg:
		m.tickCount++
		if m.tickCount%4 == 0 {
			m.frameIndex = (m.frameIndex + 1) % len(catFrames)
		}
		return m, tickCmd()
	}
	return m, nil
}

func (m MainMenuModel) View() string {
	s := catFrames[m.frameIndex] + "\n\n"
	s += "Main Menu\n"
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
