package cli

import (
	"testing"
	"time"

	"expensecat/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type mockStoreForCLI struct{}

func (m *mockStoreForCLI) Close() error                        { return nil }
func (m *mockStoreForCLI) GetConfig() (*storage.Config, error) { return &storage.Config{}, nil }
func (m *mockStoreForCLI) GetCategories() ([]string, error)    { return []string{}, nil }
func (m *mockStoreForCLI) UpdateCategories([]string) error     { return nil }
func (m *mockStoreForCLI) GetCurrency() (string, error)        { return "usd", nil }
func (m *mockStoreForCLI) UpdateCurrency(string) error         { return nil }
func (m *mockStoreForCLI) GetStartDate() (int, error)          { return 1, nil }
func (m *mockStoreForCLI) UpdateStartDate(int) error           { return nil }
func (m *mockStoreForCLI) GetRecurringExpenses() ([]storage.RecurringExpense, error) {
	return []storage.RecurringExpense{}, nil
}
func (m *mockStoreForCLI) GetRecurringExpense(string) (storage.RecurringExpense, error) {
	return storage.RecurringExpense{}, nil
}
func (m *mockStoreForCLI) AddRecurringExpense(storage.RecurringExpense) error { return nil }
func (m *mockStoreForCLI) RemoveRecurringExpense(string, bool) error          { return nil }
func (m *mockStoreForCLI) UpdateRecurringExpense(string, storage.RecurringExpense, bool) error {
	return nil
}
func (m *mockStoreForCLI) GetAllExpenses(startDate, endDate *time.Time) ([]storage.Expense, error) {
	return []storage.Expense{}, nil
}
func (m *mockStoreForCLI) GetExpense(string) (storage.Expense, error)  { return storage.Expense{}, nil }
func (m *mockStoreForCLI) AddExpense(storage.Expense) error            { return nil }
func (m *mockStoreForCLI) RemoveExpense(string) error                  { return nil }
func (m *mockStoreForCLI) AddMultipleExpenses([]storage.Expense) error { return nil }
func (m *mockStoreForCLI) RemoveMultipleExpenses([]string) error       { return nil }
func (m *mockStoreForCLI) UpdateExpense(string, storage.Expense) error { return nil }

func TestMainMenuModel_Initialization(t *testing.T) {
	model := NewMainMenuModel()

	if len(model.choices) != 3 {
		t.Errorf("Expected 3 choices, got %d", len(model.choices))
	}
	if model.cursor != 0 {
		t.Errorf("Expected cursor to be 0, got %d", model.cursor)
	}
	if model.choices[0] != "Import Expenses" {
		t.Errorf("Expected first choice to be 'Import Expenses', got %s", model.choices[0])
	}
	if model.choices[1] != "View Monthly Totals" {
		t.Errorf("Expected second choice to be 'View Monthly Totals', got %s", model.choices[1])
	}
	if model.choices[2] != "Quit" {
		t.Errorf("Expected third choice to be 'Quit', got %s", model.choices[2])
	}
}

func TestImportModel_Initialization(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportModel(store, "./data")

	if model.basePath != "./data" {
		t.Errorf("Expected basePath to be './data', got %s", model.basePath)
	}
	if model.selectedFile != "" {
		t.Errorf("Expected selectedFile to be empty, got %s", model.selectedFile)
	}
	if model.cursor != 0 {
		t.Errorf("Expected cursor to be 0, got %d", model.cursor)
	}
	if model.showResults {
		t.Error("Expected showResults to be false")
	}
}

func TestReportModel_Initialization(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewReportModel(store)

	if model.inputMonth != "" {
		t.Errorf("Expected inputMonth to be empty, got %s", model.inputMonth)
	}
	currentYear := time.Now().Year()
	expectedYear := time.Date(currentYear, time.January, 1, 0, 0, 0, 0, time.UTC)
	if model.inputYear != expectedYear.Format("2006") {
		t.Errorf("Expected inputYear to be %d, got %s", currentYear, model.inputYear)
	}
	if model.cursor != 0 {
		t.Errorf("Expected cursor to be 0, got %d", model.cursor)
	}
	if model.showResults {
		t.Error("Expected showResults to be false")
	}
}

func TestSortCategoryTotals(t *testing.T) {
	totals := []storage.CategoryTotal{
		{Category: "Zoo", Total: 10.00},
		{Category: "Food", Total: 5.00},
		{Category: "Banking", Total: 15.00},
	}

	sorted := SortCategoryTotals(totals)

	expected := []string{"Banking", "Food", "Zoo"}
	for i, exp := range expected {
		if sorted[i].Category != exp {
			t.Errorf("Expected category %s at index %d, got %s", exp, i, sorted[i].Category)
		}
	}
}

func TestNavigationConstants(t *testing.T) {
	if ScreenMainMenu != 0 {
		t.Errorf("Expected ScreenMainMenu to be 0, got %d", ScreenMainMenu)
	}
	if ScreenImport != 1 {
		t.Errorf("Expected ScreenImport to be 1, got %d", ScreenImport)
	}
	if ScreenReport != 2 {
		t.Errorf("Expected ScreenReport to be 2, got %d", ScreenReport)
	}
}

func TestMainMenuModel_Navigation(t *testing.T) {
	model := NewMainMenuModel()

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = newModel.(MainMenuModel)
	if model.cursor != 1 {
		t.Errorf("After down arrow, cursor = %d, want 1", model.cursor)
	}

	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = newModel.(MainMenuModel)
	if model.cursor != 2 {
		t.Errorf("After second down arrow, cursor = %d, want 2", model.cursor)
	}

	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	model = newModel.(MainMenuModel)
	if model.cursor != 1 {
		t.Errorf("After up arrow, cursor = %d, want 1", model.cursor)
	}

	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	model = newModel.(MainMenuModel)
	if model.cursor != 0 {
		t.Errorf("After second up arrow, cursor = %d, want 0", model.cursor)
	}

	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	model = newModel.(MainMenuModel)
	if model.cursor != 0 {
		t.Errorf("At top, cursor should stay at 0, got %d", model.cursor)
	}
}

func TestMainMenuModel_AtBottom(t *testing.T) {
	model := NewMainMenuModel()
	model.cursor = 2

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = newModel.(MainMenuModel)
	if model.cursor != 2 {
		t.Errorf("At bottom, cursor should stay at 2, got %d", model.cursor)
	}
}

func TestImportModel_CursorMovement(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportModel(store, "./data")

	if model.cursor != 0 {
		t.Errorf("Initial cursor = %d, want 0", model.cursor)
	}

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = newModel.(ImportModel)
	if model.cursor != 1 {
		t.Errorf("After down arrow, cursor = %d, want 1", model.cursor)
	}

	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = newModel.(ImportModel)
	if model.cursor != 2 {
		t.Errorf("After second down arrow, cursor = %d, want 2", model.cursor)
	}

	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	model = newModel.(ImportModel)
	if model.cursor != 1 {
		t.Errorf("After up arrow, cursor = %d, want 1", model.cursor)
	}
}

func TestImportModel_FileSelectionState(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportModel(store, "./data")

	if model.fileListVisible {
		t.Error("Initial fileListVisible should be false")
	}

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = newModel.(ImportModel)
	if !model.fileListVisible {
		t.Error("After selecting 'Select File', fileListVisible should be true")
	}
	if model.fileCursor != 0 {
		t.Errorf("fileCursor = %d, want 0", model.fileCursor)
	}
}

func TestReportModel_Navigation(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewReportModel(store)

	if model.cursor != 0 {
		t.Errorf("Initial cursor = %d, want 0", model.cursor)
	}

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = newModel.(ReportModel)
	if model.cursor != 1 {
		t.Errorf("After down arrow, cursor = %d, want 1", model.cursor)
	}

	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = newModel.(ReportModel)
	if model.cursor != 2 {
		t.Errorf("After second down arrow, cursor = %d, want 2", model.cursor)
	}

	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = newModel.(ReportModel)
	if model.cursor != 3 {
		t.Errorf("After third down arrow, cursor = %d, want 3", model.cursor)
	}

	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = newModel.(ReportModel)
	if model.cursor != 3 {
		t.Errorf("At bottom, cursor should stay at 3, got %d", model.cursor)
	}

	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	model = newModel.(ReportModel)
	if model.cursor != 2 {
		t.Errorf("After up arrow, cursor = %d, want 2", model.cursor)
	}
}

func TestReportModel_MonthInput(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewReportModel(store)

	model.cursor = 0
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = newModel.(ReportModel)
	if model.inputtingFor != 1 {
		t.Errorf("After selecting Month, inputtingFor = %d, want 1", model.inputtingFor)
	}

	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("3")})
	model = newModel.(ReportModel)
	if model.inputMonth != "3" {
		t.Errorf("inputMonth = %q, want '3'", model.inputMonth)
	}

	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("0")})
	model = newModel.(ReportModel)
	if model.inputMonth != "30" {
		t.Errorf("inputMonth = %q, want '30'", model.inputMonth)
	}
}

func TestReportModel_YearInput(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewReportModel(store)

	model.cursor = 1
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = newModel.(ReportModel)
	if model.inputtingFor != 2 {
		t.Errorf("After selecting Year, inputtingFor = %d, want 2", model.inputtingFor)
	}

	model.inputYear = ""
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("2")})
	model = newModel.(ReportModel)
	if model.inputYear != "2" {
		t.Errorf("inputYear = %q, want '2'", model.inputYear)
	}

	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("0")})
	model = newModel.(ReportModel)
	if model.inputYear != "20" {
		t.Errorf("inputYear = %q, want '20'", model.inputYear)
	}

	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("2")})
	model = newModel.(ReportModel)
	if model.inputYear != "202" {
		t.Errorf("inputYear = %q, want '202'", model.inputYear)
	}
}

func TestReportModel_InputPersistence(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewReportModel(store)

	model.inputMonth = "3"
	model.cursor = 1
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = newModel.(ReportModel)
	if model.inputMonth != "3" {
		t.Errorf("inputMonth should persist after switching to Year, got %q", model.inputMonth)
	}

	model.inputYear = "2024"
	model.cursor = 0
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = newModel.(ReportModel)
	if model.inputMonth != "3" {
		t.Errorf("inputMonth should persist after switching to Month, got %q", model.inputMonth)
	}
	if model.inputYear != "2024" {
		t.Errorf("inputYear should persist after switching to Month, got %q", model.inputYear)
	}
}
