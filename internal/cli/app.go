package cli

import (
	"expensecat/internal/storage"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const defaultBasePath = "~/Downloads/Expenses-Mar-2026-Exports"

type AppModel struct {
	currentScreen int
	store         storage.Storage
	basePath      string
	mainMenu      MainMenuModel
	importModel   ImportModel
	reportModel   ReportModel
}

func NewAppModel(store storage.Storage) AppModel {
	return AppModel{
		currentScreen: ScreenMainMenu,
		store:         store,
		basePath:      getBasePath(),
		mainMenu:      NewMainMenuModel(),
		importModel:   NewImportModel(store, getBasePath()),
		reportModel:   NewReportModel(store),
	}
}

func getBasePath() string {
	path := os.Getenv("EXPENSE_BASE_PATH")
	if path == "" {
		return defaultBasePath
	}
	return path
}

func (m AppModel) Init() tea.Cmd {
	return m.mainMenu.Init()
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case NavigationMsg:
		m.currentScreen = msg.Destination
		switch msg.Destination {
		case ScreenImport:
			m.importModel = NewImportModel(m.store, m.basePath)
		case ScreenReport:
			m.reportModel = NewReportModel(m.store)
		}
		return m, nil

	case tea.KeyMsg:
		switch m.currentScreen {
		case ScreenMainMenu:
			newModel, cmd := m.mainMenu.Update(msg)
			m.mainMenu = newModel.(MainMenuModel)
			if m.mainMenu.IsQuitting() {
				return m, tea.Quit
			}
			return m, cmd

		case ScreenImport:
			newModel, cmd := m.importModel.Update(msg)
			m.importModel = newModel.(ImportModel)
			if m.importModel.IsQuitting() {
				return m, tea.Quit
			}
			return m, cmd

		case ScreenReport:
			newModel, cmd := m.reportModel.Update(msg)
			m.reportModel = newModel.(ReportModel)
			if m.reportModel.IsQuitting() {
				return m, tea.Quit
			}
			return m, cmd
		}
	}

	switch m.currentScreen {
	case ScreenMainMenu:
		newModel, _ := m.mainMenu.Update(msg)
		m.mainMenu = newModel.(MainMenuModel)
	case ScreenImport:
		newModel, _ := m.importModel.Update(msg)
		m.importModel = newModel.(ImportModel)
	case ScreenReport:
		newModel, _ := m.reportModel.Update(msg)
		m.reportModel = newModel.(ReportModel)
	}

	return m, nil
}

func (m AppModel) View() string {
	switch m.currentScreen {
	case ScreenMainMenu:
		return m.mainMenu.View()
	case ScreenImport:
		return m.importModel.View()
	case ScreenReport:
		return m.reportModel.View()
	default:
		return "Unknown screen\n"
	}
}

func RunApp(store storage.Storage) error {
	model := NewAppModel(store)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running application: %w", err)
	}

	return nil
}
