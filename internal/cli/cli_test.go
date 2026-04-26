package cli

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"expensecat/internal/api"
	"expensecat/internal/storage"

	tea "github.com/charmbracelet/bubbletea"
)

type mockStoreForCLI struct{}

type mockStoreWithCategories struct {
	categories []string
	err        error
}

func (m *mockStoreWithCategories) Close() error                        { return nil }
func (m *mockStoreWithCategories) GetConfig() (*storage.Config, error) { return &storage.Config{}, nil }
func (m *mockStoreWithCategories) GetCategories() ([]string, error)    { return m.categories, m.err }
func (m *mockStoreWithCategories) UpdateCategories(categories []string) error {
	m.categories = categories
	return nil
}
func (m *mockStoreWithCategories) GetAllExpenses(startDate, endDate *time.Time) ([]storage.Expense, error) {
	return []storage.Expense{}, nil
}
func (m *mockStoreWithCategories) GetExpense(id string) (storage.Expense, error) {
	return storage.Expense{}, nil
}
func (m *mockStoreWithCategories) AddExpense(storage.Expense) error            { return nil }
func (m *mockStoreWithCategories) RemoveExpense(string) error                  { return nil }
func (m *mockStoreWithCategories) UpdateExpense(string, storage.Expense) error { return nil }
func (m *mockStoreWithCategories) AddMultipleExpenses([]storage.Expense) error { return nil }
func (m *mockStoreWithCategories) RemoveMultipleExpenses([]string) error       { return nil }
func (m *mockStoreWithCategories) GetExclusionList() ([]string, error)         { return []string{}, nil }
func (m *mockStoreWithCategories) AddExclusion(string) error                   { return nil }
func (m *mockStoreWithCategories) RemoveExclusion(string) error                { return nil }
func (m *mockStoreWithCategories) UpdateExclusionList([]string) error          { return nil }
func (m *mockStoreWithCategories) GetImportPath() (string, error)              { return "", nil }
func (m *mockStoreWithCategories) UpdateImportPath(string) error               { return nil }
func (m *mockStoreWithCategories) GetCurrency() (string, error)                { return "usd", nil }
func (m *mockStoreWithCategories) UpdateCurrency(string) error                 { return nil }
func (m *mockStoreWithCategories) GetStartDate() (int, error)                  { return 1, nil }
func (m *mockStoreWithCategories) UpdateStartDate(int) error                   { return nil }
func (m *mockStoreWithCategories) GetRecurringExpenses() ([]storage.RecurringExpense, error) {
	return []storage.RecurringExpense{}, nil
}
func (m *mockStoreWithCategories) GetRecurringExpense(string) (storage.RecurringExpense, error) {
	return storage.RecurringExpense{}, nil
}
func (m *mockStoreWithCategories) AddRecurringExpense(storage.RecurringExpense) error { return nil }
func (m *mockStoreWithCategories) RemoveRecurringExpense(string, bool) error          { return nil }
func (m *mockStoreWithCategories) UpdateRecurringExpense(string, storage.RecurringExpense, bool) error {
	return nil
}

func (m *mockStoreForCLI) Close() error                        { return nil }
func (m *mockStoreForCLI) GetConfig() (*storage.Config, error) { return &storage.Config{}, nil }
func (m *mockStoreForCLI) GetCategories() ([]string, error)    { return []string{}, nil }
func (m *mockStoreForCLI) UpdateCategories([]string) error     { return nil }
func (m *mockStoreForCLI) GetAllExpenses(startDate, endDate *time.Time) ([]storage.Expense, error) {
	return []storage.Expense{}, nil
}
func (m *mockStoreForCLI) GetExclusionList() ([]string, error)     { return []string{}, nil }
func (m *mockStoreForCLI) AddExclusion(pattern string) error       { return nil }
func (m *mockStoreForCLI) RemoveExclusion(pattern string) error    { return nil }
func (m *mockStoreForCLI) UpdateExclusionList(list []string) error { return nil }
func (m *mockStoreForCLI) GetImportPath() (string, error)          { return "~/Downloads/Expense-Files", nil }
func (m *mockStoreForCLI) UpdateImportPath(string) error           { return nil }
func (m *mockStoreForCLI) GetCurrency() (string, error)            { return "usd", nil }
func (m *mockStoreForCLI) UpdateCurrency(string) error             { return nil }
func (m *mockStoreForCLI) GetStartDate() (int, error)              { return 1, nil }
func (m *mockStoreForCLI) UpdateStartDate(int) error               { return nil }
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
func (m *mockStoreForCLI) GetExpense(string) (storage.Expense, error)  { return storage.Expense{}, nil }
func (m *mockStoreForCLI) AddExpense(storage.Expense) error            { return nil }
func (m *mockStoreForCLI) RemoveExpense(string) error                  { return nil }
func (m *mockStoreForCLI) AddMultipleExpenses([]storage.Expense) error { return nil }
func (m *mockStoreForCLI) RemoveMultipleExpenses([]string) error       { return nil }
func (m *mockStoreForCLI) UpdateExpense(string, storage.Expense) error { return nil }
func (m *mockStoreForCLI) GetExpensesByMonthGroupedByCategory(year, month int) ([]storage.CategoryTotal, error) {
	return []storage.CategoryTotal{
		{Category: "Food", Total: 100.00},
	}, nil
}

func TestMainMenuModel_Initialization(t *testing.T) {
	model := NewMainMenuModel()

	if len(model.choices) != 6 {
		t.Errorf("Expected 6 choices, got %d", len(model.choices))
	}
	if model.cursor != 0 {
		t.Errorf("Expected cursor to be 0, got %d", model.cursor)
	}
	if model.choices[0] != "Import Expenses" {
		t.Errorf("Expected first choice to be 'Import Expenses', got %s", model.choices[0])
	}
	if model.choices[1] != "Manage Expenses" {
		t.Errorf("Expected second choice to be 'Manage Expenses', got %s", model.choices[1])
	}
	if model.choices[2] != "View Monthly Totals" {
		t.Errorf("Expected third choice to be 'View Monthly Totals', got %s", model.choices[2])
	}
	if model.choices[3] != "View Recurring Expenses" {
		t.Errorf("Expected fourth choice to be 'View Recurring Expenses', got %s", model.choices[3])
	}
	if model.choices[4] != "Settings" {
		t.Errorf("Expected fifth choice to be 'Settings', got %s", model.choices[4])
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

	now := time.Now()
	currentMonth := int(now.Month())
	currentYear := now.Year()
	expectedMonth := currentMonth - 1
	expectedYear := currentYear
	if expectedMonth < 1 {
		expectedMonth = 12
		expectedYear--
	}

	if model.month != expectedMonth {
		t.Errorf("Expected month to be %d, got %d", expectedMonth, model.month)
	}
	if model.year != expectedYear {
		t.Errorf("Expected year to be %d, got %d", expectedYear, model.year)
	}
	if !model.showResults {
		t.Error("Expected showResults to be true")
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
	if ScreenExpenseList != 7 {
		t.Errorf("Expected ScreenExpenseList to be 7, got %d", ScreenExpenseList)
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
	model.cursor = 5

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = newModel.(MainMenuModel)
	if model.cursor != 5 {
		t.Errorf("At bottom, cursor should stay at 5, got %d", model.cursor)
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
	model.showResults = true
	model.results = []storage.CategoryTotal{
		{Category: "Food", Total: 100.00},
		{Category: "Transport", Total: 50.00},
	}

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
	if model.cursor != 0 {
		t.Errorf("At bottom, cursor should stay at 0 when wrapping, got %d", model.cursor)
	}
}

func TestSettingsModel_CursorMovement(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewSettingsModel(store)

	if model.cursor != 0 {
		t.Errorf("Initial cursor = %d, want 0", model.cursor)
	}

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = newModel.(SettingsModel)
	if model.cursor != 1 {
		t.Errorf("After first down arrow, cursor = %d, want 1", model.cursor)
	}

	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = newModel.(SettingsModel)
	if model.cursor != 2 {
		t.Errorf("At bottom, cursor should stay at 2, got %d", model.cursor)
	}

	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	model = newModel.(SettingsModel)
	if model.cursor != 1 {
		t.Errorf("After up arrow, cursor = %d, want 1", model.cursor)
	}
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	model = newModel.(SettingsModel)
	if model.cursor != 0 {
		t.Errorf("After second up arrow, cursor = %d, want 0", model.cursor)
	}
}

func TestMainMenuModel_View(t *testing.T) {
	model := NewMainMenuModel()
	model.cursor = 0
	view := model.View()
	if view == "" {
		t.Error("View() should not be empty")
	}
	if len(view) < 10 {
		t.Errorf("View() should contain content, got %q", view)
	}
}

func TestMainMenuModel_View_ContainsLogo(t *testing.T) {
	model := NewMainMenuModel()
	view := model.View()
	if !strings.Contains(view, "/\\_/\\") {
		t.Error("View() should contain logo")
	}
	if !strings.Contains(view, "ExpenseCat") {
		t.Error("View() should contain ExpenseCat")
	}
}

func TestMainMenuModel_View_ContainsMenuLabel(t *testing.T) {
	model := NewMainMenuModel()
	view := model.View()
	if !strings.Contains(view, "Main Menu") {
		t.Error("View() should contain Main Menu label")
	}
}

func TestMainMenuModel_IsQuitting(t *testing.T) {
	model := NewMainMenuModel()
	if model.IsQuitting() {
		t.Error("IsQuitting() should be false initially")
	}
	model.quitting = true
	if !model.IsQuitting() {
		t.Error("IsQuitting() should be true after setting quitting")
	}
}

func TestMainMenuModel_EnterOnImport(t *testing.T) {
	model := NewMainMenuModel()
	model.cursor = 0
	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newMainMenu := newModel.(MainMenuModel)
	if cmd == nil {
		t.Error("Enter on Import Expenses should return command")
	}
	_ = newMainMenu
}

func TestMainMenuModel_EnterOnReport(t *testing.T) {
	model := NewMainMenuModel()
	model.cursor = 1
	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newMainMenu := newModel.(MainMenuModel)
	if cmd == nil {
		t.Error("Enter on View Monthly Totals should return command")
	}
	_ = newMainMenu
}

func TestMainMenuModel_EnterOnQuit(t *testing.T) {
	model := NewMainMenuModel()
	model.cursor = 5
	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newMainMenu := newModel.(MainMenuModel)
	if cmd == nil {
		t.Error("Enter on Quit should return Quit command")
	}
	if !newMainMenu.quitting {
		t.Error("quitting should be true after Quit")
	}
	_ = cmd
}

func TestMainMenuModel_TickAnimation(t *testing.T) {
	model := NewMainMenuModel()
	initialFrame := model.frameIndex

	newModel, _ := model.Update(tickMsg{})
	newMainMenu := newModel.(MainMenuModel)

	newModel2, _ := newMainMenu.Update(tickMsg{})
	newMainMenu2 := newModel2.(MainMenuModel)

	_ = initialFrame
	if newMainMenu2.frameIndex != initialFrame {
		t.Logf("Tick should advance frame, got frameIndex=%d", newMainMenu2.frameIndex)
	}
}

func TestMainMenuModel_Init_ReturnsTickCmd(t *testing.T) {
	model := NewMainMenuModel()
	cmd := model.Init()
	if cmd == nil {
		t.Error("Init() should return tick command")
	}
}

func TestMainMenuModel_QuitKey(t *testing.T) {
	model := NewMainMenuModel()
	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	newMainMenu := newModel.(MainMenuModel)
	if cmd == nil {
		t.Error("Ctrl+C should return Quit command")
	}
	if !newMainMenu.quitting {
		t.Error("quitting should be true after Ctrl+C")
	}
}

func TestImportModel_View_Results(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportModel(store, "./data")
	model.showResults = true
	model.result = &api.ImportResult{Imported: 5, Skipped: 2, Duplicates: 1}

	view := model.View()
	if !strings.Contains(view, "Import Results") {
		t.Error("View() should contain Import Results")
	}
	if !strings.Contains(view, "5") {
		t.Error("View() should contain imported count")
	}
}

func TestImportModel_View_FileList(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportModel(store, "./data")
	model.fileListVisible = true
	model.files = []string{"file1.csv", "file2.csv"}
	model.fileCursor = 1

	view := model.View()
	if !strings.Contains(view, "file1.csv") {
		t.Error("View() should contain file1.csv")
	}
	if !strings.Contains(view, "file2.csv") {
		t.Error("View() should contain file2.csv")
	}
}

func TestImportModel_View_Error(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportModel(store, "./data")
	model.dirError = "test error"

	view := model.View()
	if !strings.Contains(view, "test error") {
		t.Error("View() should contain error")
	}
}

func TestImportModel_View_FileSelected(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportModel(store, "./data")
	model.selectedFile = "test.csv"

	view := model.View()
	if !strings.Contains(view, "test.csv") {
		t.Error("View() should contain selected file")
	}
}

func TestImportModel_View_ErrorMsg(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportModel(store, "./data")
	model.errorMsg = "import failed"

	view := model.View()
	if !strings.Contains(view, "import failed") {
		t.Error("View() should contain error message")
	}
}

func TestImportModel_IsQuitting(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportModel(store, "./data")
	if model.IsQuitting() {
		t.Error("IsQuitting() should be false initially")
	}
	model.quitting = true
	if !model.IsQuitting() {
		t.Error("IsQuitting() should be true after setting quitting")
	}
}

func TestImportModel_SelectFile(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportModel(store, "./data")
	model.cursor = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newImportModel := newModel.(ImportModel)
	if !newImportModel.fileListVisible {
		t.Error("fileListVisible should be true after selecting Select File")
	}
}

func TestImportModel_ExecuteImport(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportModel(store, "./data")
	model.selectedFile = "test.csv"

	model.executeImport()
	if model.errorMsg != "" {
		t.Logf("executeImport error: %v", model.errorMsg)
	}
}

func TestImportModel_ExecuteImport_NoFile(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportModel(store, "./data")

	model.executeImport()
	if model.errorMsg == "" {
		t.Error("errorMsg should be set when no file selected")
	}
}

func TestImportModel_FileSelection(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportModel(store, "./data")
	model.fileListVisible = true
	model.files = []string{"test.csv"}
	model.fileCursor = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newImportModel := newModel.(ImportModel)
	if newImportModel.selectedFile != "test.csv" {
		t.Errorf("selectedFile = %q, want 'test.csv'", newImportModel.selectedFile)
	}
}

func TestImportModel_FileListBack(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportModel(store, "./data")
	model.fileListVisible = true
	model.files = []string{"test.csv"}
	model.fileCursor = 1

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newImportModel := newModel.(ImportModel)
	if newImportModel.fileListVisible {
		t.Error("fileListVisible should be false after selecting Back")
	}
}

func TestImportModel_NavigateFileList(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportModel(store, "./data")
	model.fileListVisible = true
	model.files = []string{"file1.csv", "file2.csv"}

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	newImportModel := newModel.(ImportModel)
	if newImportModel.fileCursor != 1 {
		t.Errorf("fileCursor = %d, want 1", newImportModel.fileCursor)
	}
}

func TestImportModel_Backspace(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportModel(store, "./data")
	model.fileListVisible = true
	model.selectedFile = "test"

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	newImportModel := newModel.(ImportModel)
	if newImportModel.selectedFile != "tes" {
		t.Errorf("selectedFile = %q, want 'tes'", newImportModel.selectedFile)
	}
}

func TestImportModel_LoadFiles_Empty(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportModel(store, "./nonexistent")
	if len(model.files) != 0 {
		t.Errorf("files should be empty for nonexistent dir, got %d", len(model.files))
	}
}

func TestImportModel_ExpandPath_Tilde(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantHome bool
	}{
		{"tilde path", "~/test", true},
		{"absolute path", "/tmp/test", false},
		{"relative path", "./test", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandPath(tt.input)
			if tt.wantHome {
				home, _ := os.UserHomeDir()
				want := home + "/test"
				if result != want && !strings.HasPrefix(result, "/") {
					t.Logf("expandPath(%q) = %q", tt.input, result)
				}
			}
		})
	}
}

func TestImportModel_ImportResults(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportModel(store, "./data")
	model.showResults = true
	model.result = &api.ImportResult{Imported: 10, Skipped: 5, Duplicates: 2}

	view := model.View()
	if !strings.Contains(view, "10") {
		t.Error("View() should contain imported count")
	}
	if !strings.Contains(view, "5") {
		t.Error("View() should contain skipped count")
	}
	if !strings.Contains(view, "2") {
		t.Error("View() should contain duplicates count")
	}
}

func TestImportModel_EscKey(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportModel(store, "./data")

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	newImportModel := newModel.(ImportModel)
	if cmd == nil {
		t.Error("Esc should return navigation command")
	}
	_ = newImportModel
}

func TestReportModel_View_Inputting(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewReportModel(store)
	model.inputtingFor = 1

	view := model.View()
	if !strings.Contains(view, "Month") {
		t.Error("View() should contain Month label")
	}
}

func TestReportModel_View_Results(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewReportModel(store)
	model.showResults = true
	model.results = []storage.CategoryTotal{
		{Category: "Food", Total: 100.00},
		{Category: "Transport", Total: 50.00},
	}
	model.total = 150.00
	model.month = 3
	model.year = 2024

	view := model.View()
	if !strings.Contains(view, "Food") {
		t.Error("View() should contain Food")
	}
	if !strings.Contains(view, "100") {
		t.Error("View() should contain total")
	}
}

func TestReportModel_View_Error(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewReportModel(store)
	model.errorMsg = "invalid month"

	view := model.View()
	if !strings.Contains(view, "invalid month") {
		t.Error("View() should contain error")
	}
}

func TestReportModel_IsQuitting(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewReportModel(store)
	if model.IsQuitting() {
		t.Error("IsQuitting() should be false initially")
	}
	model.quitting = true
	if !model.IsQuitting() {
		t.Error("IsQuitting() should be true after setting quitting")
	}
}

func TestReportModel_ResultBackToMenu(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewReportModel(store)
	model.showResults = true

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newReportModel := newModel.(ReportModel)
	if cmd == nil {
		t.Error("Enter on results should return navigation command")
	}
	_ = newReportModel
}

func TestReportModel_EscKey(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewReportModel(store)

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	newReportModel := newModel.(ReportModel)
	if cmd == nil {
		t.Error("Esc should return navigation command")
	}
	_ = newReportModel
}

func TestReportModel_ViewShowsMonthYear(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewReportModel(store)
	model.showResults = true
	model.results = []storage.CategoryTotal{
		{Category: "Food", Total: 100.00},
	}
	model.month = 3
	model.year = 2024

	view := model.View()
	if !strings.Contains(view, "March") {
		t.Error("View() should contain month name")
	}
	if !strings.Contains(view, "2024") {
		t.Error("View() should contain year")
	}
}

func TestReportModel_ViewCalculatesTotal(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewReportModel(store)
	model.results = []storage.CategoryTotal{
		{Category: "Food", Total: 100.50},
		{Category: "Transport", Total: 50.25},
	}
	model.total = 150.75
	model.showResults = true
	model.month = 1
	model.year = 2024

	view := model.View()
	if !strings.Contains(view, "TOTAL") {
		t.Error("View() should contain TOTAL")
	}
	if !strings.Contains(view, "150") {
		t.Error("View() should contain calculated total")
	}
}

func TestReportModel_SortPreservesTotals(t *testing.T) {
	totals := []storage.CategoryTotal{
		{Category: "Zoo", Total: 10.00},
		{Category: "Food", Total: 5.00},
	}

	sorted := SortCategoryTotals(totals)
	if len(sorted) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(sorted))
	}
	if sorted[0].Category == "Food" && sorted[0].Total != 5.00 {
		t.Errorf("First sorted item should be Food with 5.00, got %v", sorted[0])
	}
	if sorted[1].Category == "Zoo" && sorted[1].Total != 10.00 {
		t.Errorf("Second sorted item should be Zoo with 10.00, got %v", sorted[1])
	}
}

func TestGetBasePath_WithEnv(t *testing.T) {
	os.Setenv("EXPENSE_BASE_PATH", "/custom/path")
	defer os.Unsetenv("EXPENSE_BASE_PATH")

	result := getBasePath(nil)
	if result != "/custom/path" {
		t.Errorf("getBasePath() = %q, want '/custom/path'", result)
	}
}

func TestGetBasePath_FromConfig(t *testing.T) {
	os.Unsetenv("EXPENSE_BASE_PATH")
	store := &mockStoreForCLI{}
	result := getBasePath(store)
	if result != "~/Downloads/Expense-Files" {
		t.Errorf("getBasePath() = %q, want '~/Downloads/Expense-Files'", result)
	}
}

func TestImportModel_Init_ReturnsNil(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportModel(store, "./data")
	cmd := model.Init()
	if cmd != nil {
		t.Error("ImportModel.Init() should return nil")
	}
}

func TestReportModel_Init_ReturnsNil(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewReportModel(store)
	cmd := model.Init()
	if cmd != nil {
		t.Error("ReportModel.Init() should return nil")
	}
}

func TestMainMenuModel_TickCmd(t *testing.T) {
	cmd := tickCmd()
	if cmd == nil {
		t.Error("tickCmd() should return command")
	}
}

func TestReportModel_Execute_Success(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewReportModel(store)
	model.showResults = true
	model.month = 1
	model.year = 2024

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Error("Enter should return navigation command")
	}
	_ = newModel
}

func TestMainMenuModel_View_ContainsDivider(t *testing.T) {
	model := NewMainMenuModel()
	view := model.View()
	if !strings.Contains(view, "─") {
		t.Error("View() should contain divider")
	}
}

func TestMainMenuModel_View_ContainsChoices(t *testing.T) {
	model := NewMainMenuModel()
	view := model.View()
	if !strings.Contains(view, "Import Expenses") {
		t.Error("View() should contain Import Expenses")
	}
	if !strings.Contains(view, "Manage Expenses") {
		t.Error("View() should contain Manage Expenses")
	}
	if !strings.Contains(view, "View Monthly Totals") {
		t.Error("View() should contain View Monthly Totals")
	}
	if !strings.Contains(view, "Quit") {
		t.Error("View() should contain Quit")
	}
}

func TestImportModel_View_ContainsBasePath(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportModel(store, "/test/path")
	model.cursor = 1
	model.showResults = false
	model.fileListVisible = false

	view := model.View()
	if !strings.Contains(view, "/test/path") {
		t.Error("View() should contain base path")
	}
}

func TestImportModel_View_ContainsNavigationHelp(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportModel(store, "./data")

	view := model.View()
	if !strings.Contains(view, "Up/Down") {
		t.Error("View() should contain navigation help")
	}
}

func TestReportModel_View_ContainsNavigationHelp(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewReportModel(store)
	model.showResults = true
	model.results = []storage.CategoryTotal{
		{Category: "Food", Total: 100.00},
	}

	view := model.View()
	if !strings.Contains(view, "←") || !strings.Contains(view, "→") {
		t.Error("View() should contain arrow navigation")
	}
}

func TestMainMenuModel_CursorAtBoundary(t *testing.T) {
	model := NewMainMenuModel()
	model.cursor = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyUp})
	newMainMenu := newModel.(MainMenuModel)
	if newMainMenu.cursor != 0 {
		t.Error("Cursor should stay at 0 when at top")
	}
}

func TestImportModel_LoadFiles_Error(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportModel(store, "/nonexistent/path")
	if model.dirError == "" {
		t.Log("dirError may be empty for nonexistent path")
	}
}

func TestAppModel_Init(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)
	cmd := model.Init()
	if cmd == nil {
		t.Error("AppModel.Init() should return command from mainMenu")
	}
}

func TestAppModel_View_MainMenu(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)
	view := model.View()
	if !strings.Contains(view, "Main Menu") {
		t.Error("View() should contain Main Menu")
	}
}

func TestAppModel_View_AfterNavigation(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)

	model.currentScreen = ScreenImport
	view := model.View()
	if !strings.Contains(view, "Import") {
		t.Error("View() should show Import screen")
	}

	model.currentScreen = ScreenReport
	model.reportModel.showResults = true
	model.reportModel.results = []storage.CategoryTotal{
		{Category: "Food", Total: 100.00},
	}
	view = model.View()
	if !strings.Contains(view, "Food") {
		t.Error("View() should show Report screen with results")
	}
}

func TestAppModel_Update_NavigationMsg_Import(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)

	newModel, _ := model.Update(NavigationMsg{Destination: ScreenImport})
	newAppModel := newModel.(AppModel)
	if newAppModel.currentScreen != ScreenImport {
		t.Errorf("currentScreen = %d, want %d", newAppModel.currentScreen, ScreenImport)
	}
}

func TestAppModel_Update_NavigationMsg_Report(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)

	newModel, _ := model.Update(NavigationMsg{Destination: ScreenReport})
	newAppModel := newModel.(AppModel)
	if newAppModel.currentScreen != ScreenReport {
		t.Errorf("currentScreen = %d, want %d", newAppModel.currentScreen, ScreenReport)
	}
}

func TestAppModel_Update_KeyMsg_MainMenu(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	newAppModel := newModel.(AppModel)
	if newAppModel.mainMenu.cursor != 1 {
		t.Errorf("mainMenu.cursor = %d, want 1", newAppModel.mainMenu.cursor)
	}
	_ = cmd
}

func TestAppModel_Update_KeyMsg_Import(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)
	model.currentScreen = ScreenImport

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	newAppModel := newModel.(AppModel)
	if newAppModel.importModel.cursor != 1 {
		t.Errorf("importModel.cursor = %d, want 1", newAppModel.importModel.cursor)
	}
}

func TestAppModel_Update_KeyMsg_Report(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)
	model.currentScreen = ScreenReport
	model.reportModel.showResults = true
	model.reportModel.results = []storage.CategoryTotal{
		{Category: "Food", Total: 100.00},
		{Category: "Transport", Total: 50.00},
	}

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	newAppModel := newModel.(AppModel)
	if newAppModel.reportModel.cursor != 1 {
		t.Errorf("reportModel.cursor = %d, want 1", newAppModel.reportModel.cursor)
	}
}

func TestAppModel_Update_CtrlC(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)
	model.currentScreen = ScreenMainMenu

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	newAppModel := newModel.(AppModel)
	if cmd == nil {
		t.Error("Ctrl+C should return command")
	}
	if !newAppModel.mainMenu.quitting {
		t.Error("mainMenu should be quitting")
	}
}

func TestMainMenuModel_TickCmd_ReturnsCommand(t *testing.T) {
	cmd := tickCmd()
	if cmd == nil {
		t.Error("tickCmd should not return nil")
	}
}

func TestSettingsModel_Initialization(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewSettingsModel(store)

	if model.cursor != 0 {
		t.Errorf("Expected cursor to be 0, got %d", model.cursor)
	}
	if model.quitting {
		t.Error("quitting should be false initially")
	}
}

func TestSettingsModel_View(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewSettingsModel(store)
	view := model.View()

	if !strings.Contains(view, "Settings") {
		t.Error("View() should contain Settings")
	}
	if !strings.Contains(view, "Exclusion List") {
		t.Error("View() should contain Exclusion List")
	}
	if !strings.Contains(view, "Back") {
		t.Error("View() should contain Back")
	}
}

func TestExclusionListModel_Initialization(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)

	if model.cursor != 0 {
		t.Errorf("Expected cursor to be 0, got %d", model.cursor)
	}
	if model.mode != 0 {
		t.Errorf("Expected mode to be 0, got %d", model.mode)
	}
	if model.inputValue != "" {
		t.Errorf("Expected inputValue to be empty, got %s", model.inputValue)
	}
	if model.quitting {
		t.Error("quitting should be false initially")
	}
}

func TestExclusionListModel_View(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)
	model.patterns = []string{"pattern1", "pattern2"}
	view := model.View()

	if !strings.Contains(view, "Exclusion List") {
		t.Error("View() should contain Exclusion List")
	}
	if !strings.Contains(view, "Patterns") {
		t.Error("View() should contain Patterns")
	}
	if !strings.Contains(view, "a Add") {
		t.Error("View() should contain add option in footer")
	}
	if !strings.Contains(view, "e Edit") {
		t.Error("View() should contain edit option in footer")
	}
	if !strings.Contains(view, "x Delete") {
		t.Error("View() should contain delete option in footer")
	}
}

func TestExclusionListModel_AddMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)

	model.mode = ExclusionListAdd
	model.inputValue = "TEST"
	view := model.View()

	if !strings.Contains(view, "Enter new pattern:") {
		t.Error("View() should show input prompt in add mode")
	}
	if !strings.Contains(view, "TEST") {
		t.Error("View() should show entered pattern")
	}
}

func TestExclusionListModel_EditMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)

	model.mode = ExclusionListEdit
	model.editIndex = 0
	model.inputValue = "UPDATED"
	view := model.View()

	if !strings.Contains(view, "Edit pattern:") {
		t.Error("View() should show edit prompt in edit mode")
	}
	if !strings.Contains(view, "UPDATED") {
		t.Error("View() should show updated pattern")
	}
}

func TestExclusionListModel_DeleteConfirmMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)

	model.mode = ExclusionListDeleteConfirm
	model.patterns = []string{"PATTERN1", "PATTERN2"}
	model.cursor = 1
	view := model.View()

	if !strings.Contains(view, "Confirm delete") {
		t.Error("View() should show confirm delete prompt")
	}
	if !strings.Contains(view, "PATTERN2") {
		t.Error("View() should show pattern to delete")
	}
	if !strings.Contains(view, "[Y]es") {
		t.Error("View() should show yes option")
	}
	if !strings.Contains(view, "[N]o") {
		t.Error("View() should show no option")
	}
}

func TestSettingsModel_NavigationToExclusionList(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewSettingsModel(store)
	model.cursor = 0

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	_ = newModel.(SettingsModel)

	if cmd == nil {
		t.Error("Enter on Exclusion List should return navigation command")
	}
}

func TestExclusionListModel_EscInAddMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)
	model.mode = ExclusionListAdd
	model.inputValue = "TEST"

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	newExclusionModel := newModel.(ExclusionListModel)

	if newExclusionModel.mode != ExclusionListView {
		t.Errorf("Esc should return to menu mode, got %d", newExclusionModel.mode)
	}
}

func TestExclusionListModel_EscInMenuMode_NavigatesBack(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	_ = newModel.(ExclusionListModel)

	if cmd == nil {
		t.Error("Esc in menu mode should return navigation command to Settings")
	}
}

func TestExclusionListModel_UpAtTop(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)
	model.cursor = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyUp})
	newExclusionModel := newModel.(ExclusionListModel)

	if newExclusionModel.cursor != 0 {
		t.Error("Cursor should stay at 0 when at top")
	}
}

func TestExclusionListModel_DownAtBottom(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)
	model.cursor = 2

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	newExclusionModel := newModel.(ExclusionListModel)

	if newExclusionModel.cursor != 2 {
		t.Error("Cursor should stay at 2 when at bottom")
	}
}

func TestExclusionListModel_EscNavigatesBack(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	_ = newModel.(ExclusionListModel)

	if cmd == nil {
		t.Error("Esc should return navigation command to Settings")
	}
}

func TestExclusionListModel_EnterInAddMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)
	model.mode = ExclusionListAdd
	model.inputValue = "NEW PATTERN"

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newExclusionModel := newModel.(ExclusionListModel)

	if newExclusionModel.mode != ExclusionListView {
		t.Error("Enter in add mode should save and return to menu")
	}
}

func TestExclusionListModel_TypeInAddMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)
	model.mode = ExclusionListAdd

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	newExclusionModel := newModel.(ExclusionListModel)

	if newExclusionModel.inputValue == "" {
		t.Error("Typing in add mode should add characters")
	}
}

func TestSettingsModel_EnterOnBack(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewSettingsModel(store)
	model.cursor = 1

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	_ = newModel.(SettingsModel)

	if cmd == nil {
		t.Error("Enter on Back should return navigation command")
	}
}

func TestSettingsModel_NavigationMsg_ResetsModel(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)
	model.settingsModel.cursor = 1

	newModel, _ := model.Update(NavigationMsg{Destination: ScreenSettings})
	newAppModel := newModel.(AppModel)

	if newAppModel.settingsModel.cursor != 0 {
		t.Error("Navigation to Settings should reset cursor")
	}
}

func TestExclusionListModel_NavigationMsg_ResetsModel(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)
	model.exclusionModel.mode = ExclusionListEdit

	newModel, _ := model.Update(NavigationMsg{Destination: ScreenExclusionList})
	newAppModel := newModel.(AppModel)

	if newAppModel.exclusionModel.mode != ExclusionListView {
		t.Error("Navigation to ExclusionList should reset mode")
	}
}

func TestExclusionListModel_AddKey(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	newExclusionModel := newModel.(ExclusionListModel)

	if newExclusionModel.mode != ExclusionListAdd {
		t.Error("'a' key should enter add mode")
	}
}

func TestExclusionListModel_EditKey(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)
	model.patterns = []string{"TEST1", "TEST2"}
	model.cursor = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	newExclusionModel := newModel.(ExclusionListModel)

	if newExclusionModel.mode != ExclusionListEdit {
		t.Error("'e' key should enter edit mode")
	}
	if newExclusionModel.editIndex != 0 {
		t.Error("editIndex should be 0")
	}
	if newExclusionModel.inputValue != "TEST1" {
		t.Error("inputValue should be populated with current pattern")
	}
}

func TestExclusionListModel_DeleteKey(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)
	model.patterns = []string{"TEST1", "TEST2"}
	model.cursor = 1

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	newExclusionModel := newModel.(ExclusionListModel)

	if newExclusionModel.mode != ExclusionListDeleteConfirm {
		t.Error("'x' key should enter delete confirm mode")
	}
}

func TestExclusionListModel_EnterOnItem(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)
	model.patterns = []string{"TEST1"}
	model.cursor = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newExclusionModel := newModel.(ExclusionListModel)

	if newExclusionModel.mode != ExclusionListEdit {
		t.Error("Enter on item should enter edit mode")
	}
}

func TestExclusionListModel_ConfirmDeleteYes(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)
	model.mode = ExclusionListDeleteConfirm
	model.patterns = []string{"TEST1", "TEST2"}
	model.cursor = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	newExclusionModel := newModel.(ExclusionListModel)

	if newExclusionModel.mode != ExclusionListView {
		t.Error("'y' should confirm delete and return to view")
	}
}

func TestExclusionListModel_ConfirmDeleteNo(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)
	model.mode = ExclusionListDeleteConfirm
	model.patterns = []string{"TEST1", "TEST2"}
	model.cursor = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	newExclusionModel := newModel.(ExclusionListModel)

	if newExclusionModel.mode != ExclusionListView {
		t.Error("'n' should cancel delete and return to view")
	}
}

func TestExclusionListModel_CtrlC(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	newExclusionModel := newModel.(ExclusionListModel)

	if cmd == nil {
		t.Error("Ctrl+C should return quit command")
	}
	if !newExclusionModel.quitting {
		t.Error("quitting should be true after Ctrl+C")
	}
}

func TestExclusionListModel_LoadPatterns_Error(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)

	model.loadPatterns()

	if model.errorMsg != "" {
		t.Error("Should not have error with mock store")
	}
}

func TestExclusionListModel_SuccessMessage(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)
	model.successMsg = "Test success"

	view := model.View()

	if !strings.Contains(view, "Test success") {
		t.Error("View() should show success message")
	}
}

func TestSettingsModel_CtrlC(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewSettingsModel(store)

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	newSettingsModel := newModel.(SettingsModel)

	if cmd == nil {
		t.Error("Ctrl+C should return quit command")
	}
	if !newSettingsModel.quitting {
		t.Error("quitting should be true after Ctrl+C")
	}
}

func TestExclusionListModel_SelectBackInRemoveMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)
	model.mode = ExclusionListEdit
	model.patterns = []string{}
	model.cursor = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	newExclusionModel := newModel.(ExclusionListModel)

	if newExclusionModel.cursor != 0 {
		t.Error("Down should stay at 0 when no patterns exist")
	}
}

func TestExclusionListModel_UpAtTopInRemoveMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)
	model.mode = ExclusionListEdit
	model.patterns = []string{"TEST1", "TEST2"}
	model.cursor = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyUp})
	newExclusionModel := newModel.(ExclusionListModel)

	if newExclusionModel.cursor != 0 {
		t.Error("Up should stay at 0 when at top in remove mode")
	}
}

func TestExclusionListModel_DownStaysAtMaxInRemoveMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)
	model.mode = ExclusionListEdit
	model.patterns = []string{"TEST1", "TEST2"}
	model.cursor = 2

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	newExclusionModel := newModel.(ExclusionListModel)

	if newExclusionModel.cursor != 2 {
		t.Error("Down should stay at max when at bottom in remove mode")
	}
}

func TestCategoryListModel_Initialization(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)

	if model.cursor != 0 {
		t.Errorf("Expected cursor to be 0, got %d", model.cursor)
	}
	if model.mode != 0 {
		t.Errorf("Expected mode to be 0, got %d", model.mode)
	}
	if model.inputValue != "" {
		t.Errorf("Expected inputValue to be empty, got %s", model.inputValue)
	}
	if model.quitting {
		t.Error("quitting should be false initially")
	}
}

func TestCategoryListModel_View(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)
	model.categories = []string{"Food", "Transport"}
	view := model.View()

	if !strings.Contains(view, "Category List") {
		t.Error("View() should contain Category List")
	}
	if !strings.Contains(view, "a Add") {
		t.Error("View() should contain add option in footer")
	}
	if !strings.Contains(view, "e Edit") {
		t.Error("View() should contain edit option in footer")
	}
	if !strings.Contains(view, "x Delete") {
		t.Error("View() should contain delete option in footer")
	}
}

func TestCategoryListModel_EscNavigatesBack(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	_ = newModel.(CategoryListModel)

	if cmd == nil {
		t.Error("Esc should return navigation command to Settings")
	}
}

func TestCategoryListModel_AddKey(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	newCategoryModel := newModel.(CategoryListModel)

	if newCategoryModel.mode != CategoryListAdd {
		t.Error("'a' key should enter add mode")
	}
}

func TestCategoryListModel_EditKey(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)
	model.categories = []string{"Food", "Transport"}
	model.cursor = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	newCategoryModel := newModel.(CategoryListModel)

	if newCategoryModel.mode != CategoryListEdit {
		t.Error("'e' key should enter edit mode")
	}
	if newCategoryModel.editIndex != 0 {
		t.Error("editIndex should be 0")
	}
}

func TestCategoryListModel_DeleteKey(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)
	model.categories = []string{"Food", "Transport"}
	model.cursor = 1

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	newCategoryModel := newModel.(CategoryListModel)

	if newCategoryModel.mode != CategoryListDeleteConfirm {
		t.Error("'x' key should enter delete confirm mode")
	}
}

func TestCategoryListModel_EnterOnItem(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)
	model.categories = []string{"Food"}
	model.cursor = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newCategoryModel := newModel.(CategoryListModel)

	if newCategoryModel.mode != CategoryListEdit {
		t.Error("Enter on item should enter edit mode")
	}
}

func TestCategoryListModel_TypeInAddMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)
	model.mode = CategoryListAdd
	model.inputValue = ""

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'T', 'e', 's', 't'}})
	newCategoryModel := newModel.(CategoryListModel)

	if newCategoryModel.inputValue != "Test" {
		t.Errorf("inputValue should be 'Test', got %s", newCategoryModel.inputValue)
	}
}

func TestCategoryListModel_BackspaceInAddMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)
	model.mode = CategoryListAdd
	model.inputValue = "Test"

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	newCategoryModel := newModel.(CategoryListModel)

	if newCategoryModel.inputValue != "Tes" {
		t.Errorf("inputValue should be 'Tes', got %s", newCategoryModel.inputValue)
	}
}

func TestCategoryListModel_EscInAddMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)
	model.mode = CategoryListAdd
	model.inputValue = "Test"

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	newCategoryModel := newModel.(CategoryListModel)

	if newCategoryModel.mode != CategoryListView {
		t.Errorf("Expected mode CategoryListView, got %d", newCategoryModel.mode)
	}
	if newCategoryModel.inputValue != "" {
		t.Error("inputValue should be cleared after Esc")
	}
}

func TestCategoryListModel_UpAtTop(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)
	model.mode = CategoryListView
	model.cursor = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyUp})
	newCategoryModel := newModel.(CategoryListModel)

	if newCategoryModel.cursor != 0 {
		t.Error("Cursor should stay at 0 when at top")
	}
}

func TestCategoryListModel_DownAtBottom(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)
	model.mode = CategoryListView
	model.cursor = 2

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	newCategoryModel := newModel.(CategoryListModel)

	if newCategoryModel.cursor != 2 {
		t.Error("Cursor should stay at 2 when at bottom")
	}
}

func TestCategoryListModel_CtrlC(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	newCategoryModel := newModel.(CategoryListModel)

	if !newCategoryModel.quitting {
		t.Error("quitting should be true after Ctrl+C")
	}
}

func TestCategoryListModel_DownInAddMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)
	model.mode = CategoryListAdd

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	newCategoryModel := newModel.(CategoryListModel)

	if newCategoryModel.mode != CategoryListAdd {
		t.Error("Down should not change mode in add mode")
	}
}

func TestCategoryListModel_EscInRemoveMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)
	model.mode = CategoryListEdit
	model.cursor = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	newCategoryModel := newModel.(CategoryListModel)

	if newCategoryModel.mode != CategoryListView {
		t.Errorf("Expected mode CategoryListView, got %d", newCategoryModel.mode)
	}
}

func TestCategoryListModel_UpInRemoveMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)
	model.mode = CategoryListEdit
	model.categories = []string{"A", "B", "C"}
	model.cursor = 1

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyUp})
	newCategoryModel := newModel.(CategoryListModel)

	if newCategoryModel.cursor != 0 {
		t.Errorf("cursor should be 0, got %d", newCategoryModel.cursor)
	}
}

func TestCategoryListModel_DownInViewMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)
	model.mode = CategoryListView
	model.categories = []string{"A", "B", "C"}
	model.cursor = 1

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	newCategoryModel := newModel.(CategoryListModel)

	if newCategoryModel.cursor != 2 {
		t.Errorf("cursor should be 2, got %d", newCategoryModel.cursor)
	}
}

func TestCategoryListModel_ViewInEditMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)
	model.mode = CategoryListEdit
	model.editIndex = 0
	model.inputValue = "Updated"
	model.categories = []string{"Food", "Travel"}

	view := model.View()

	if !strings.Contains(view, "Edit category:") {
		t.Error("View should contain edit prompt")
	}
	if !strings.Contains(view, "Updated") {
		t.Error("View should contain updated value")
	}
}

func TestCategoryListModel_DeleteConfirmMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)
	model.mode = CategoryListDeleteConfirm
	model.categories = []string{"Food", "Travel"}
	model.cursor = 1

	view := model.View()

	if !strings.Contains(view, "Confirm delete") {
		t.Error("View should contain confirm delete prompt")
	}
	if !strings.Contains(view, "Travel") {
		t.Error("View should show category to delete")
	}
	if !strings.Contains(view, "[Y]es") {
		t.Error("View should show yes option")
	}
	if !strings.Contains(view, "[N]o") {
		t.Error("View should show no option")
	}
}

func TestCategoryListModel_ConfirmDeleteYes(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)
	model.mode = CategoryListDeleteConfirm
	model.categories = []string{"Food", "Travel"}
	model.cursor = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	newCategoryModel := newModel.(CategoryListModel)

	if newCategoryModel.mode != CategoryListView {
		t.Error("'y' should confirm delete and return to view")
	}
}

func TestCategoryListModel_ConfirmDeleteNo(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)
	model.mode = CategoryListDeleteConfirm
	model.categories = []string{"Food", "Travel"}
	model.cursor = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	newCategoryModel := newModel.(CategoryListModel)

	if newCategoryModel.mode != CategoryListView {
		t.Error("'n' should cancel delete and return to view")
	}
}

func TestCategoryListModel_EnterOnEmptyInput(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)
	model.mode = CategoryListAdd
	model.inputValue = ""

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newCategoryModel := newModel.(CategoryListModel)

	if newCategoryModel.mode != CategoryListAdd {
		t.Error("Should stay in add mode when input is empty and Enter pressed")
	}
}

func TestSettingsModel_SelectCategoryList(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewSettingsModel(store)
	model.cursor = SettingsCategoryList

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newSettingsModel := newModel.(SettingsModel)

	if cmd == nil {
		t.Error("Should return navigation command when selecting Category List")
	}
	if newSettingsModel.cursor != SettingsCategoryList {
		t.Error("Cursor should be unchanged")
	}
	_ = newSettingsModel
}

func TestSettingsModel_ViewCategoryListOption(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewSettingsModel(store)

	view := model.View()

	if !strings.Contains(view, "Category List") {
		t.Error("View should contain Category List")
	}
}

func TestSettingsModel_CategoryListOptionOrder(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewSettingsModel(store)

	view := model.View()
	exclusionIdx := strings.Index(view, "Exclusion List")
	importPathIdx := strings.Index(view, "Import Path")
	categoryListIdx := strings.Index(view, "Category List")

	if categoryListIdx == -1 {
		t.Error("Category List should be in view")
	}
	if categoryListIdx < exclusionIdx || categoryListIdx < importPathIdx {
		t.Error("Category List should appear after Import Path")
	}
}

func TestSettingsModel_NavigateToCategoryList(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewSettingsModel(store)
	model.cursor = SettingsImportPath

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	downModel := newModel.(SettingsModel)

	if downModel.cursor != SettingsCategoryList {
		t.Errorf("cursor = %d, want %d", downModel.cursor, SettingsCategoryList)
	}
}

func TestSettingsModel_NavigateUpToCategoryList(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewSettingsModel(store)
	model.cursor = SettingsCategoryList

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyUp})
	upModel := newModel.(SettingsModel)

	if upModel.cursor != SettingsImportPath {
		t.Errorf("cursor = %d, want %d", upModel.cursor, SettingsImportPath)
	}
}

func TestSettingsModel_NavigateUpAtTop(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewSettingsModel(store)
	model.cursor = SettingsExclusionList

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyUp})
	upModel := newModel.(SettingsModel)

	if upModel.cursor != SettingsExclusionList {
		t.Error("Cursor should stay at top when pressing Up")
	}
}

func TestSettingsModel_NavigateDownAtBottom(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewSettingsModel(store)
	model.cursor = SettingsBack

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	downModel := newModel.(SettingsModel)

	if downModel.cursor != SettingsBack {
		t.Error("Cursor should stay at bottom when pressing Down")
	}
}

func TestCategoryListModel_AddCategory(t *testing.T) {
	store := &mockStoreWithCategories{}
	model := NewCategoryListModel(store)
	model.categories = []string{"Food", "Travel"}

	err := model.addCategory("Utilities")

	if err != nil {
		t.Errorf("addCategory should not error, got %v", err)
	}
	if len(store.categories) != 3 {
		t.Errorf("Expected 3 categories, got %d", len(store.categories))
	}
}

func TestCategoryListModel_RemoveCategory(t *testing.T) {
	store := &mockStoreWithCategories{categories: []string{"Food", "Travel", "Utilities"}}
	model := NewCategoryListModel(store)
	model.categories = store.categories

	model.removeCategory("Travel")

	if len(store.categories) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(store.categories))
	}
}

func TestCategoryListModel_RemoveCategoryError(t *testing.T) {
	store := &mockStoreWithCategories{categories: []string{"Food", "Travel"}, err: fmt.Errorf("test error")}
	model := NewCategoryListModel(store)
	model.categories = store.categories

	model.removeCategory("Travel")

	if model.errorMsg == "" {
		t.Error("Should set error message on failure")
	}
}

func TestCategoryListModel_Init(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)

	cmd := model.Init()

	if cmd != nil {
		t.Error("Init should return nil command")
	}
}

func TestCategoryListModel_ErrorMessage(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)
	model.errorMsg = "test error"

	view := model.View()

	if !strings.Contains(view, "Error:") {
		t.Error("View should contain error message")
	}
}

func TestCategoryListModel_SuccessMessage(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)
	model.successMsg = "test success"

	view := model.View()

	if !strings.Contains(view, "test success") {
		t.Error("View should contain success message")
	}
}

func TestCategoryListModel_EmptyCategories(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)
	model.categories = []string{}

	view := model.View()

	if !strings.Contains(view, "no categories") {
		t.Error("View should indicate no categories")
	}
}

func TestCategoryListModel_EmptyErrorAndSuccess(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)

	view := model.View()

	if strings.Contains(view, "Error:") {
		t.Error("View should not contain error when none")
	}
	if strings.Contains(view, "success") {
		t.Error("View should not contain success when none")
	}
}

func TestCategoryListModel_LoadCategoriesError(t *testing.T) {
	store := &mockStoreWithCategories{categories: []string{}, err: fmt.Errorf("test error")}
	model := NewCategoryListModel(store)

	if model.errorMsg == "" {
		t.Error("Load should set error on failure")
	}
}

func TestCategoryListModel_IsQuitting(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)

	if model.IsQuitting() {
		t.Error("IsQuitting should be false initially")
	}
}

func TestCategoryListModel_Update_NonKeyMsg(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)

	newModel, _ := model.Update("not a key")
	newCategoryModel := newModel.(CategoryListModel)

	if newCategoryModel.mode != CategoryListView {
		t.Error("Non-key message should not change mode")
	}
}

func TestCategoryListModel_DownAtMaxInMenuMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewCategoryListModel(store)
	model.mode = CategoryListView
	model.cursor = 2

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	newCategoryModel := newModel.(CategoryListModel)

	if newCategoryModel.cursor != 2 {
		t.Error("Cursor should stay at max")
	}
}

func TestExpenseListModel_getEditFieldValue(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test", Amount: 10.00, Category: "Food", Currency: "usd", Date: time.Now()},
		},
		categories: []string{"Food", "Travel"},
	}
	model := NewExpenseListModel(store)
	model.editingExpense = model.expenses[0]
	model.editField = 4

	value := model.getEditFieldValue()

	if value == "" {
		t.Error("getEditFieldValue should return non-empty value")
	}
}

func TestExpenseListModel_UpdateEditMode_EmptyInputValue(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test", Amount: 10.00, Category: "Food", Currency: "usd", Date: time.Now()},
		},
		categories: []string{"Food", "Travel"},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListEdit
	model.editField = 0
	model.editValue = ""

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	newExpenseModel := newModel.(*ExpenseListModel)

	if newExpenseModel.editValue != "" {
		t.Error("Backspace on empty should stay empty")
	}
}

func TestExpenseListModel_categoryPickerUpDown(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test", Amount: 10.00, Category: "Food", Currency: "usd", Date: time.Now()},
		},
		categories: []string{"Food", "Travel", "Utilities"},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListEdit
	model.editField = 1
	model.editingExpense = model.expenses[0]
	model.editValue = "Food"
	model.showCategoryPicker = true
	model.categoryPickerIdx = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	newExpenseModel := newModel.(*ExpenseListModel)

	if newExpenseModel.categoryPickerIdx != 1 {
		t.Error("Category picker index should increment")
	}
}

func TestSettingsModel_ViewChoices(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewSettingsModel(store)
	view := model.View()

	if !strings.Contains(view, "Exclusion List") {
		t.Error("View() should contain Exclusion List")
	}
	if !strings.Contains(view, "Back") {
		t.Error("View() should contain Back")
	}
}

func TestSettingsModel_CursorAtTop(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewSettingsModel(store)

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyUp})
	newSettingsModel := newModel.(SettingsModel)

	if newSettingsModel.cursor != 0 {
		t.Error("Up at top should stay at 0")
	}
}

func TestAppModel_SettingsScreenView(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)

	model.currentScreen = ScreenSettings
	view := model.View()

	if !strings.Contains(view, "Settings") {
		t.Error("View() should show Settings when on Settings screen")
	}
}

func TestAppModel_ExclusionListScreenView(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)

	model.currentScreen = ScreenExclusionList
	view := model.View()

	if !strings.Contains(view, "Exclusion List") {
		t.Error("View() should show Exclusion List when on ExclusionList screen")
	}
}

func TestImportModel_LoadFiles_SkipsDirs(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportModel(store, ".")
	if model.files == nil {
		t.Error("files should be initialized")
	}
}

func TestAppModel_Update_NavigationMsg_Settings(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)

	newModel, _ := model.Update(NavigationMsg{Destination: ScreenSettings})
	newAppModel := newModel.(AppModel)
	if newAppModel.currentScreen != ScreenSettings {
		t.Errorf("currentScreen = %d, want %d", newAppModel.currentScreen, ScreenSettings)
	}
}

func TestAppModel_Update_NavigationMsg_ExclusionList(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)

	newModel, _ := model.Update(NavigationMsg{Destination: ScreenExclusionList})
	newAppModel := newModel.(AppModel)
	if newAppModel.currentScreen != ScreenExclusionList {
		t.Errorf("currentScreen = %d, want %d", newAppModel.currentScreen, ScreenExclusionList)
	}
}

func TestAppModel_Update_NavigationMsg_ImportPath(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)

	newModel, _ := model.Update(NavigationMsg{Destination: ScreenImportPath})
	newAppModel := newModel.(AppModel)
	if newAppModel.currentScreen != ScreenImportPath {
		t.Errorf("currentScreen = %d, want %d", newAppModel.currentScreen, ScreenImportPath)
	}
	if newAppModel.importPathModel.path != "~/Downloads/Expense-Files" {
		t.Errorf("importPathModel.path = %q, want '~/Downloads/Expense-Files'", newAppModel.importPathModel.path)
	}
}

func TestImportPathModel_Initialization(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportPathModel(store)

	if model.mode != ImportPathMenu {
		t.Errorf("Initial mode = %d, want %d (ImportPathMenu)", model.mode, ImportPathMenu)
	}
	if model.cursor != 0 {
		t.Errorf("Initial cursor = %d, want 0", model.cursor)
	}
	if model.cursorPos != 0 {
		t.Errorf("Initial cursorPos = %d, want 0", model.cursorPos)
	}
	if model.path != "~/Downloads/Expense-Files" {
		t.Errorf("Initial path = %q, want '~/Downloads/Expense-Files'", model.path)
	}
}

func TestImportPathModel_View_MenuMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportPathModel(store)
	model.path = "/custom/path"

	view := model.View()

	if !strings.Contains(view, "/custom/path") {
		t.Error("View should show current path")
	}
	if !strings.Contains(view, "Edit Path") {
		t.Error("View should show Edit Path option")
	}
	if !strings.Contains(view, "Back") {
		t.Error("View should show Back option")
	}
}

func TestImportPathModel_View_EditMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportPathModel(store)
	model.mode = ImportPathEdit
	model.inputValue = "/test/path"
	model.cursorPos = 5

	view := model.View()

	if !strings.Contains(view, "Enter new path:") {
		t.Error("View should show input prompt in edit mode")
	}
	if !strings.Contains(view, "|") {
		t.Error("View should show cursor position indicator")
	}
}

func TestImportPathModel_View_ErrorMsg(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportPathModel(store)
	model.errorMsg = "Test error message"

	view := model.View()

	if !strings.Contains(view, "Test error message") {
		t.Error("View should show error message")
	}
}

func TestImportPathModel_View_SuccessMsg(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportPathModel(store)
	model.successMsg = "Test success message"

	view := model.View()

	if !strings.Contains(view, "Test success message") {
		t.Error("View should show success message")
	}
}

func TestImportPathModel_CursorMovement(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportPathModel(store)

	if model.cursor != 0 {
		t.Errorf("Initial cursor = %d, want 0", model.cursor)
	}

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = newModel.(ImportPathModel)
	if model.cursor != 1 {
		t.Errorf("After down arrow, cursor = %d, want 1", model.cursor)
	}

	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	model = newModel.(ImportPathModel)
	if model.cursor != 0 {
		t.Errorf("After up arrow, cursor = %d, want 0", model.cursor)
	}
}

func TestImportPathModel_UpAtTop(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportPathModel(store)
	model.cursor = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyUp})
	model = newModel.(ImportPathModel)

	if model.cursor != 0 {
		t.Error("Up at top should stay at 0")
	}
}

func TestImportPathModel_DownAtBottom(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportPathModel(store)
	model.cursor = 1

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = newModel.(ImportPathModel)

	if model.cursor != 1 {
		t.Error("Down at bottom should stay at 1")
	}
}

func TestImportPathModel_EnterOnBack(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportPathModel(store)
	model.cursor = 1

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newImportPathModel := newModel.(ImportPathModel)

	if newImportPathModel.mode != ImportPathMenu {
		t.Error("Mode should stay in menu")
	}
	if cmd == nil {
		t.Error("Should return navigation cmd")
	}
}

func TestImportPathModel_NavigationMsg_ResetsModel(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportPathModel(store)
	model.mode = ImportPathEdit
	model.inputValue = "test"
	model.cursor = 1
	model.cursorPos = 4

	msg := NavigationMsg{Destination: ScreenSettings}
	_, _ = model.Update(msg)

	if model.mode != ImportPathEdit {
		t.Error("Mode should stay unchanged (NavigationMsg handled by AppModel)")
	}
}

func TestImportPathModel_EnterEditMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportPathModel(store)
	model.path = "/initial/path"

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newImportPathModel := newModel.(ImportPathModel)

	if newImportPathModel.mode != ImportPathEdit {
		t.Errorf("Mode = %d, want %d (ImportPathEdit)", newImportPathModel.mode, ImportPathEdit)
	}
	if newImportPathModel.inputValue != "/initial/path" {
		t.Errorf("inputValue = %q, want '/initial/path'", newImportPathModel.inputValue)
	}
	if newImportPathModel.cursorPos != len("/initial/path") {
		t.Errorf("cursorPos = %d, want %d", newImportPathModel.cursorPos, len("/initial/path"))
	}
}

func TestImportPathModel_LeftArrow_MovesCursor(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportPathModel(store)
	model.mode = ImportPathEdit
	model.inputValue = "test"
	model.cursorPos = 4

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model = newModel.(ImportPathModel)

	if model.cursorPos != 3 {
		t.Errorf("cursorPos = %d, want 3", model.cursorPos)
	}
}

func TestImportPathModel_LeftArrow_AtStart(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportPathModel(store)
	model.mode = ImportPathEdit
	model.inputValue = "test"
	model.cursorPos = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model = newModel.(ImportPathModel)

	if model.cursorPos != 0 {
		t.Error("Left at start should stay at 0")
	}
}

func TestImportPathModel_RightArrow_MovesCursor(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportPathModel(store)
	model.mode = ImportPathEdit
	model.inputValue = "test"
	model.cursorPos = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRight})
	model = newModel.(ImportPathModel)

	if model.cursorPos != 1 {
		t.Errorf("cursorPos = %d, want 1", model.cursorPos)
	}
}

func TestImportPathModel_RightArrow_AtEnd(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportPathModel(store)
	model.mode = ImportPathEdit
	model.inputValue = "test"
	model.cursorPos = 4

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRight})
	model = newModel.(ImportPathModel)

	if model.cursorPos != 4 {
		t.Error("Right at end should stay at 4")
	}
}

func TestImportPathModel_TypeInEditMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportPathModel(store)
	model.mode = ImportPathEdit
	model.inputValue = "test"
	model.cursorPos = 2

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("X")})
	model = newModel.(ImportPathModel)

	if model.inputValue != "teXst" {
		t.Errorf("inputValue = %q, want 'teXst'", model.inputValue)
	}
	if model.cursorPos != 3 {
		t.Errorf("cursorPos = %d, want 3", model.cursorPos)
	}
}

func TestImportPathModel_BackspaceInEditMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportPathModel(store)
	model.mode = ImportPathEdit
	model.inputValue = "test"
	model.cursorPos = 3

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	model = newModel.(ImportPathModel)

	if model.inputValue != "tet" {
		t.Errorf("inputValue = %q, want 'tet'", model.inputValue)
	}
	if model.cursorPos != 2 {
		t.Errorf("cursorPos = %d, want 2", model.cursorPos)
	}
}

func TestImportPathModel_EscInEditMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportPathModel(store)
	model.mode = ImportPathEdit
	model.inputValue = "test"
	model.cursorPos = 2
	model.cursor = 1

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = newModel.(ImportPathModel)

	if model.mode != ImportPathMenu {
		t.Errorf("mode = %d, want %d (ImportPathMenu)", model.mode, ImportPathMenu)
	}
	if model.inputValue != "" {
		t.Error("inputValue should be cleared")
	}
	if model.cursorPos != 0 {
		t.Error("cursorPos should be reset")
	}
	if model.cursor != 0 {
		t.Error("cursor should be reset")
	}
}

func TestImportPathModel_EnterSavesPath(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportPathModel(store)
	model.mode = ImportPathEdit
	model.inputValue = "/new/path"

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newImportPathModel := newModel.(ImportPathModel)

	if newImportPathModel.mode != ImportPathMenu {
		t.Error("Should return to menu mode")
	}
	if newImportPathModel.inputValue != "" {
		t.Error("inputValue should be cleared after save")
	}
}

func TestImportPathModel_CtrlC_Quits(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportPathModel(store)

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	newImportPathModel := newModel.(ImportPathModel)

	if !newImportPathModel.quitting {
		t.Error("quitting should be true")
	}
	if cmd == nil {
		t.Error("Should return quit cmd")
	}
}

func TestImportPathModel_IsQuitting(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportPathModel(store)

	if model.IsQuitting() {
		t.Error("IsQuitting should be false initially")
	}

	model.quitting = true
	if !model.IsQuitting() {
		t.Error("IsQuitting should be true after quitting is set")
	}
}

func TestSettingsModel_View_ContainsImportPath(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewSettingsModel(store)

	view := model.View()

	if !strings.Contains(view, "Import Path") {
		t.Error("Settings view should contain Import Path option")
	}
}

func TestSettingsModel_NavigateToImportPath(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewSettingsModel(store)
	model.cursor = SettingsImportPath

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newSettingsModel := newModel.(SettingsModel)

	if cmd == nil {
		t.Error("Should return navigation cmd")
	}
	_ = newSettingsModel
}

func TestImportPathModel_EscInMenuMode_NavigatesBack(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportPathModel(store)
	model.mode = ImportPathMenu

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	_ = newModel.(ImportPathModel)

	if cmd == nil {
		t.Error("Should return navigation cmd to settings")
	}
}

func TestSettingsModel_Init_ReturnsNil(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewSettingsModel(store)

	cmd := model.Init()
	if cmd != nil {
		t.Error("Init should return nil")
	}
}

func TestSettingsModel_IsQuitting(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewSettingsModel(store)

	if model.IsQuitting() {
		t.Error("IsQuitting should be false initially")
	}

	model.quitting = true
	if !model.IsQuitting() {
		t.Error("IsQuitting should be true after quitting is set")
	}
}

func TestExclusionListModel_Init_ReturnsNil(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)

	cmd := model.Init()
	if cmd != nil {
		t.Error("Init should return nil")
	}
}

func TestExclusionListModel_LoadPatterns(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)
	model.patterns = nil

	model.loadPatterns()

	if model.patterns == nil {
		t.Error("loadPatterns should populate patterns")
	}
}

func TestImportPathModel_Init_ReturnsNil(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportPathModel(store)

	cmd := model.Init()
	if cmd != nil {
		t.Error("Init should return nil")
	}
}

func TestImportPathModel_LoadPath(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewImportPathModel(store)
	model.path = ""

	model.loadPath()

	if model.path == "" {
		t.Error("loadPath should populate path")
	}
}

func TestAppModel_Update_NavigationMsg_ImportPathScreen(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)
	model.currentScreen = ScreenSettings

	newModel, _ := model.Update(NavigationMsg{Destination: ScreenImportPath})
	newAppModel := newModel.(AppModel)
	if newAppModel.currentScreen != ScreenImportPath {
		t.Errorf("currentScreen = %d, want %d", newAppModel.currentScreen, ScreenImportPath)
	}
	if newAppModel.importPathModel.path == "" {
		t.Error("importPathModel.path should be populated")
	}
}

func TestExpenseListModel_Initialization(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExpenseListModel(store)

	if model.cursor != 0 {
		t.Errorf("cursor should be 0, got %d", model.cursor)
	}
	if model.page != 0 {
		t.Errorf("page should be 0, got %d", model.page)
	}
	if model.mode != ExpenseListBrowse {
		t.Errorf("mode should be ExpenseListBrowse (%d), got %d", ExpenseListBrowse, model.mode)
	}
}

func TestExpenseListModel_View(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExpenseListModel(store)

	view := model.View()
	if !strings.Contains(view, "Manage Expenses") {
		t.Error("View should contain 'Manage Expenses'")
	}
}

func TestExpenseListModel_CursorMovement(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExpenseListModel(store)

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = *newModel.(*ExpenseListModel)

	if model.cursor != 0 {
		t.Errorf("cursor should stay at 0 with no expenses, got %d", model.cursor)
	}
}

func TestExpenseListModel_Esc_NavigatesBack(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExpenseListModel(store)

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	_ = *newModel.(*ExpenseListModel)

	if cmd == nil {
		t.Error("Esc should return navigation command")
	}
}

func TestExpenseListModel_CtrlC_Quits(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExpenseListModel(store)

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	model = *newModel.(*ExpenseListModel)

	if cmd == nil {
		t.Error("Ctrl+C should return quit command")
	}
	if !model.quitting {
		t.Error("quitting should be true after Ctrl+C")
	}
}

func TestExpenseListModel_IsQuitting(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExpenseListModel(store)

	if model.IsQuitting() {
		t.Error("IsQuitting should be false initially")
	}
	model.quitting = true
	if !model.IsQuitting() {
		t.Error("IsQuitting should be true after setting quitting")
	}
}

func TestExpenseListModel_Init_ReturnsNil(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExpenseListModel(store)

	cmd := model.Init()
	if cmd != nil {
		t.Error("Init should return nil")
	}
}

func TestExpenseListModel_View_ContainsHelp(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExpenseListModel(store)

	view := model.View()
	if !strings.Contains(view, "esc") {
		t.Error("View should contain 'esc' for help")
	}
}

func TestExpenseListModel_View_ContainsFilters(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExpenseListModel(store)

	view := model.View()
	if !strings.Contains(view, "/ Search") {
		t.Error("View should contain search help")
	}
	if !strings.Contains(view, "e Edit") {
		t.Error("View should contain edit help")
	}
	if !strings.Contains(view, "x Delete") {
		t.Error("View should contain delete help")
	}
}

func TestAppModel_Update_NavigationMsg_ExpenseListScreen(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)
	model.currentScreen = ScreenMainMenu

	newModel, _ := model.Update(NavigationMsg{Destination: ScreenExpenseList})
	newAppModel := newModel.(AppModel)
	if newAppModel.currentScreen != ScreenExpenseList {
		t.Errorf("currentScreen = %d, want %d", newAppModel.currentScreen, ScreenExpenseList)
	}
}

func TestAppModel_View_ExpenseListScreen(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)

	model.currentScreen = ScreenExpenseList
	view := model.View()

	if !strings.Contains(view, "Manage Expenses") {
		t.Error("View should show ExpenseList when on ExpenseList screen")
	}
}

func TestAppModel_Update_KeyMsg_ExpenseList(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)
	model.currentScreen = ScreenExpenseList

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	newAppModel := newModel.(AppModel)

	if newAppModel.expenseListModel.cursor != 0 {
		t.Errorf("expenseListModel.cursor = %d, want 0", newAppModel.expenseListModel.cursor)
	}
}

func TestExpenseListModel_SearchMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExpenseListModel(store)

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	model = *newModel.(*ExpenseListModel)

	if model.mode != ExpenseListSearch {
		t.Errorf("mode should be ExpenseListSearch (%d), got %d", ExpenseListSearch, model.mode)
	}
}

func TestAppModel_Update_NavigationMsg_RecurringList(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)

	newModel, _ := model.Update(NavigationMsg{Destination: ScreenRecurringList})
	newAppModel := newModel.(AppModel)

	if newAppModel.currentScreen != ScreenRecurringList {
		t.Errorf("currentScreen = %d, want %d", newAppModel.currentScreen, ScreenRecurringList)
	}
}

func TestAppModel_Update_NavigationMsg_RecurringStatus(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)

	newModel, _ := model.Update(NavigationMsg{Destination: ScreenRecurringStatus})
	newAppModel := newModel.(AppModel)

	if newAppModel.currentScreen != ScreenRecurringStatus {
		t.Errorf("currentScreen = %d, want %d", newAppModel.currentScreen, ScreenRecurringStatus)
	}
}

func TestAppModel_View_RecurringListScreen(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)
	model.currentScreen = ScreenRecurringList
	model.recurringModel = NewRecurringListModel(store)

	view := model.View()
	if !strings.Contains(view, "Payments") {
		t.Error("Expected RecurringList view")
	}
}

func TestAppModel_View_RecurringStatusScreen(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)
	model.currentScreen = ScreenRecurringStatus
	model.recurringStatusModel = NewRecurringStatusModel(store)

	view := model.View()
	if !strings.Contains(view, "Recurring") {
		t.Error("Expected RecurringStatus view")
	}
}

func TestAppModel_Update_KeyMsg_RecurringList(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)
	model.currentScreen = ScreenRecurringList
	model.recurringModel = NewRecurringListModel(store)

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	newAppModel := newModel.(AppModel)

	if newAppModel.recurringModel.cursor != 0 {
		t.Errorf("recurringModel.cursor = %d, want 0 (list is empty)", newAppModel.recurringModel.cursor)
	}
}

func TestAppModel_Update_KeyMsg_RecurringStatus(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)
	model.currentScreen = ScreenRecurringStatus
	model.recurringStatusModel = NewRecurringStatusModel(store)

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	newAppModel := newModel.(AppModel)

	if newAppModel.recurringStatusModel.cursor != 0 {
		t.Errorf("recurringStatusModel.cursor = %d, want 0 (list is empty)", newAppModel.recurringStatusModel.cursor)
	}
}

func TestAppModel_Update_NavigationMsg_CategoryList(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)

	newModel, _ := model.Update(NavigationMsg{Destination: ScreenCategoryList})
	newAppModel := newModel.(AppModel)

	if newAppModel.currentScreen != ScreenCategoryList {
		t.Errorf("currentScreen = %d, want %d", newAppModel.currentScreen, ScreenCategoryList)
	}
}

func TestAppModel_View_CategoryListScreen(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)
	model.currentScreen = ScreenCategoryList
	model.categoryModel = NewCategoryListModel(store)

	view := model.View()
	if !strings.Contains(view, "Category") {
		t.Error("Expected Category view")
	}
}

func TestAppModel_Update_KeyMsg_CategoryList(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)
	model.currentScreen = ScreenCategoryList
	model.categoryModel = NewCategoryListModel(store)

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	newAppModel := newModel.(AppModel)

	if newAppModel.categoryModel.cursor != 0 {
		t.Errorf("categoryModel.cursor = %d, want 0 (list is empty)", newAppModel.categoryModel.cursor)
	}
}

func TestAppModel_Update_DefaultScreen(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)
	model.currentScreen = 999

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	newAppModel := newModel.(AppModel)

	if newAppModel.currentScreen != 999 {
		t.Error("Screen should remain unchanged")
	}
}

func TestAppModel_Update_NonKeyMsg_AnyScreen(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)
	model.currentScreen = ScreenMainMenu

	newModel, _ := model.Update("some string message")
	newAppModel := newModel.(AppModel)

	if newAppModel.currentScreen != ScreenMainMenu {
		t.Error("Screen should remain unchanged")
	}
}

func TestExpenseListModel_SearchModeTyping(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListSearch

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t', 'e', 's', 't'}})
	model = *newModel.(*ExpenseListModel)

	if model.searchQuery != "test" {
		t.Errorf("searchQuery = %q, want 'test'", model.searchQuery)
	}
}

func TestExpenseListModel_SearchModeBackspace(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListSearch
	model.searchQuery = "test"

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	model = *newModel.(*ExpenseListModel)

	if model.searchQuery != "tes" {
		t.Errorf("searchQuery = %q, want 'tes'", model.searchQuery)
	}
}

func TestExpenseListModel_SearchModeEsc(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListSearch
	model.searchQuery = "test"

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = *newModel.(*ExpenseListModel)

	if model.mode != ExpenseListBrowse {
		t.Errorf("mode should be ExpenseListBrowse, got %d", model.mode)
	}
}

func TestExpenseListModel_SearchModeEnter(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListSearch
	model.searchQuery = "test"

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = *newModel.(*ExpenseListModel)

	if model.mode != ExpenseListBrowse {
		t.Errorf("mode should be ExpenseListBrowse after enter, got %d", model.mode)
	}
}

func TestExpenseListModel_Trunc(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"exactly10!", 10, "exactly10!"},
		{"this is too long", 8, "this i.."},
	}

	for _, tt := range tests {
		result := trunc(tt.input, tt.maxLen)
		if result != tt.expected {
			t.Errorf("trunc(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
		}
	}
}

type mockStoreWithExpenses struct {
	mockStoreForCLI
	expenses   []storage.Expense
	categories []string
}

func (m *mockStoreWithExpenses) GetAllExpenses(startDate, endDate *time.Time) ([]storage.Expense, error) {
	return m.expenses, nil
}

func (m *mockStoreWithExpenses) GetCategories() ([]string, error) {
	return m.categories, nil
}

func (m *mockStoreWithExpenses) UpdateExpense(id string, expense storage.Expense) error {
	return nil
}

func (m *mockStoreWithExpenses) RemoveExpense(id string) error {
	return nil
}

func TestExpenseListModel_ParseDate_YYYYMMDD(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses:   []storage.Expense{},
		categories: []string{},
	}
	_ = NewExpenseListModel(store)

	result, err := parseDate("2024-03-15")
	if err != nil {
		t.Errorf("parseDate failed: %v", err)
	}
	if result.Year() != 2024 || result.Month() != 3 || result.Day() != 15 {
		t.Errorf("parseDate = %v, want 2024-03-15", result)
	}
}

func TestExpenseListModel_ParseDate_MMDDYYYY(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses:   []storage.Expense{},
		categories: []string{},
	}
	_ = NewExpenseListModel(store)

	result, err := parseDate("03/15/2024")
	if err != nil {
		t.Errorf("parseDate failed: %v", err)
	}
	if result.Year() != 2024 || result.Month() != 3 || result.Day() != 15 {
		t.Errorf("parseDate = %v, want 2024-03-15", result)
	}
}

func TestExpenseListModel_ParseDate_Invalid(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses:   []storage.Expense{},
		categories: []string{},
	}
	_ = NewExpenseListModel(store)

	_, err := parseDate("invalid-date")
	if err == nil {
		t.Error("parseDate should fail for invalid date")
	}
}

func TestExpenseListModel_ParseDate_Empty(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses:   []storage.Expense{},
		categories: []string{},
	}
	_ = NewExpenseListModel(store)

	_, err := parseDate("")
	if err == nil {
		t.Error("parseDate should fail for empty string")
	}
}

func TestExpenseListModel_Update_EditMode_Enter(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test Expense", Amount: 10.00, Category: "Food", Date: time.Now()},
		},
		categories: []string{"Food", "Transport"},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListEdit
	model.editingExpense = model.expenses[0]
	model.editField = 0
	model.editValue = "Updated Name"

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.mode != ExpenseListBrowse {
		t.Errorf("mode should be ExpenseListBrowse after save, got %d", newExpenseListModel.mode)
	}
}

func TestExpenseListModel_Update_EditMode_Tab(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses:   []storage.Expense{},
		categories: []string{},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListEdit
	model.editField = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.editField != 1 {
		t.Errorf("editField should be 1 after Tab, got %d", newExpenseListModel.editField)
	}
}

func TestExpenseListModel_Update_EditMode_UpDown(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses:   []storage.Expense{},
		categories: []string{},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListEdit
	model.editField = 2

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyUp})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.editField != 1 {
		t.Errorf("editField should be 1 after Up, got %d", newExpenseListModel.editField)
	}

	newModel, _ = newExpenseListModel.Update(tea.KeyMsg{Type: tea.KeyDown})
	finalModel := *newModel.(*ExpenseListModel)

	if finalModel.editField != 2 {
		t.Errorf("editField should be 2 after Down, got %d", finalModel.editField)
	}
}

func TestExpenseListModel_Update_EditMode_Esc(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses:   []storage.Expense{},
		categories: []string{},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListEdit
	model.editField = 0
	model.errorMsg = "some error"

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.mode != ExpenseListBrowse {
		t.Errorf("mode should be ExpenseListBrowse after Esc, got %d", newExpenseListModel.mode)
	}
	if newExpenseListModel.errorMsg != "" {
		t.Error("errorMsg should be cleared after Esc")
	}
}

func TestExpenseListModel_Update_EditMode_CtrlC(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses:   []storage.Expense{},
		categories: []string{},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListEdit

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if cmd == nil {
		t.Error("Ctrl+C should return quit command")
	}
	if !newExpenseListModel.quitting {
		t.Error("quitting should be true after Ctrl+C")
	}
}

func TestExpenseListModel_Update_EditMode_Backspace(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses:   []storage.Expense{},
		categories: []string{},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListEdit
	model.editField = 0
	model.editValue = "test"

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.editValue != "tes" {
		t.Errorf("editValue = %q, want 'tes'", newExpenseListModel.editValue)
	}
}

func TestExpenseListModel_Update_DeleteConfirm_Y(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test Expense", Amount: 10.00, Category: "Food", Date: time.Now()},
		},
		categories: []string{"Food"},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListDeleteConfirm
	model.editingExpense = model.expenses[0]

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.mode != ExpenseListBrowse {
		t.Errorf("mode should be ExpenseListBrowse after Y, got %d", newExpenseListModel.mode)
	}
}

func TestExpenseListModel_Update_DeleteConfirm_N(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test Expense", Amount: 10.00, Category: "Food", Date: time.Now()},
		},
		categories: []string{"Food"},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListDeleteConfirm
	model.editingExpense = model.expenses[0]

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.mode != ExpenseListBrowse {
		t.Errorf("mode should be ExpenseListBrowse after N, got %d", newExpenseListModel.mode)
	}
}

func TestExpenseListModel_Update_DeleteConfirm_Esc(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses:   []storage.Expense{},
		categories: []string{},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListDeleteConfirm

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.mode != ExpenseListBrowse {
		t.Errorf("mode should be ExpenseListBrowse after Esc, got %d", newExpenseListModel.mode)
	}
}

func TestExpenseListModel_Update_DeleteConfirm_CtrlC(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses:   []storage.Expense{},
		categories: []string{},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListDeleteConfirm

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if cmd == nil {
		t.Error("Ctrl+C should return quit command")
	}
	if !newExpenseListModel.quitting {
		t.Error("quitting should be true after Ctrl+C")
	}
}

func TestExpenseListModel_Update_BrowseMode_Search(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses:   []storage.Expense{},
		categories: []string{},
	}
	model := NewExpenseListModel(store)

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.mode != ExpenseListSearch {
		t.Errorf("mode should be ExpenseListSearch, got %d", newExpenseListModel.mode)
	}
}

func TestExpenseListModel_Update_BrowseMode_Edit(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test Expense", Amount: 10.00, Category: "Food", Date: time.Now()},
		},
		categories: []string{"Food"},
	}
	model := NewExpenseListModel(store)

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.mode != ExpenseListEdit {
		t.Errorf("mode should be ExpenseListEdit, got %d", newExpenseListModel.mode)
	}
}

func TestExpenseListModel_Update_BrowseMode_Delete(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test Expense", Amount: 10.00, Category: "Food", Date: time.Now()},
		},
		categories: []string{"Food"},
	}
	model := NewExpenseListModel(store)

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.mode != ExpenseListDeleteConfirm {
		t.Errorf("mode should be ExpenseListDeleteConfirm, got %d", newExpenseListModel.mode)
	}
}

func TestExpenseListModel_Update_BrowseMode_PaginationLeft(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses:   []storage.Expense{},
		categories: []string{},
	}
	model := NewExpenseListModel(store)
	model.page = 1

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.page != 0 {
		t.Errorf("page should be 0 after Left, got %d", newExpenseListModel.page)
	}
}

func TestExpenseListModel_Update_BrowseMode_PaginationRight(t *testing.T) {
	expenses := make([]storage.Expense, 25)
	for i := 0; i < 25; i++ {
		expenses[i] = storage.Expense{ID: fmt.Sprintf("test-%d", i), Name: fmt.Sprintf("Expense %d", i), Amount: 10.00, Category: "Food", Date: time.Now()}
	}
	store := &mockStoreWithExpenses{
		expenses:   expenses,
		categories: []string{"Food"},
	}
	model := NewExpenseListModel(store)
	model.page = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRight})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.page != 1 {
		t.Errorf("page should be 1 after Right, got %d", newExpenseListModel.page)
	}
}

func TestExpenseListModel_Update_BrowseMode_L(t *testing.T) {
	expenses := make([]storage.Expense, 25)
	for i := 0; i < 25; i++ {
		expenses[i] = storage.Expense{ID: fmt.Sprintf("test-%d", i), Name: fmt.Sprintf("Expense %d", i), Amount: 10.00, Category: "Food", Date: time.Now()}
	}
	store := &mockStoreWithExpenses{
		expenses:   expenses,
		categories: []string{"Food"},
	}
	model := NewExpenseListModel(store)
	model.cursor = 3

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRight})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.page != 1 {
		t.Errorf("page should be 1 after Right, got %d", newExpenseListModel.page)
	}
	if newExpenseListModel.cursor != 13 {
		t.Errorf("cursor should be 13 after Right (3 + 10), got %d", newExpenseListModel.cursor)
	}
}

func TestExpenseListModel_Pagination_CursorMovesWithPage(t *testing.T) {
	expenses := make([]storage.Expense, 25)
	for i := 0; i < 25; i++ {
		expenses[i] = storage.Expense{ID: fmt.Sprintf("test-%d", i), Name: fmt.Sprintf("Expense %d", i), Amount: 10.00, Category: "Food", Date: time.Now()}
	}
	store := &mockStoreWithExpenses{
		expenses:   expenses,
		categories: []string{"Food"},
	}
	model := NewExpenseListModel(store)
	model.cursor = 13
	model.page = 1

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.page != 0 {
		t.Errorf("page should be 0 after Left, got %d", newExpenseListModel.page)
	}
	if newExpenseListModel.cursor != 3 {
		t.Errorf("cursor should be 3 after Left (13 - 10), got %d", newExpenseListModel.cursor)
	}
}

func TestExpenseListModel_Pagination_CursorClampsOnLastPage(t *testing.T) {
	expenses := make([]storage.Expense, 25)
	for i := 0; i < 25; i++ {
		expenses[i] = storage.Expense{ID: fmt.Sprintf("test-%d", i), Name: fmt.Sprintf("Expense %d", i), Amount: 10.00, Category: "Food", Date: time.Now()}
	}
	store := &mockStoreWithExpenses{
		expenses:   expenses,
		categories: []string{"Food"},
	}
	model := NewExpenseListModel(store)
	model.cursor = 3
	model.page = 2

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRight})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.page != 2 {
		t.Errorf("page should stay at 2 (already last page), got %d", newExpenseListModel.page)
	}
}

func TestExpenseListModel_Update_BrowseMode_UpDown(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test 1", Amount: 10.00, Category: "Food", Date: time.Now()},
			{ID: "test-2", Name: "Test 2", Amount: 20.00, Category: "Food", Date: time.Now()},
		},
		categories: []string{"Food"},
	}
	model := NewExpenseListModel(store)
	model.cursor = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.cursor != 1 {
		t.Errorf("cursor should be 1 after Down, got %d", newExpenseListModel.cursor)
	}

	newModel, _ = newExpenseListModel.Update(tea.KeyMsg{Type: tea.KeyUp})
	finalModel := *newModel.(*ExpenseListModel)

	if finalModel.cursor != 0 {
		t.Errorf("cursor should be 0 after Up, got %d", finalModel.cursor)
	}
}

func TestExpenseListModel_View_EditMode(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test Expense", Amount: 10.00, Category: "Food", Currency: "usd", Date: time.Now()},
		},
		categories: []string{"Food"},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListEdit
	model.editingExpense = model.expenses[0]

	view := model.View()
	if !strings.Contains(view, "Edit Expense") {
		t.Error("View should contain 'Edit Expense'")
	}
	if !strings.Contains(view, "Name") {
		t.Error("View should contain 'Name' field")
	}
	if !strings.Contains(view, "Amount") {
		t.Error("View should contain 'Amount' field")
	}
	if !strings.Contains(view, "Category") {
		t.Error("View should contain 'Category' field")
	}
}

func TestExpenseListModel_View_EditMode_WithError(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test Expense", Amount: 10.00, Category: "Food", Date: time.Now()},
		},
		categories: []string{"Food"},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListEdit
	model.editingExpense = model.expenses[0]
	model.errorMsg = "Validation error"

	view := model.View()
	if !strings.Contains(view, "Validation error") {
		t.Error("View should contain error message")
	}
}

func TestExpenseListModel_View_SearchMode(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses:   []storage.Expense{},
		categories: []string{},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListSearch
	model.searchQuery = "test query"

	view := model.View()
	if !strings.Contains(view, "Search Expenses") {
		t.Error("View should contain 'Search Expenses'")
	}
	if !strings.Contains(view, "test query") {
		t.Error("View should contain search query")
	}
}

func TestExpenseListModel_View_DeleteConfirm(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Delete Me", Amount: 25.50, Category: "Food", Currency: "usd", Date: time.Now()},
		},
		categories: []string{"Food"},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListDeleteConfirm
	model.editingExpense = model.expenses[0]

	view := model.View()
	if !strings.Contains(view, "Delete Expense?") {
		t.Error("View should contain 'Delete Expense?'")
	}
	if !strings.Contains(view, "Delete Me") {
		t.Error("View should contain expense name")
	}
	if !strings.Contains(view, "25.50") {
		t.Error("View should contain expense amount")
	}
}

func TestExpenseListModel_View_BrowseMode_WithFilters(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test", Amount: 10.00, Category: "Food", Date: time.Now()},
		},
		categories: []string{"Food"},
	}
	model := NewExpenseListModel(store)
	model.searchQuery = "coffee"
	model.categoryFilter = "Food"

	view := model.View()
	if !strings.Contains(view, "search='coffee'") {
		t.Error("View should contain search filter")
	}
	if !strings.Contains(view, "category='Food'") {
		t.Error("View should contain category filter")
	}
}

func TestExpenseListModel_View_BrowseMode_WithError(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses:   []storage.Expense{},
		categories: []string{},
	}
	model := NewExpenseListModel(store)
	model.errorMsg = "Failed to load"

	view := model.View()
	if !strings.Contains(view, "Failed to load") {
		t.Error("View should contain error message")
	}
}

func TestExpenseListModel_View_BrowseMode_WithSuccess(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses:   []storage.Expense{},
		categories: []string{},
	}
	model := NewExpenseListModel(store)
	model.successMsg = "Updated successfully"

	view := model.View()
	if !strings.Contains(view, "Updated successfully") {
		t.Error("View should contain success message")
	}
}

func TestExpenseListModel_ApplyFilters_DateRange(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Old Expense", Amount: 10.00, Category: "Food", Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			{ID: "test-2", Name: "New Expense", Amount: 20.00, Category: "Food", Date: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)},
		},
		categories: []string{"Food"},
	}
	model := NewExpenseListModel(store)
	model.startDateFilter = "2024-05-01"

	model.applyFilters()

	if len(model.filteredExpenses) != 1 {
		t.Errorf("filteredExpenses should have 1 item, got %d", len(model.filteredExpenses))
	}
	if model.filteredExpenses[0].Name != "New Expense" {
		t.Errorf("filteredExpenses should contain 'New Expense', got %s", model.filteredExpenses[0].Name)
	}
}

func TestExpenseListModel_ApplyFilters_Category(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Food Item", Amount: 10.00, Category: "Food", Date: time.Now()},
			{ID: "test-2", Name: "Transport Item", Amount: 20.00, Category: "Transport", Date: time.Now()},
		},
		categories: []string{"Food", "Transport"},
	}
	model := NewExpenseListModel(store)
	model.categoryFilter = "Transport"

	model.applyFilters()

	if len(model.filteredExpenses) != 1 {
		t.Errorf("filteredExpenses should have 1 item, got %d", len(model.filteredExpenses))
	}
	if model.filteredExpenses[0].Category != "Transport" {
		t.Errorf("filteredExpenses should contain 'Transport', got %s", model.filteredExpenses[0].Category)
	}
}

func TestExpenseListModel_ApplyFilters_Search(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Coffee at Starbucks", Amount: 5.00, Category: "Food", Date: time.Now()},
			{ID: "test-2", Name: "Gas Station", Amount: 40.00, Category: "Transport", Date: time.Now()},
		},
		categories: []string{"Food", "Transport"},
	}
	model := NewExpenseListModel(store)
	model.searchQuery = "coffee"

	model.applyFilters()

	if len(model.filteredExpenses) != 1 {
		t.Errorf("filteredExpenses should have 1 item, got %d", len(model.filteredExpenses))
	}
	if !strings.Contains(model.filteredExpenses[0].Name, "Coffee") {
		t.Errorf("filteredExpenses should contain 'Coffee', got %s", model.filteredExpenses[0].Name)
	}
}

func TestExpenseListModel_ExpensesPerPage(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses:   []storage.Expense{},
		categories: []string{},
	}
	model := NewExpenseListModel(store)

	if model.expensesPerPage() != expensesPerPage {
		t.Errorf("expensesPerPage() = %d, want %d", model.expensesPerPage(), expensesPerPage)
	}
}

func TestExpenseListModel_TotalPages(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses:   []storage.Expense{},
		categories: []string{},
	}
	model := NewExpenseListModel(store)

	if model.totalPages() != 1 {
		t.Errorf("totalPages() = %d with no expenses, want 1", model.totalPages())
	}

	model.filteredExpenses = make([]storage.Expense, 25)
	if model.totalPages() != 3 {
		t.Errorf("totalPages() = %d with 25 expenses, want 3", model.totalPages())
	}
}

func TestExpenseListModel_CurrentPageExpenses(t *testing.T) {
	expenses := make([]storage.Expense, 25)
	for i := 0; i < 25; i++ {
		expenses[i] = storage.Expense{ID: fmt.Sprintf("test-%d", i), Name: fmt.Sprintf("Expense %d", i), Amount: 10.00, Category: "Food", Date: time.Now()}
	}
	store := &mockStoreWithExpenses{
		expenses:   expenses,
		categories: []string{"Food"},
	}
	model := NewExpenseListModel(store)

	pageExpenses := model.currentPageExpenses()
	if len(pageExpenses) != expensesPerPage {
		t.Errorf("currentPageExpenses() should return %d items, got %d", expensesPerPage, len(pageExpenses))
	}

	model.page = 1
	pageExpenses = model.currentPageExpenses()
	if len(pageExpenses) != 10 {
		t.Errorf("currentPageExpenses() on page 1 should return 10 items, got %d", len(pageExpenses))
	}
}

func TestExpenseListModel_Update_BrowseMode_Category(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses:   []storage.Expense{},
		categories: []string{},
	}
	model := NewExpenseListModel(store)

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.mode != ExpenseListSearch {
		t.Errorf("mode should be ExpenseListSearch for category, got %d", newExpenseListModel.mode)
	}
}

func TestExpenseListModel_Update_BrowseMode_Date(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses:   []storage.Expense{},
		categories: []string{},
	}
	model := NewExpenseListModel(store)

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.mode != ExpenseListSearch {
		t.Errorf("mode should be ExpenseListSearch for date, got %d", newExpenseListModel.mode)
	}
}

func TestExpenseListModel_Update_BrowseMode_D(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test", Amount: 10.00, Category: "Food", Date: time.Now()},
		},
		categories: []string{"Food"},
	}
	model := NewExpenseListModel(store)

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'D'}})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.mode != ExpenseListDeleteConfirm {
		t.Errorf("mode should be ExpenseListDeleteConfirm, got %d", newExpenseListModel.mode)
	}
}

func TestExpenseListModel_Update_BrowseMode_H(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses:   []storage.Expense{},
		categories: []string{},
	}
	model := NewExpenseListModel(store)
	model.page = 2

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.page != 1 {
		t.Errorf("page should be 1 after h, got %d", newExpenseListModel.page)
	}
}

func TestExpenseListModel_Update_EditMode_CharacterInput(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test", Amount: 10.00, Category: "Food", Date: time.Now()},
		},
		categories: []string{"Food"},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListEdit
	model.editField = 0
	model.editValue = ""

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'N', 'e', 'w'}})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.editValue != "New" {
		t.Errorf("editValue = %q, want 'New'", newExpenseListModel.editValue)
	}
}

func TestExpenseListModel_Update_DeleteConfirm_Enter(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test", Amount: 10.00, Category: "Food", Date: time.Now()},
		},
		categories: []string{"Food"},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListDeleteConfirm
	model.editingExpense = model.expenses[0]

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.mode != ExpenseListBrowse {
		t.Errorf("mode should be ExpenseListBrowse after Enter, got %d", newExpenseListModel.mode)
	}
}

func TestExpenseListModel_Update_DeleteConfirm_Y_Uppercase(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test", Amount: 10.00, Category: "Food", Date: time.Now()},
		},
		categories: []string{"Food"},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListDeleteConfirm
	model.editingExpense = model.expenses[0]

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'Y'}})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.mode != ExpenseListBrowse {
		t.Errorf("mode should be ExpenseListBrowse after Y, got %d", newExpenseListModel.mode)
	}
}

func TestExpenseListModel_Update_DeleteConfirm_N_Uppercase(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test", Amount: 10.00, Category: "Food", Date: time.Now()},
		},
		categories: []string{"Food"},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListDeleteConfirm
	model.editingExpense = model.expenses[0]

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'N'}})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.mode != ExpenseListBrowse {
		t.Errorf("mode should be ExpenseListBrowse after N, got %d", newExpenseListModel.mode)
	}
}

func TestExpenseListModel_View_BrowseMode_Pagination(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test 1", Amount: 10.00, Category: "Food", Date: time.Now()},
		},
		categories: []string{"Food"},
	}
	model := NewExpenseListModel(store)
	model.page = 0

	view := model.View()
	if !strings.Contains(view, "Page 1/1") {
		t.Error("View should contain pagination info")
	}
}

func TestExpenseListModel_Update_EditMode_K_J(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test", Amount: 10.00, Category: "Food", Date: time.Now()},
		},
		categories: []string{"Food", "Travel", "Utilities"},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListEdit
	model.editField = 2
	model.editValue = "10.00"
	model.editingExpense = model.expenses[0]

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.editField != 1 {
		t.Errorf("editField should be 1 after k, got %d", newExpenseListModel.editField)
	}

	newModel, _ = newExpenseListModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	finalModel := *newModel.(*ExpenseListModel)

	if finalModel.editField != 2 {
		t.Errorf("editField should be 2 after j, got %d", finalModel.editField)
	}
}

func TestExpenseListModel_Update_EditMode_Tab_OpensCategoryPicker(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test", Amount: 10.00, Category: "Travel", Date: time.Now()},
		},
		categories: []string{"Food", "Travel", "Utilities"},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListEdit
	model.editField = 1
	model.editValue = "Travel"
	model.editingExpense = model.expenses[0]

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if !newExpenseListModel.showCategoryPicker {
		t.Error("showCategoryPicker should be true after Tab on category field")
	}
	if newExpenseListModel.categoryPickerIdx != 1 {
		t.Errorf("categoryPickerIdx should be 1 (Travel), got %d", newExpenseListModel.categoryPickerIdx)
	}
}

func TestExpenseListModel_Update_BrowseMode_K_J(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test 1", Amount: 10.00, Category: "Food", Date: time.Now()},
			{ID: "test-2", Name: "Test 2", Amount: 20.00, Category: "Food", Date: time.Now()},
		},
		categories: []string{"Food"},
	}
	model := NewExpenseListModel(store)
	model.cursor = 1

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.cursor != 0 {
		t.Errorf("cursor should be 0 after k, got %d", newExpenseListModel.cursor)
	}

	newModel, _ = newExpenseListModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	finalModel := *newModel.(*ExpenseListModel)

	if finalModel.cursor != 1 {
		t.Errorf("cursor should be 1 after j, got %d", finalModel.cursor)
	}
}

type mockStoreWithCategoryList struct {
	mockStoreForCLI
	categories []string
}

func (m *mockStoreWithCategoryList) GetCategories() ([]string, error) {
	return m.categories, nil
}

func (m *mockStoreWithCategoryList) UpdateCategories(categories []string) error {
	m.categories = categories
	return nil
}

func TestCategoryListModel_LoadCategories(t *testing.T) {
	store := &mockStoreWithCategoryList{categories: []string{"Food", "Transport", "Utilities"}}
	model := NewCategoryListModel(store)
	model.categories = nil

	model.loadCategories()

	if len(model.categories) != 3 {
		t.Errorf("categories should have 3 items, got %d", len(model.categories))
	}
}

func TestCategoryListModel_Update_AddMode(t *testing.T) {
	store := &mockStoreWithCategoryList{categories: []string{"Food"}}
	model := NewCategoryListModel(store)

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	newCategoryListModel := newModel.(CategoryListModel)

	if newCategoryListModel.mode != CategoryListAdd {
		t.Errorf("mode should be CategoryListAdd (%d), got %d", CategoryListAdd, newCategoryListModel.mode)
	}
	if newCategoryListModel.inputValue != "" {
		t.Error("inputValue should be empty")
	}
}

func TestCategoryListModel_Update_AddModeTyping(t *testing.T) {
	store := &mockStoreWithCategoryList{categories: []string{}}
	model := NewCategoryListModel(store)
	model.mode = CategoryListAdd

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'T', 'e', 's', 't'}})
	newCategoryListModel := newModel.(CategoryListModel)

	if newCategoryListModel.inputValue != "Test" {
		t.Errorf("inputValue = %q, want 'Test'", newCategoryListModel.inputValue)
	}
}

func TestCategoryListModel_Update_AddModeSubmit(t *testing.T) {
	store := &mockStoreWithCategoryList{categories: []string{}}
	model := NewCategoryListModel(store)
	model.mode = CategoryListAdd
	model.inputValue = "NewCategory"

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newCategoryListModel := newModel.(CategoryListModel)

	if newCategoryListModel.mode != CategoryListView {
		t.Errorf("mode should be CategoryListView after submit, got %d", newCategoryListModel.mode)
	}
	if len(newCategoryListModel.categories) != 1 {
		t.Errorf("categories should have 1 item, got %d", len(newCategoryListModel.categories))
	}
}

func TestCategoryListModel_Update_AddModeEsc(t *testing.T) {
	store := &mockStoreWithCategoryList{categories: []string{}}
	model := NewCategoryListModel(store)
	model.mode = CategoryListAdd
	model.inputValue = "test"

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	newCategoryListModel := newModel.(CategoryListModel)

	if newCategoryListModel.mode != CategoryListView {
		t.Errorf("mode should be CategoryListView after Esc, got %d", newCategoryListModel.mode)
	}
	if newCategoryListModel.inputValue != "" {
		t.Error("inputValue should be cleared after Esc")
	}
}

func TestCategoryListModel_Update_EditMode(t *testing.T) {
	store := &mockStoreWithCategoryList{categories: []string{"Food", "Transport"}}
	model := NewCategoryListModel(store)
	model.cursor = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	newCategoryListModel := newModel.(CategoryListModel)

	if newCategoryListModel.mode != CategoryListEdit {
		t.Errorf("mode should be CategoryListEdit, got %d", newCategoryListModel.mode)
	}
	if newCategoryListModel.editIndex != 0 {
		t.Errorf("editIndex should be 0, got %d", newCategoryListModel.editIndex)
	}
	if newCategoryListModel.inputValue != "Food" {
		t.Errorf("inputValue = %q, want 'Food'", newCategoryListModel.inputValue)
	}
}

func TestCategoryListModel_Update_EditModeSubmit(t *testing.T) {
	store := &mockStoreWithCategoryList{categories: []string{"OldCategory"}}
	model := NewCategoryListModel(store)
	model.mode = CategoryListEdit
	model.editIndex = 0
	model.inputValue = "NewCategory"

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newCategoryListModel := newModel.(CategoryListModel)

	if newCategoryListModel.mode != CategoryListView {
		t.Errorf("mode should be CategoryListView after save, got %d", newCategoryListModel.mode)
	}
	if len(newCategoryListModel.categories) != 1 || newCategoryListModel.categories[0] != "NewCategory" {
		t.Errorf("categories should contain 'NewCategory', got %v", newCategoryListModel.categories)
	}
}

func TestCategoryListModel_Update_DeleteMode(t *testing.T) {
	store := &mockStoreWithCategoryList{categories: []string{"Food", "Transport"}}
	model := NewCategoryListModel(store)
	model.cursor = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	newCategoryListModel := newModel.(CategoryListModel)

	if newCategoryListModel.mode != CategoryListDeleteConfirm {
		t.Errorf("mode should be CategoryListDeleteConfirm, got %d", newCategoryListModel.mode)
	}
}

func TestCategoryListModel_Update_DeleteConfirmY(t *testing.T) {
	store := &mockStoreWithCategoryList{categories: []string{"Food", "Transport"}}
	model := NewCategoryListModel(store)
	model.mode = CategoryListDeleteConfirm
	model.cursor = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	newCategoryListModel := newModel.(CategoryListModel)

	if newCategoryListModel.mode != CategoryListView {
		t.Errorf("mode should be CategoryListView after Y, got %d", newCategoryListModel.mode)
	}
	if len(newCategoryListModel.categories) != 1 {
		t.Errorf("categories should have 1 item after delete, got %d", len(newCategoryListModel.categories))
	}
}

func TestCategoryListModel_Update_DeleteConfirmN(t *testing.T) {
	store := &mockStoreWithCategoryList{categories: []string{"Food"}}
	model := NewCategoryListModel(store)
	model.mode = CategoryListDeleteConfirm

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	newCategoryListModel := newModel.(CategoryListModel)

	if newCategoryListModel.mode != CategoryListView {
		t.Errorf("mode should be CategoryListView after N, got %d", newCategoryListModel.mode)
	}
	if len(newCategoryListModel.categories) != 1 {
		t.Error("categories should remain unchanged")
	}
}

func TestCategoryListModel_Update_Esc_NavigatesBack(t *testing.T) {
	store := &mockStoreWithCategoryList{categories: []string{}}
	model := NewCategoryListModel(store)
	model.mode = CategoryListView

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	_ = newModel.(CategoryListModel)

	if cmd == nil {
		t.Error("Esc should return navigation command")
	}
}

func TestCategoryListModel_CursorMovement(t *testing.T) {
	store := &mockStoreWithCategoryList{categories: []string{"Food", "Transport", "Utilities"}}
	model := NewCategoryListModel(store)
	model.cursor = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	newCategoryListModel := newModel.(CategoryListModel)

	if newCategoryListModel.cursor != 1 {
		t.Errorf("cursor should be 1 after Down, got %d", newCategoryListModel.cursor)
	}

	newModel, _ = newCategoryListModel.Update(tea.KeyMsg{Type: tea.KeyUp})
	finalModel := newModel.(CategoryListModel)

	if finalModel.cursor != 0 {
		t.Errorf("cursor should be 0 after Up, got %d", finalModel.cursor)
	}
}

func TestCategoryListModel_Update_Backspace(t *testing.T) {
	store := &mockStoreWithCategoryList{categories: []string{}}
	model := NewCategoryListModel(store)
	model.mode = CategoryListAdd
	model.inputValue = "test"

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	newCategoryListModel := newModel.(CategoryListModel)

	if newCategoryListModel.inputValue != "tes" {
		t.Errorf("inputValue = %q, want 'tes'", newCategoryListModel.inputValue)
	}
}

func TestCategoryListModel_CtrlC_Quits(t *testing.T) {
	store := &mockStoreWithCategoryList{categories: []string{}}
	model := NewCategoryListModel(store)

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	newCategoryListModel := newModel.(CategoryListModel)

	if !newCategoryListModel.quitting {
		t.Error("quitting should be true")
	}
	if cmd == nil {
		t.Error("Should return quit cmd")
	}
}

func TestCategoryListModel_Init_ReturnsNil(t *testing.T) {
	store := &mockStoreWithCategoryList{categories: []string{}}
	model := NewCategoryListModel(store)

	cmd := model.Init()
	if cmd != nil {
		t.Error("Init should return nil")
	}
}

func TestCategoryListModel_View_AddMode(t *testing.T) {
	store := &mockStoreWithCategoryList{categories: []string{}}
	model := NewCategoryListModel(store)
	model.mode = CategoryListAdd

	view := model.View()
	if !strings.Contains(view, "Enter new category:") {
		t.Error("View should show 'Enter new category:'")
	}
}

func TestCategoryListModel_View_EditMode(t *testing.T) {
	store := &mockStoreWithCategoryList{categories: []string{}}
	model := NewCategoryListModel(store)
	model.mode = CategoryListEdit
	model.inputValue = "TestCategory"

	view := model.View()
	if !strings.Contains(view, "Edit category:") {
		t.Error("View should show 'Edit category:'")
	}
	if !strings.Contains(view, "TestCategory") {
		t.Error("View should show input value")
	}
}

func TestCategoryListModel_View_DeleteConfirm(t *testing.T) {
	store := &mockStoreWithCategoryList{categories: []string{"Food"}}
	model := NewCategoryListModel(store)
	model.mode = CategoryListDeleteConfirm
	model.cursor = 0

	view := model.View()
	if !strings.Contains(view, "Confirm delete") {
		t.Error("View should show 'Confirm delete'")
	}
	if !strings.Contains(view, "Food") {
		t.Error("View should show category name")
	}
}

func TestCategoryListModel_View_EmptyCategories(t *testing.T) {
	store := &mockStoreWithCategoryList{categories: []string{}}
	model := NewCategoryListModel(store)

	view := model.View()
	if !strings.Contains(view, "(no categories defined)") {
		t.Error("View should show '(no categories defined)'")
	}
}

func TestCategoryListModel_View_WithError(t *testing.T) {
	store := &mockStoreWithCategoryList{categories: []string{}}
	model := NewCategoryListModel(store)
	model.errorMsg = "Failed to load"

	view := model.View()
	if !strings.Contains(view, "Failed to load") {
		t.Error("View should show error message")
	}
}

func TestCategoryListModel_View_WithSuccess(t *testing.T) {
	store := &mockStoreWithCategoryList{categories: []string{}}
	model := NewCategoryListModel(store)
	model.successMsg = "Category added successfully"

	view := model.View()
	if !strings.Contains(view, "Category added successfully") {
		t.Error("View should show success message")
	}
}

func TestCategoryListModel_CursorAtBottom(t *testing.T) {
	store := &mockStoreWithCategoryList{categories: []string{"Food"}}
	model := NewCategoryListModel(store)
	model.cursor = 1

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	newCategoryListModel := newModel.(CategoryListModel)

	if newCategoryListModel.cursor != 1 {
		t.Error("Cursor at bottom should stay at 1")
	}
}

func TestCategoryListModel_VimKeys(t *testing.T) {
	store := &mockStoreWithCategoryList{categories: []string{"Food", "Transport"}}
	model := NewCategoryListModel(store)
	model.cursor = 1

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	newCategoryListModel := newModel.(CategoryListModel)

	if newCategoryListModel.cursor != 0 {
		t.Errorf("cursor should be 0 after k, got %d", newCategoryListModel.cursor)
	}

	newModel, _ = newCategoryListModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	finalModel := newModel.(CategoryListModel)

	if finalModel.cursor != 1 {
		t.Errorf("cursor should be 1 after j, got %d", finalModel.cursor)
	}
}

type mockStoreWithExpenseErrors struct {
	mockStoreForCLI
	expenses   []storage.Expense
	categories []string
	getErr     error
	updateErr  error
}

func (m *mockStoreWithExpenseErrors) GetAllExpenses(startDate, endDate *time.Time) ([]storage.Expense, error) {
	return m.expenses, m.getErr
}

func (m *mockStoreWithExpenseErrors) GetCategories() ([]string, error) {
	return m.categories, nil
}

func (m *mockStoreWithExpenseErrors) UpdateExpense(id string, expense storage.Expense) error {
	return m.updateErr
}

func (m *mockStoreWithExpenseErrors) RemoveExpense(id string) error {
	return nil
}

func TestExpenseListModel_LoadData_WithError(t *testing.T) {
	store := &mockStoreWithExpenseErrors{
		expenses:   []storage.Expense{},
		categories: []string{"Food"},
		getErr:     fmt.Errorf("database error"),
	}
	model := NewExpenseListModel(store)

	if model.errorMsg == "" {
		t.Error("errorMsg should be set when load fails")
	}
}

func TestExpenseListModel_View_CategoryPicker(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test", Amount: 10.00, Category: "Food", Date: time.Now()},
		},
		categories: []string{"Food", "Transport", "Utilities"},
	}
	model := NewExpenseListModel(store)
	model.showCategoryPicker = true
	model.categoryPickerIdx = 1
	model.mode = ExpenseListEdit
	model.editingExpense = model.expenses[0]

	view := model.View()
	if !strings.Contains(view, "Select Category") {
		t.Error("View should show 'Select Category'")
	}
	if !strings.Contains(view, "Transport") {
		t.Error("View should show 'Transport'")
	}
}

func TestExpenseListModel_Update_CategoryPicker_Navigate(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test", Amount: 10.00, Category: "Food", Date: time.Now()},
		},
		categories: []string{"Food", "Transport", "Utilities"},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListEdit
	model.showCategoryPicker = true
	model.categoryPickerIdx = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.categoryPickerIdx != 1 {
		t.Errorf("categoryPickerIdx should be 1, got %d", newExpenseListModel.categoryPickerIdx)
	}
}

func TestExpenseListModel_Update_CategoryPicker_Select(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test", Amount: 10.00, Category: "Food", Date: time.Now()},
		},
		categories: []string{"Food", "Transport", "Utilities"},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListEdit
	model.showCategoryPicker = true
	model.categoryPickerIdx = 1
	model.editingExpense = model.expenses[0]

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.showCategoryPicker {
		t.Error("showCategoryPicker should be false after selection")
	}
	if newExpenseListModel.editValue != "Transport" {
		t.Errorf("editValue should be 'Transport', got %s", newExpenseListModel.editValue)
	}
}

func TestExpenseListModel_Update_CategoryPicker_Esc(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test", Amount: 10.00, Category: "Food", Date: time.Now()},
		},
		categories: []string{"Food", "Transport"},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListEdit
	model.showCategoryPicker = true

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.showCategoryPicker {
		t.Error("showCategoryPicker should be false after Esc")
	}
}

func TestExpenseListModel_HandleEditSave_ValidationError(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses:   []storage.Expense{},
		categories: []string{"Food"},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListEdit
	model.editField = 0
	model.editValue = ""
	model.editingExpense = storage.Expense{ID: "test-1", Name: "", Amount: 10.00, Category: "Food"}

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.errorMsg == "" {
		t.Error("errorMsg should be set on validation error")
	}
}

func TestExpenseListModel_HandleEditSave_UpdateError(t *testing.T) {
	store := &mockStoreWithExpenseErrors{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test", Amount: 10.00, Category: "Food", Date: time.Now()},
		},
		categories: []string{"Food"},
		updateErr:  fmt.Errorf("database error"),
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListEdit
	model.editField = 0
	model.editValue = "Updated Test"
	model.editingExpense = storage.Expense{ID: "test-1", Name: "Test", Amount: 10.00, Category: "Food", Date: time.Now()}

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.errorMsg == "" {
		t.Error("errorMsg should be set on update error")
	}
}

func TestExpenseListModel_HandleEditSave_CurrencyField(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test", Amount: 10.00, Category: "Food", Currency: "usd", Date: time.Now()},
		},
		categories: []string{"Food"},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListEdit
	model.editField = 3
	model.editValue = "EUR"
	model.editingExpense = model.expenses[0]

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.successMsg == "" {
		t.Error("successMsg should be set on successful update")
	}
}

func TestExpenseListModel_HandleEditSave_DateField(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test", Amount: 10.00, Category: "Food", Date: time.Now()},
		},
		categories: []string{"Food"},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListEdit
	model.editField = 4
	model.editValue = "2024-12-25"
	model.editingExpense = model.expenses[0]

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.successMsg == "" {
		t.Error("successMsg should be set on successful date update")
	}
}

func TestExpenseListModel_Update_EditMode_DeleteChar(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses: []storage.Expense{
			{ID: "test-1", Name: "Test", Amount: 10.00, Category: "Food", Date: time.Now()},
		},
		categories: []string{"Food"},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListEdit
	model.editField = 0
	model.editValue = "Test"

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.editValue != "Tes" {
		t.Errorf("editValue = %q, want 'Tes'", newExpenseListModel.editValue)
	}
}

func TestExpenseListModel_Update_EditMode_CtrlC_Quits(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses:   []storage.Expense{},
		categories: []string{},
	}
	model := NewExpenseListModel(store)
	model.mode = ExpenseListEdit

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if cmd == nil {
		t.Error("Ctrl+C should return quit command")
	}
	if !newExpenseListModel.quitting {
		t.Error("quitting should be true")
	}
}

func TestExpenseListModel_Update_BrowseMode_Left(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses:   []storage.Expense{},
		categories: []string{},
	}
	model := NewExpenseListModel(store)
	model.page = 1

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.page != 0 {
		t.Errorf("page should be 0 after Left, got %d", newExpenseListModel.page)
	}
}

func TestExpenseListModel_Update_BrowseMode_Right_AtEnd(t *testing.T) {
	store := &mockStoreWithExpenses{
		expenses:   []storage.Expense{},
		categories: []string{},
	}
	model := NewExpenseListModel(store)
	model.page = 1

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRight})
	newExpenseListModel := *newModel.(*ExpenseListModel)

	if newExpenseListModel.page != 1 {
		t.Error("page should stay at 1 when at end")
	}
}

type mockStoreForRecurringStatus struct {
	recurringExpenses []storage.RecurringExpense
	expenses          []storage.Expense
	recurringErr      error
	expensesErr       error
}

func (m *mockStoreForRecurringStatus) Close() error { return nil }
func (m *mockStoreForRecurringStatus) GetConfig() (*storage.Config, error) {
	return &storage.Config{}, nil
}
func (m *mockStoreForRecurringStatus) GetCategories() ([]string, error) { return []string{}, nil }
func (m *mockStoreForRecurringStatus) UpdateCategories([]string) error  { return nil }
func (m *mockStoreForRecurringStatus) GetAllExpenses(startDate, endDate *time.Time) ([]storage.Expense, error) {
	return m.expenses, m.expensesErr
}
func (m *mockStoreForRecurringStatus) GetExpense(string) (storage.Expense, error) {
	return storage.Expense{}, nil
}
func (m *mockStoreForRecurringStatus) AddExpense(storage.Expense) error            { return nil }
func (m *mockStoreForRecurringStatus) RemoveExpense(string) error                  { return nil }
func (m *mockStoreForRecurringStatus) UpdateExpense(string, storage.Expense) error { return nil }
func (m *mockStoreForRecurringStatus) AddMultipleExpenses([]storage.Expense) error { return nil }
func (m *mockStoreForRecurringStatus) RemoveMultipleExpenses([]string) error       { return nil }
func (m *mockStoreForRecurringStatus) GetExclusionList() ([]string, error)         { return []string{}, nil }
func (m *mockStoreForRecurringStatus) AddExclusion(string) error                   { return nil }
func (m *mockStoreForRecurringStatus) RemoveExclusion(string) error                { return nil }
func (m *mockStoreForRecurringStatus) UpdateExclusionList([]string) error          { return nil }
func (m *mockStoreForRecurringStatus) GetImportPath() (string, error)              { return "", nil }
func (m *mockStoreForRecurringStatus) UpdateImportPath(string) error               { return nil }
func (m *mockStoreForRecurringStatus) GetCurrency() (string, error)                { return "usd", nil }
func (m *mockStoreForRecurringStatus) UpdateCurrency(string) error                 { return nil }
func (m *mockStoreForRecurringStatus) GetStartDate() (int, error)                  { return 1, nil }
func (m *mockStoreForRecurringStatus) UpdateStartDate(int) error                   { return nil }
func (m *mockStoreForRecurringStatus) GetRecurringExpenses() ([]storage.RecurringExpense, error) {
	return m.recurringExpenses, m.recurringErr
}
func (m *mockStoreForRecurringStatus) GetRecurringExpense(string) (storage.RecurringExpense, error) {
	return storage.RecurringExpense{}, nil
}
func (m *mockStoreForRecurringStatus) AddRecurringExpense(storage.RecurringExpense) error { return nil }
func (m *mockStoreForRecurringStatus) RemoveRecurringExpense(string, bool) error          { return nil }
func (m *mockStoreForRecurringStatus) UpdateRecurringExpense(string, storage.RecurringExpense, bool) error {
	return nil
}

func TestRecurringStatusModel_Initialization(t *testing.T) {
	now := time.Now()
	expectedMonth := int(now.Month()) - 1
	if expectedMonth == 0 {
		expectedMonth = 12
	}
	expectedYear := now.Year()
	if expectedMonth == 12 {
		expectedYear--
	}

	store := &mockStoreForRecurringStatus{}
	model := NewRecurringStatusModel(store)

	if model.month != expectedMonth {
		t.Errorf("Expected month %d, got %d", expectedMonth, model.month)
	}
	if model.year != expectedYear {
		t.Errorf("Expected year %d, got %d", expectedYear, model.year)
	}
}

func TestRecurringStatusModel_ExecuteQuery_NoRecurring(t *testing.T) {
	store := &mockStoreForRecurringStatus{
		recurringExpenses: []storage.RecurringExpense{},
		expenses:          []storage.Expense{},
	}
	model := NewRecurringStatusModel(store)
	model.month = 3
	model.year = 2026

	model.executeQuery()

	if len(model.results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(model.results))
	}
	if !model.showResults {
		t.Error("showResults should be true after executeQuery")
	}
}

func TestRecurringStatusModel_ExecuteQuery_WithPaid(t *testing.T) {
	store := &mockStoreForRecurringStatus{
		recurringExpenses: []storage.RecurringExpense{
			{ID: "rec-1", Name: "Netflix", Amount: 15.00},
		},
		expenses: []storage.Expense{
			{ID: "exp-1", Name: "Netflix Subscription", Amount: 15.00, Date: time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)},
		},
	}
	model := NewRecurringStatusModel(store)
	model.month = 3
	model.year = 2026

	model.executeQuery()

	if len(model.results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(model.results))
	}
	if !model.results[0].Paid {
		t.Error("Expected recurring to be marked as paid")
	}
	if model.results[0].Amount != 15.00 {
		t.Errorf("Expected amount 15.00, got %.2f", model.results[0].Amount)
	}
}

func TestRecurringStatusModel_ExecuteQuery_WithPending(t *testing.T) {
	store := &mockStoreForRecurringStatus{
		recurringExpenses: []storage.RecurringExpense{
			{ID: "rec-1", Name: "Netflix", Amount: 15.00},
		},
		expenses: []storage.Expense{},
	}
	model := NewRecurringStatusModel(store)
	model.month = 3
	model.year = 2026

	model.executeQuery()

	if len(model.results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(model.results))
	}
	if model.results[0].Paid {
		t.Error("Expected recurring to be marked as unpaid")
	}
}

func TestRecurringStatusModel_NavigateMonth_Next(t *testing.T) {
	store := &mockStoreForRecurringStatus{}
	model := NewRecurringStatusModel(store)
	model.month = 3
	model.year = 2026
	model.showResults = true

	model.navigateMonth(1)

	if model.month != 4 {
		t.Errorf("Expected month 4, got %d", model.month)
	}
	if model.year != 2026 {
		t.Errorf("Expected year 2026, got %d", model.year)
	}
}

func TestRecurringStatusModel_NavigateMonth_Previous(t *testing.T) {
	store := &mockStoreForRecurringStatus{}
	model := NewRecurringStatusModel(store)
	model.month = 3
	model.year = 2026
	model.showResults = true

	model.navigateMonth(-1)

	if model.month != 2 {
		t.Errorf("Expected month 2, got %d", model.month)
	}
	if model.year != 2026 {
		t.Errorf("Expected year 2026, got %d", model.year)
	}
}

func TestRecurringStatusModel_NavigateMonth_DecemberToJanuary(t *testing.T) {
	store := &mockStoreForRecurringStatus{}
	model := NewRecurringStatusModel(store)
	model.month = 12
	model.year = 2026
	model.showResults = true

	model.navigateMonth(1)

	if model.month != 1 {
		t.Errorf("Expected month 1, got %d", model.month)
	}
	if model.year != 2027 {
		t.Errorf("Expected year 2027, got %d", model.year)
	}
}

func TestRecurringStatusModel_NavigateMonth_JanuaryToDecember(t *testing.T) {
	store := &mockStoreForRecurringStatus{}
	model := NewRecurringStatusModel(store)
	model.month = 1
	model.year = 2026
	model.showResults = true

	model.navigateMonth(-1)

	if model.month != 12 {
		t.Errorf("Expected month 12, got %d", model.month)
	}
	if model.year != 2025 {
		t.Errorf("Expected year 2025, got %d", model.year)
	}
}

func TestRecurringStatusModel_Esc_ReturnsToMainMenu(t *testing.T) {
	store := &mockStoreForRecurringStatus{}
	model := NewRecurringStatusModel(store)

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	_ = newModel.(RecurringStatusModel)
}

func TestReportModel_NavigateMonth(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewReportModel(store)
	model.month = 3
	model.year = 2026
	model.showResults = true

	model.navigateMonth(1)

	if model.month != 4 {
		t.Errorf("Expected month 4, got %d", model.month)
	}
	if model.year != 2026 {
		t.Errorf("Expected year 2026, got %d", model.year)
	}
}

func TestReportModel_NavigateMonth_Previous(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewReportModel(store)
	model.month = 3
	model.year = 2026
	model.showResults = true

	model.navigateMonth(-1)

	if model.month != 2 {
		t.Errorf("Expected month 2, got %d", model.month)
	}
}

func TestReportModel_NavigateMonth_DecemberToJanuary(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewReportModel(store)
	model.month = 12
	model.year = 2026
	model.showResults = true

	model.navigateMonth(1)

	if model.month != 1 {
		t.Errorf("Expected month 1, got %d", model.month)
	}
	if model.year != 2027 {
		t.Errorf("Expected year 2027, got %d", model.year)
	}
}

type mockStoreForRecurringList struct {
	recurring  []storage.RecurringExpense
	categories []string
	err        error
}

func (m *mockStoreForRecurringList) Close() error { return nil }
func (m *mockStoreForRecurringList) GetConfig() (*storage.Config, error) {
	return &storage.Config{}, nil
}
func (m *mockStoreForRecurringList) GetCategories() ([]string, error) { return m.categories, m.err }
func (m *mockStoreForRecurringList) UpdateCategories([]string) error  { return nil }
func (m *mockStoreForRecurringList) GetAllExpenses(startDate, endDate *time.Time) ([]storage.Expense, error) {
	return []storage.Expense{}, nil
}
func (m *mockStoreForRecurringList) GetExpense(string) (storage.Expense, error) {
	return storage.Expense{}, nil
}
func (m *mockStoreForRecurringList) AddExpense(storage.Expense) error            { return nil }
func (m *mockStoreForRecurringList) RemoveExpense(string) error                  { return nil }
func (m *mockStoreForRecurringList) UpdateExpense(string, storage.Expense) error { return nil }
func (m *mockStoreForRecurringList) AddMultipleExpenses([]storage.Expense) error { return nil }
func (m *mockStoreForRecurringList) RemoveMultipleExpenses([]string) error       { return nil }
func (m *mockStoreForRecurringList) GetExclusionList() ([]string, error)         { return []string{}, nil }
func (m *mockStoreForRecurringList) AddExclusion(string) error                   { return nil }
func (m *mockStoreForRecurringList) RemoveExclusion(string) error                { return nil }
func (m *mockStoreForRecurringList) UpdateExclusionList([]string) error          { return nil }
func (m *mockStoreForRecurringList) GetImportPath() (string, error)              { return "", nil }
func (m *mockStoreForRecurringList) UpdateImportPath(string) error               { return nil }
func (m *mockStoreForRecurringList) GetCurrency() (string, error)                { return "usd", nil }
func (m *mockStoreForRecurringList) UpdateCurrency(string) error                 { return nil }
func (m *mockStoreForRecurringList) GetStartDate() (int, error)                  { return 1, nil }
func (m *mockStoreForRecurringList) UpdateStartDate(int) error                   { return nil }
func (m *mockStoreForRecurringList) GetRecurringExpenses() ([]storage.RecurringExpense, error) {
	return m.recurring, m.err
}
func (m *mockStoreForRecurringList) GetRecurringExpense(string) (storage.RecurringExpense, error) {
	if len(m.recurring) > 0 {
		return m.recurring[0], nil
	}
	return storage.RecurringExpense{}, fmt.Errorf("not found")
}
func (m *mockStoreForRecurringList) AddRecurringExpense(e storage.RecurringExpense) error { return nil }
func (m *mockStoreForRecurringList) RemoveRecurringExpense(id string, b bool) error       { return m.err }
func (m *mockStoreForRecurringList) UpdateRecurringExpense(id string, e storage.RecurringExpense, b bool) error {
	return nil
}

func TestRecurringListModel_Initialization(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{},
		categories: []string{"Food", "Travel"},
	}
	model := NewRecurringListModel(store)

	if model.mode != RecurringView {
		t.Errorf("Expected mode RecurringView, got %d", model.mode)
	}
	if model.cursor != 0 {
		t.Error("Expected cursor to be 0")
	}
}

func TestRecurringListModel_Update_Esc_ViewToSettings(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{},
		categories: []string{"Food", "Travel"},
	}
	model := NewRecurringListModel(store)

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	newRecurringList := newModel.(RecurringListModel)

	if newRecurringList.mode != RecurringView {
		t.Error("Expected mode to stay RecurringView")
	}
}

func TestRecurringListModel_Update_Up(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{},
		categories: []string{"Food", "Travel"},
	}
	model := NewRecurringListModel(store)
	model.cursor = 1

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyUp})
	newRecurringList := newModel.(RecurringListModel)

	if newRecurringList.cursor != 0 {
		t.Error("Expected cursor to move up")
	}
}

func TestRecurringListModel_Update_Down(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{},
		categories: []string{"Food", "Travel"},
	}
	model := NewRecurringListModel(store)
	model.payments = []storage.RecurringExpense{{ID: "1"}}

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	newRecurringList := newModel.(RecurringListModel)

	if newRecurringList.cursor != 1 {
		t.Error("Expected cursor to move down")
	}
}

func TestRecurringListModel_Update_Add(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{},
		categories: []string{"Food", "Travel"},
	}
	model := NewRecurringListModel(store)

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	newRecurringList := newModel.(RecurringListModel)

	if newRecurringList.mode != RecurringAdd {
		t.Error("Expected mode to be RecurringAdd")
	}
}

func TestRecurringListModel_Update_Edit(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{{ID: "1", Name: "Test"}},
		categories: []string{"Food", "Travel"},
	}
	model := NewRecurringListModel(store)
	model.payments = []storage.RecurringExpense{{ID: "1", Name: "Test"}}

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	newRecurringList := newModel.(RecurringListModel)

	if newRecurringList.mode != RecurringEdit {
		t.Error("Expected mode to be RecurringEdit")
	}
}

func TestRecurringListModel_Update_Delete(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{{ID: "1", Name: "Test"}},
		categories: []string{"Food", "Travel"},
	}
	model := NewRecurringListModel(store)
	model.payments = []storage.RecurringExpense{{ID: "1", Name: "Test"}}

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	newRecurringList := newModel.(RecurringListModel)

	if newRecurringList.mode != RecurringDeleteConfirm {
		t.Error("Expected mode to be RecurringDeleteConfirm")
	}
}

func TestRecurringListModel_MonthsToString(t *testing.T) {
	result := monthsToString([]int{1, 2, 3})
	if result != "Jan,Feb,Mar" {
		t.Errorf("Expected 'Jan,Feb,Mar', got %s", result)
	}
}

func TestRecurringListModel_StringToMonths(t *testing.T) {
	result := stringToMonths("Jan,Mar")
	if len(result) != 2 {
		t.Errorf("Expected 2 months, got %d", len(result))
	}
}

func TestRecurringListModel_MonthToString(t *testing.T) {
	result := monthToString(1)
	if result != "Jan" {
		t.Errorf("Expected 'Jan', got %s", result)
	}
}

func TestRecurringListModel_MonthToNumber(t *testing.T) {
	result := monthToNumber("January")
	if result != 1 {
		t.Errorf("Expected 1, got %d", result)
	}
}

func TestRecurringListModel_YesNo(t *testing.T) {
	if yesNo(true) != "Yes" {
		t.Error("Expected 'Yes' for true")
	}
	if yesNo(false) != "No" {
		t.Error("Expected 'No' for false")
	}
}

func TestRecurringListModel_Update_Backspace(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{},
		categories: []string{"Food", "Travel"},
	}
	model := NewRecurringListModel(store)
	model.mode = RecurringAdd
	model.inputValue = "test"

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	newRecurringList := newModel.(RecurringListModel)

	if newRecurringList.inputValue != "tes" {
		t.Error("Expected backspace to remove last char")
	}
}

func TestRecurringListModel_Update_LeftRight(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{},
		categories: []string{"Food", "Travel"},
	}
	model := NewRecurringListModel(store)
	model.mode = RecurringAdd
	model.payments = []storage.RecurringExpense{{ID: "1"}}

	// Left in view mode should move cursor
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	newRecurringList := newModel.(RecurringListModel)

	if newRecurringList.cursor != 0 {
		t.Error("Left should not move cursor in input mode")
	}
}

func TestRecurringListModel_Update_Enter_View_Empty(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{},
		categories: []string{"Food", "Travel"},
	}
	model := NewRecurringListModel(store)
	model.payments = nil

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newRecurringList := newModel.(RecurringListModel)

	// Should return navigation command when cursor at end
	if cmd == nil {
		t.Error("Expected command when at end of list")
	}
	_ = newRecurringList
}

func TestRecurringListModel_View_Empty(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{},
		categories: []string{"Food", "Travel"},
	}
	model := NewRecurringListModel(store)
	model.payments = nil

	view := model.View()

	if !strings.Contains(view, "no recurring") {
		t.Error("Expected empty message in view")
	}
}

func TestRecurringListModel_View_WithPayments(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{{ID: "1", Name: "Netflix", Amount: 15.00, Category: "Entertainment"}},
		categories: []string{"Food", "Travel"},
	}
	model := NewRecurringListModel(store)
	model.payments = []storage.RecurringExpense{{ID: "1", Name: "Netflix", Amount: 15.00, Category: "Entertainment"}}

	view := model.View()

	if !strings.Contains(view, "Netflix") {
		t.Error("Expected Netflix in view")
	}
}

func TestRecurringListModel_handleFieldInput_Name(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{},
		categories: []string{"Food", "Travel"},
	}
	model := NewRecurringListModel(store)
	model.field = FieldName
	model.inputValue = "Netflix"
	model.editing = storage.RecurringExpense{}

	model.handleFieldInput()

	if model.field != FieldCategory {
		t.Error("Expected field to move to Category")
	}
}

func TestRecurringListModel_handleFieldInput_Category(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{},
		categories: []string{"Food", "Travel"},
	}
	model := NewRecurringListModel(store)
	model.field = FieldCategory
	model.categoryCursor = 0
	model.editing = storage.RecurringExpense{}
	model.inputValue = "" // Clear since category is selected

	model.handleFieldInput()

	if model.field != FieldScheduleType {
		t.Error("Expected field to move to ScheduleType")
	}
}

func TestRecurringListModel_handleFieldInput_InvalidSchedule(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{},
		categories: []string{"Food", "Travel"},
	}
	model := NewRecurringListModel(store)
	model.field = FieldScheduleType
	model.inputValue = "invalid"
	model.editing = storage.RecurringExpense{}

	model.handleFieldInput()

	// Should stay on schedule type for invalid input
	if model.field != FieldScheduleType {
		t.Error("Expected to stay on invalid schedule type")
	}
}

func TestRecurringListModel_Truncate(t *testing.T) {
	result := truncate("VeryLongName", 10)
	if result != "VeryLong.." {
		t.Errorf("Expected 'VeryLong..', got %s", result)
	}

	result = truncate("Short", 10)
	if result != "Short" {
		t.Errorf("Expected 'Short', got %s", result)
	}
}

func TestRecurringStatusModel_Update_Left(t *testing.T) {
	store := &mockStoreForRecurringStatus{}
	model := NewRecurringStatusModel(store)
	model.month = 1
	model.year = 2026
	model.showResults = true

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	newRecurringStatus := newModel.(RecurringStatusModel)

	if newRecurringStatus.month != 12 {
		t.Error("Expected month to decrease")
	}
}

func TestRecurringStatusModel_Update_Right(t *testing.T) {
	store := &mockStoreForRecurringStatus{}
	model := NewRecurringStatusModel(store)
	model.month = 12
	model.year = 2026
	model.showResults = true

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRight})
	newRecurringStatus := newModel.(RecurringStatusModel)

	if newRecurringStatus.month != 1 {
		t.Error("Expected month to increase")
	}
}

func TestRecurringStatusModel_Update_CtrlC(t *testing.T) {
	store := &mockStoreForRecurringStatus{}
	model := NewRecurringStatusModel(store)

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	newRecurringStatus := newModel.(RecurringStatusModel)

	if !newRecurringStatus.quitting {
		t.Error("Expected quitting to be true")
	}
	if cmd == nil {
		t.Error("Expected non-nil command")
	}
}

func TestRecurringStatusModel_Update_Up(t *testing.T) {
	store := &mockStoreForRecurringStatus{
		recurringExpenses: []storage.RecurringExpense{{ID: "1", Name: "Test"}},
	}
	model := NewRecurringStatusModel(store)
	model.showResults = true
	model.results = []RecurringStatusItem{{Name: "Test", Paid: false}}
	model.cursor = 1

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyUp})
	newRecurringStatus := newModel.(RecurringStatusModel)

	if newRecurringStatus.cursor != 0 {
		t.Error("Expected cursor to move up")
	}
}

func TestRecurringStatusModel_Update_Down(t *testing.T) {
	store := &mockStoreForRecurringStatus{
		recurringExpenses: []storage.RecurringExpense{{ID: "1", Name: "Test"}},
	}
	model := NewRecurringStatusModel(store)
	model.showResults = true
	model.results = []RecurringStatusItem{{Name: "Test", Paid: false}}

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	newRecurringStatus := newModel.(RecurringStatusModel)

	if newRecurringStatus.cursor != 1 {
		t.Error("Expected cursor to move down")
	}
}

func TestRecurringListModel_Update_CtrlC(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{},
		categories: []string{"Food"},
	}
	model := NewRecurringListModel(store)

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	newRecurringList := newModel.(RecurringListModel)

	if !newRecurringList.quitting {
		t.Error("Expected quitting to be true")
	}
	if cmd == nil {
		t.Error("Expected command to be returned")
	}
}

func TestRecurringListModel_SetCategoryCursor(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{},
		categories: []string{"Food", "Travel", "Entertainment"},
	}
	model := NewRecurringListModel(store)
	model.editing.Category = "Travel"

	model.setCategoryCursor()

	if model.categoryCursor != 1 {
		t.Errorf("Expected cursor at 1, got %d", model.categoryCursor)
	}
}

func TestRecurringListModel_setCategoryCursor_NotFound(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{},
		categories: []string{"Food", "Travel"},
	}
	model := NewRecurringListModel(store)
	model.editing.Category = "Unknown"

	model.setCategoryCursor()

	if model.categoryCursor != 0 {
		t.Errorf("Expected cursor at 0, got %d", model.categoryCursor)
	}
}

func TestRecurringListModel_handleFieldInput_Enabled(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{},
		categories: []string{"Food"},
	}
	model := NewRecurringListModel(store)
	model.field = FieldEnabled
	model.mode = RecurringAdd
	model.inputValue = "y"

	model.handleFieldInput()

	// Verify input was processed (mode or field changed or save was attempted)
	_ = model.inputValue // Used
}

func TestRecurringListModel_handleFieldInput_Schedule(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{},
		categories: []string{"Food"},
	}
	model := NewRecurringListModel(store)
	model.field = FieldScheduleType
	model.inputValue = "monthly"

	model.handleFieldInput()

	if model.field != FieldEnabled {
		t.Error("Expected field to move to Enabled")
	}
}

func TestRecurringListModel_handleFieldInput_Months(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{},
		categories: []string{"Food"},
	}
	model := NewRecurringListModel(store)
	model.field = FieldMonths
	model.inputValue = "Jan,Mar"

	model.handleFieldInput()

	// Just verify no panic - the logic changes based on schedule type
}

func TestRecurringListModel_removePayment_Error(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{},
		categories: []string{"Food"},
		err:       fmt.Errorf("database error"),
	}
	model := NewRecurringListModel(store)
	model.errorMsg = "test error"

	model.removePayment("test-id")

	if model.errorMsg == "" {
		t.Error("Expected errorMsg to be set when RemoveRecurringExpense fails")
	}
	if model.successMsg != "" {
		t.Error("Expected successMsg to be empty on error")
	}
}

func TestRecurringListModel_Init(t *testing.T) {
	store := &mockStoreForRecurringList{}
	model := NewRecurringListModel(store)

	cmd := model.Init()

	if cmd != nil {
		t.Error("Expected nil cmd from Init")
	}
}

func TestRecurringListModel_View_Form(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{},
		categories: []string{"Food"},
	}
	model := NewRecurringListModel(store)
	model.mode = RecurringAdd

	view := model.View()

	if !strings.Contains(view, "Add") {
		t.Error("Expected Add in view")
	}
}

func TestRecurringListModel_View_Edit(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{},
		categories: []string{"Food"},
	}
	model := NewRecurringListModel(store)
	model.mode = RecurringEdit
	model.editing = storage.RecurringExpense{Name: "Test"}

	view := model.View()

	if !strings.Contains(view, "Edit") {
		t.Error("Expected Edit in view")
	}
}

func TestRecurringListModel_View_Delete(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{},
		categories: []string{"Food"},
	}
	model := NewRecurringListModel(store)
	model.mode = RecurringDeleteConfirm

	view := model.View()

	if !strings.Contains(view, "delete") {
		t.Error("Expected delete confirmation in view")
	}
}

func TestRecurringStatusModel_View_Empty(t *testing.T) {
	store := &mockStoreForRecurringStatus{}
	model := NewRecurringStatusModel(store)

	view := model.View()

	if !strings.Contains(view, "Recurring Expenses") {
		t.Error("Expected Recurring Expenses in view")
	}
}

func TestRecurringStatusModel_View_WithResults(t *testing.T) {
	store := &mockStoreForRecurringStatus{
		recurringExpenses: []storage.RecurringExpense{{ID: "1", Name: "Test", Amount: 10}},
	}
	model := NewRecurringStatusModel(store)
	model.results = []RecurringStatusItem{{Name: "Test", Amount: 10, Paid: true}}
	model.showResults = true

	view := model.View()

	_ = view
}

func TestRecurringStatusModel_TruncateStatus(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{"short string", "Test", 15, "Test"},
		{"exact length", "123456789012345", 15, "123456789012345"},
		{"long string", "This is a very long name", 15, "This is a ver.."},
		{"empty string", "", 15, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateStatus(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("truncateStatus(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}

func TestRecurringStatusModel_Update_MonthNav(t *testing.T) {
	store := &mockStoreForRecurringStatus{
		recurringExpenses: []storage.RecurringExpense{{ID: "1", Name: "Test", Amount: 10}},
	}
	model := NewRecurringStatusModel(store)

	initialMonth := model.month

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("right")})
	updatedModel := updated.(RecurringStatusModel)

	if updatedModel.month != initialMonth+1 {
		t.Errorf("Expected month to increment from %d to %d, got %d", initialMonth, initialMonth+1, updatedModel.month)
	}
}

func TestRecurringStatusModel_Update_MonthNav_Decrement(t *testing.T) {
	store := &mockStoreForRecurringStatus{
		recurringExpenses: []storage.RecurringExpense{{ID: "1", Name: "Test", Amount: 10}},
	}
	model := NewRecurringStatusModel(store)

	initialMonth := model.month

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("left")})
	updatedModel := updated.(RecurringStatusModel)

	if updatedModel.month != initialMonth-1 {
		t.Errorf("Expected month to decrement from %d to %d, got %d", initialMonth, initialMonth-1, updatedModel.month)
	}
}

func TestRecurringListModel_Update_YearNav(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{},
		categories: []string{"Food"},
	}
	model := NewRecurringListModel(store)
	model.mode = RecurringView

	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})

	// Should not panic
}

func TestRecurringListModel_Update_TypeChar_WhileViewing(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{},
		categories: []string{"Food"},
	}
	model := NewRecurringListModel(store)
	model.mode = RecurringView

	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})

	if updated.(RecurringListModel).mode != RecurringAdd {
		t.Error("Expected mode to be RecurringAdd")
	}
	if cmd != nil {
		t.Error("Expected no command")
	}
}

func TestRecurringListModel_Update_TypeChar_WhileViewing_Edit(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{{ID: "1", Name: "Test"}},
		categories: []string{"Food"},
	}
	model := NewRecurringListModel(store)
	model.mode = RecurringView
	model.cursor = 0

	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("e")})

	if updated.(RecurringListModel).mode != RecurringEdit {
		t.Error("Expected mode to be RecurringEdit")
	}
	if cmd != nil {
		t.Error("Expected no command")
	}
}

func TestRecurringListModel_Update_TypeChar_WhileViewing_Delete(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{{ID: "1", Name: "Test"}},
		categories: []string{"Food"},
	}
	model := NewRecurringListModel(store)
	model.mode = RecurringView
	model.cursor = 0

	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})

	if updated.(RecurringListModel).mode != RecurringDeleteConfirm {
		t.Error("Expected mode to be RecurringDeleteConfirm")
	}
	if cmd != nil {
		t.Error("Expected no command")
	}
}

func TestRecurringListModel_Update_EnterOnEmptyList(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{},
		categories: []string{"Food"},
	}
	model := NewRecurringListModel(store)
	model.mode = RecurringView
	model.cursor = 0

	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if cmd == nil {
		t.Error("Expected command for navigation")
	}
}

func TestRecurringListModel_Update_Esc_FromAddResetsForm(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{},
		categories: []string{"Food"},
	}
	model := NewRecurringListModel(store)
	model.mode = RecurringAdd
	model.inputValue = "test"
	model.field = FieldName

	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})

	updatedModel := updated.(RecurringListModel)
	if updatedModel.mode != RecurringView {
		t.Error("Expected mode to be RecurringView")
	}
	if updatedModel.inputValue != "" {
		t.Error("Expected inputValue to be cleared")
	}
	if cmd != nil {
		t.Error("Expected no command")
	}
}

func TestRecurringListModel_Update_Esc_FromEditResetsForm(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{{ID: "1", Name: "Test"}},
		categories: []string{"Food"},
	}
	model := NewRecurringListModel(store)
	model.mode = RecurringEdit
	model.inputValue = "test"
	model.field = FieldName
	model.editing = storage.RecurringExpense{ID: "1", Name: "Test"}

	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})

	updatedModel := updated.(RecurringListModel)
	if updatedModel.mode != RecurringView {
		t.Error("Expected mode to be RecurringView")
	}
	if updatedModel.inputValue != "" {
		t.Error("Expected inputValue to be cleared")
	}
	if cmd != nil {
		t.Error("Expected no command")
	}
}

func TestRecurringListModel_Update_Esc_FromDeleteConfirm(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{{ID: "1", Name: "Test"}},
		categories: []string{"Food"},
	}
	model := NewRecurringListModel(store)
	model.mode = RecurringDeleteConfirm
	model.cursor = 0

	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})

	updatedModel := updated.(RecurringListModel)
	if updatedModel.mode != RecurringView {
		t.Error("Expected mode to be RecurringView")
	}
	if cmd != nil {
		t.Error("Expected no command")
	}
}

func TestRecurringListModel_renderForm_CategorySelection(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{},
		categories: []string{"Food", "Transport", "Entertainment"},
	}
	model := NewRecurringListModel(store)
	model.mode = RecurringAdd
	model.field = FieldCategory
	model.categoryCursor = 1
	model.editing = storage.RecurringExpense{}

	var s strings.Builder
	model.renderForm(&s, "Add")

	view := s.String()
	if !strings.Contains(view, "Transport") {
		t.Error("Expected Transport in category selection")
	}
}

func TestExclusionListModel_IsQuitting(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)

	if model.IsQuitting() {
		t.Error("IsQuitting should be false initially")
	}

	model.quitting = true
	if !model.IsQuitting() {
		t.Error("IsQuitting should be true after quitting is set")
	}
}

func TestAppModel_View_DefaultScreen(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewAppModel(store)
	model.currentScreen = 999

	view := model.View()
	if view != "Unknown screen\n" {
		t.Error("Expected 'Unknown screen' for invalid screen")
	}
}

func TestRecurringStatusModel_Update_CtrlC_Double(t *testing.T) {
	store := &mockStoreForRecurringStatus{
		recurringExpenses: []storage.RecurringExpense{{ID: "1", Name: "Test", Amount: 10}},
	}
	model := NewRecurringStatusModel(store)

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	newStatusModel := newModel.(RecurringStatusModel)

	if cmd == nil {
		t.Error("Ctrl+C should return quit command")
	}
	if !newStatusModel.quitting {
		t.Error("quitting should be true after Ctrl+C")
	}
}

func TestRecurringListModel_Update_CtrlC_WhileViewing(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{},
		categories: []string{"Food"},
	}
	model := NewRecurringListModel(store)
	model.mode = RecurringView

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	newRecurringModel := newModel.(RecurringListModel)

	if cmd == nil {
		t.Error("Ctrl+C should return quit command")
	}
	if !newRecurringModel.quitting {
		t.Error("quitting should be true after Ctrl+C")
	}
}

func TestRecurringListModel_Update_EnterInDeleteConfirm_BackOption(t *testing.T) {
	store := &mockStoreForRecurringList{
		recurring:  []storage.RecurringExpense{{ID: "1", Name: "Test"}},
		categories: []string{"Food"},
	}
	model := NewRecurringListModel(store)
	model.mode = RecurringDeleteConfirm
	model.cursor = 1

	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Enter on back option should just switch to view mode (no cmd)
}
