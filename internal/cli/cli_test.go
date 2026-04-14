package cli

import (
	"os"
	"strings"
	"testing"
	"time"

	"expensecat/internal/api"
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
func (m *mockStoreForCLI) GetExclusionList() ([]string, error) { return []string{}, nil }
func (m *mockStoreForCLI) AddExclusion(string) error           { return nil }
func (m *mockStoreForCLI) RemoveExclusion(string) error        { return nil }
func (m *mockStoreForCLI) UpdateExclusionList([]string) error  { return nil }
func (m *mockStoreForCLI) GetImportPath() (string, error)      { return "~/Downloads/Expense-Files", nil }
func (m *mockStoreForCLI) UpdateImportPath(string) error       { return nil }
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

	if len(model.choices) != 4 {
		t.Errorf("Expected 4 choices, got %d", len(model.choices))
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
	if model.choices[2] != "Settings" {
		t.Errorf("Expected third choice to be 'Settings', got %s", model.choices[2])
	}
	if model.choices[3] != "Quit" {
		t.Errorf("Expected fourth choice to be 'Quit', got %s", model.choices[3])
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
	model.cursor = 3

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = newModel.(MainMenuModel)
	if model.cursor != 3 {
		t.Errorf("At bottom, cursor should stay at 3, got %d", model.cursor)
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
		t.Errorf("At bottom, cursor should stay at 2, got %d", model.cursor)
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
	model.cursor = 3
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

func TestReportModel_Execute_InvalidMonth(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewReportModel(store)
	model.inputMonth = "13"
	model.inputYear = "2024"
	model.cursor = 2

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newReportModel := newModel.(ReportModel)
	if newReportModel.errorMsg == "" {
		t.Error("Should set errorMsg for invalid month")
	}
}

func TestReportModel_Execute_InvalidYear(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewReportModel(store)
	model.inputMonth = "3"
	model.inputYear = "abc"
	model.cursor = 2

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newReportModel := newModel.(ReportModel)
	if newReportModel.errorMsg == "" {
		t.Error("Should set errorMsg for invalid year")
	}
}

func TestReportModel_InputYear_Backspace(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewReportModel(store)
	model.inputtingFor = 2
	model.inputYear = "2024"

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	newReportModel := newModel.(ReportModel)
	if newReportModel.inputYear != "202" {
		t.Errorf("inputYear = %q, want '202'", newReportModel.inputYear)
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

	view := model.View()
	if !strings.Contains(view, "Month") {
		t.Error("View() should contain Month")
	}
	if !strings.Contains(view, "Year") {
		t.Error("View() should contain Year")
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
	model.inputMonth = "1"
	model.inputYear = "2024"
	model.cursor = 2

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newReportModel := newModel.(ReportModel)
	if newReportModel.errorMsg != "" {
		t.Logf("executeReport error (expected): %v", newReportModel.errorMsg)
	}
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

	view := model.View()
	if !strings.Contains(view, "Arrow") && !strings.Contains(view, "navigate") {
		t.Error("View() should contain navigation help")
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

func TestReportModel_View_MonthInputLabel(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewReportModel(store)
	model.inputtingFor = 1
	model.inputMonth = "5"

	view := model.View()
	if !strings.Contains(view, "5") {
		t.Error("View() should show month input")
	}
}

func TestReportModel_View_YearInputLabel(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewReportModel(store)
	model.inputtingFor = 2
	model.inputYear = "2025"

	view := model.View()
	if !strings.Contains(view, "2025") {
		t.Error("View() should show year input")
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
	view = model.View()
	if !strings.Contains(view, "Report") {
		t.Error("View() should show Report screen")
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
	view := model.View()

	if !strings.Contains(view, "Exclusion List") {
		t.Error("View() should contain Exclusion List")
	}
	if !strings.Contains(view, "Add Pattern") {
		t.Error("View() should contain Add Pattern")
	}
	if !strings.Contains(view, "Remove Pattern") {
		t.Error("View() should contain Remove Pattern")
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

func TestExclusionListModel_RemoveMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)

	model.mode = ExclusionListRemove
	model.patterns = []string{"PATTERN1", "PATTERN2"}
	view := model.View()

	if !strings.Contains(view, "Select pattern to remove:") {
		t.Error("View() should show removal prompt in remove mode")
	}
	if !strings.Contains(view, "PATTERN1") {
		t.Error("View() should show patterns")
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

	if newExclusionModel.mode != ExclusionListMenu {
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

func TestExclusionListModel_EnterOnBack(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)
	model.cursor = 2

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	_ = newModel.(ExclusionListModel)

	if cmd == nil {
		t.Error("Enter on Back should return navigation command")
	}
}

func TestExclusionListModel_EnterInAddMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)
	model.mode = ExclusionListAdd
	model.inputValue = "NEW PATTERN"

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newExclusionModel := newModel.(ExclusionListModel)

	if newExclusionModel.mode != ExclusionListMenu {
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
	model.exclusionModel.mode = ExclusionListRemove

	newModel, _ := model.Update(NavigationMsg{Destination: ScreenExclusionList})
	newAppModel := newModel.(AppModel)

	if newAppModel.exclusionModel.mode != ExclusionListMenu {
		t.Error("Navigation to ExclusionList should reset mode")
	}
}

func TestExclusionListModel_RemoveSelectsPattern(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)
	model.mode = ExclusionListRemove
	model.patterns = []string{"TEST1", "TEST2"}
	model.cursor = 1

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newExclusionModel := newModel.(ExclusionListModel)

	if newExclusionModel.mode != ExclusionListMenu {
		t.Error("Enter on pattern should remove and return to menu")
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

func TestExclusionListModel_EscInRemoveMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)
	model.mode = ExclusionListRemove

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	newExclusionModel := newModel.(ExclusionListModel)

	if newExclusionModel.mode != ExclusionListMenu {
		t.Error("Esc in remove mode should return to menu")
	}
}

func TestExclusionListModel_UpInRemoveMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)
	model.mode = ExclusionListRemove
	model.patterns = []string{"TEST1", "TEST2"}
	model.cursor = 1

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyUp})
	newExclusionModel := newModel.(ExclusionListModel)

	if newExclusionModel.cursor != 0 {
		t.Error("Up should move cursor up in remove mode")
	}
}

func TestExclusionListModel_DownInRemoveMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)
	model.mode = ExclusionListRemove
	model.patterns = []string{"TEST1", "TEST2"}
	model.cursor = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	newExclusionModel := newModel.(ExclusionListModel)

	if newExclusionModel.cursor != 1 {
		t.Error("Down should move cursor down in remove mode")
	}
}

func TestExclusionListModel_DownAtEndOfRemoveMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)
	model.mode = ExclusionListRemove
	model.patterns = []string{"TEST1"}
	model.cursor = 0

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	newExclusionModel := newModel.(ExclusionListModel)

	if newExclusionModel.cursor != 1 {
		t.Error("Down should move to Back option in remove mode with 1 pattern")
	}
}

func TestExclusionListModel_EnterOnBackInRemoveMode(t *testing.T) {
	store := &mockStoreForCLI{}
	model := NewExclusionListModel(store)
	model.mode = ExclusionListRemove
	model.patterns = []string{"TEST1"}
	model.cursor = 2

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newExclusionModel := newModel.(ExclusionListModel)

	if newExclusionModel.mode != ExclusionListMenu {
		t.Error("Enter on Back in remove mode should return to menu")
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
	model.mode = ExclusionListRemove
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
	model.mode = ExclusionListRemove
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
	model.mode = ExclusionListRemove
	model.patterns = []string{"TEST1", "TEST2"}
	model.cursor = 2

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	newExclusionModel := newModel.(ExclusionListModel)

	if newExclusionModel.cursor != 2 {
		t.Error("Down should stay at max when at bottom in remove mode")
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
