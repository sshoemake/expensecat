package cli

import (
	"expensecat/internal/storage"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const defaultBasePath = "~/Downloads/Expense-Files"

type AppModel struct {
	currentScreen   int
	store           storage.Storage
	basePath        string
	mainMenu        MainMenuModel
	importModel     ImportModel
	reportModel     ReportModel
	settingsModel   SettingsModel
	exclusionModel  ExclusionListModel
	importPathModel ImportPathModel
}

func NewAppModel(store storage.Storage) AppModel {
	return AppModel{
		currentScreen:   ScreenMainMenu,
		store:           store,
		basePath:        getBasePath(store),
		mainMenu:        NewMainMenuModel(),
		importModel:     NewImportModel(store, getBasePath(store)),
		reportModel:     NewReportModel(store),
		settingsModel:   NewSettingsModel(store),
		exclusionModel:  NewExclusionListModel(store),
		importPathModel: NewImportPathModel(store),
	}
}

func getBasePath(store storage.Storage) string {
	if store != nil {
		if path, err := store.GetImportPath(); err == nil && path != "" {
			return path
		}
	}
	path := os.Getenv("EXPENSE_BASE_PATH")
	if path != "" {
		return path
	}
	return defaultBasePath
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
			m.basePath = getBasePath(m.store)
			m.importModel = NewImportModel(m.store, m.basePath)
		case ScreenReport:
			m.reportModel = NewReportModel(m.store)
		case ScreenSettings:
			m.settingsModel = NewSettingsModel(m.store)
		case ScreenExclusionList:
			m.exclusionModel = NewExclusionListModel(m.store)
		case ScreenImportPath:
			m.importPathModel = NewImportPathModel(m.store)
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

		case ScreenSettings:
			newModel, cmd := m.settingsModel.Update(msg)
			m.settingsModel = newModel.(SettingsModel)
			if m.settingsModel.IsQuitting() {
				return m, tea.Quit
			}
			return m, cmd

		case ScreenExclusionList:
			newModel, cmd := m.exclusionModel.Update(msg)
			m.exclusionModel = newModel.(ExclusionListModel)
			if m.exclusionModel.IsQuitting() {
				return m, tea.Quit
			}
			return m, cmd

		case ScreenImportPath:
			newModel, cmd := m.importPathModel.Update(msg)
			m.importPathModel = newModel.(ImportPathModel)
			if m.importPathModel.IsQuitting() {
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
	case ScreenSettings:
		newModel, _ := m.settingsModel.Update(msg)
		m.settingsModel = newModel.(SettingsModel)
	case ScreenExclusionList:
		newModel, _ := m.exclusionModel.Update(msg)
		m.exclusionModel = newModel.(ExclusionListModel)
	case ScreenImportPath:
		newModel, _ := m.importPathModel.Update(msg)
		m.importPathModel = newModel.(ImportPathModel)
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
	case ScreenSettings:
		return m.settingsModel.View()
	case ScreenExclusionList:
		return m.exclusionModel.View()
	case ScreenImportPath:
		return m.importPathModel.View()
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
