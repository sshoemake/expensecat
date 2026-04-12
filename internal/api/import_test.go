package api

import (
	"expensecat/internal/storage"
	"os"
	"testing"
	"time"
)

type mockStore struct {
	expenses []storage.Expense
	config   *storage.Config
}

func (m *mockStore) Close() error {
	return nil
}

func (m *mockStore) GetConfig() (*storage.Config, error) {
	if m.config != nil {
		return m.config, nil
	}
	return &storage.Config{}, nil
}

func (m *mockStore) GetCategories() ([]string, error) {
	return []string{}, nil
}

func (m *mockStore) UpdateCategories(categories []string) error {
	return nil
}

func (m *mockStore) GetCurrency() (string, error) {
	return "usd", nil
}

func (m *mockStore) UpdateCurrency(currency string) error {
	return nil
}

func (m *mockStore) GetStartDate() (int, error) {
	return 1, nil
}

func (m *mockStore) UpdateStartDate(startDate int) error {
	return nil
}

func (m *mockStore) GetRecurringExpenses() ([]storage.RecurringExpense, error) {
	return []storage.RecurringExpense{}, nil
}

func (m *mockStore) GetRecurringExpense(id string) (storage.RecurringExpense, error) {
	return storage.RecurringExpense{}, nil
}

func (m *mockStore) AddRecurringExpense(expense storage.RecurringExpense) error {
	return nil
}

func (m *mockStore) RemoveRecurringExpense(id string, removeAll bool) error {
	return nil
}

func (m *mockStore) UpdateRecurringExpense(id string, recurringExpense storage.RecurringExpense, updateAll bool) error {
	return nil
}

func (m *mockStore) GetAllExpenses(startDate, endDate *time.Time) ([]storage.Expense, error) {
	var result []storage.Expense
	for _, exp := range m.expenses {
		if startDate != nil && exp.Date.Before(*startDate) {
			continue
		}
		if endDate != nil && exp.Date.After(*endDate) {
			continue
		}
		result = append(result, exp)
	}
	return result, nil
}

func (m *mockStore) GetExpense(id string) (storage.Expense, error) {
	return storage.Expense{}, nil
}

func (m *mockStore) AddExpense(expense storage.Expense) error {
	m.expenses = append(m.expenses, expense)
	return nil
}

func (m *mockStore) RemoveExpense(id string) error {
	return nil
}

func (m *mockStore) AddMultipleExpenses(expenses []storage.Expense) error {
	m.expenses = append(m.expenses, expenses...)
	return nil
}

func (m *mockStore) RemoveMultipleExpenses(ids []string) error {
	return nil
}

func (m *mockStore) UpdateExpense(id string, expense storage.Expense) error {
	return nil
}

func TestFindDuplicate(t *testing.T) {
	testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		existing    []storage.Expense
		expense     storage.Expense
		wantFound   bool
		wantDescLen int
	}{
		{
			name:        "exact duplicate",
			existing:    []storage.Expense{{Name: "Coffee", Amount: 5.00, Date: testDate}},
			expense:     storage.Expense{Name: "Coffee", Amount: 5.00, Date: testDate},
			wantFound:   true,
			wantDescLen: 40,
		},
		{
			name:        "duplicate with different case in name",
			existing:    []storage.Expense{{Name: "COFFEE", Amount: 5.00, Date: testDate}},
			expense:     storage.Expense{Name: "coffee", Amount: 5.00, Date: testDate},
			wantFound:   true,
			wantDescLen: 40,
		},
		{
			name:        "no duplicate - different name",
			existing:    []storage.Expense{{Name: "Coffee", Amount: 5.00, Date: testDate}},
			expense:     storage.Expense{Name: "Tea", Amount: 5.00, Date: testDate},
			wantFound:   false,
			wantDescLen: 0,
		},
		{
			name:        "no duplicate - different amount",
			existing:    []storage.Expense{{Name: "Coffee", Amount: 5.00, Date: testDate}},
			expense:     storage.Expense{Name: "Coffee", Amount: 6.00, Date: testDate},
			wantFound:   false,
			wantDescLen: 0,
		},
		{
			name:        "no duplicate - different date",
			existing:    []storage.Expense{{Name: "Coffee", Amount: 5.00, Date: testDate}},
			expense:     storage.Expense{Name: "Coffee", Amount: 5.00, Date: testDate.Add(24 * time.Hour)},
			wantFound:   false,
			wantDescLen: 0,
		},
		{
			name:        "empty existing expenses",
			existing:    []storage.Expense{},
			expense:     storage.Expense{Name: "Coffee", Amount: 5.00, Date: testDate},
			wantFound:   false,
			wantDescLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mockStore{expenses: tt.existing}
			duplicate, desc := findDuplicate(store, tt.expense)
			if (duplicate != nil) != tt.wantFound {
				t.Errorf("findDuplicate() found = %v, want found = %v", duplicate != nil, tt.wantFound)
			}
			if len(desc) < tt.wantDescLen && tt.wantFound {
				t.Errorf("findDuplicate() desc too short, got %d chars, want at least %d", len(desc), tt.wantDescLen)
			}
		})
	}
}

func TestImportResult_Duplicates(t *testing.T) {
	result := ImportResult{
		Imported:   10,
		Skipped:    2,
		Duplicates: 5,
	}

	if result.Imported != 10 {
		t.Errorf("ImportResult.Imported = %d, want 10", result.Imported)
	}
	if result.Skipped != 2 {
		t.Errorf("ImportResult.Skipped = %d, want 2", result.Skipped)
	}
	if result.Duplicates != 5 {
		t.Errorf("ImportResult.Duplicates = %d, want 5", result.Duplicates)
	}
}

func TestParseDate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"RFC3339", "2024-01-15T10:30:00Z", false},
		{"ISO date", "2024-01-15", false},
		{"slash date", "2024/01/15", false},
		{"invalid", "not-a-date", true},
		{"space separated", "2024-01-15 10:30:00", false},
		{"single digit day", "2024-1-5", false},
		{"single digit month", "2024-1-15", false},
		{"MM/DD/YYYY", "01/15/2024", false},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseDate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDate(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestImportCSVFile_DuplicateDetection(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := tmpDir + "/test.csv"

	testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	existingExpenses := []storage.Expense{
		{Name: "Coffee", Amount: 5.00, Date: testDate},
	}
	store := &mockStore{expenses: existingExpenses}

	csvContent := "name,amount,date,category\nCoffee,5.00,2024-01-15,Groceries\nTea,3.00,2024-01-16,Groceries\n"
	if err := os.WriteFile(csvPath, []byte(csvContent), 0644); err != nil {
		t.Fatalf("Failed to write test CSV: %v", err)
	}

	result, err := ImportCSVFile(store, csvPath)
	if err != nil {
		t.Fatalf("ImportCSVFile() error = %v", err)
	}

	if result.Imported != 1 {
		t.Errorf("ImportCSVFile() Imported = %d, want 1", result.Imported)
	}
	if result.Duplicates != 1 {
		t.Errorf("ImportCSVFile() Duplicates = %d, want 1", result.Duplicates)
	}
	if result.Skipped != 0 {
		t.Errorf("ImportCSVFile() Skipped = %d, want 0", result.Skipped)
	}
}

func TestImportCSVFile_MissingRequiredColumns(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := tmpDir + "/test.csv"

	csvContent := "name,amount\nCoffee,5.00\n"
	if err := os.WriteFile(csvPath, []byte(csvContent), 0644); err != nil {
		t.Fatalf("Failed to write test CSV: %v", err)
	}

	store := &mockStore{}
	_, err := ImportCSVFile(store, csvPath)
	if err == nil {
		t.Error("ImportCSVFile() expected error for missing 'date' column, got nil")
	}
}

func TestImportCSVFile_InvalidAmount(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := tmpDir + "/test.csv"

	csvContent := "name,amount,date,category\nCoffee,invalid,2024-01-15,Food\n"
	if err := os.WriteFile(csvPath, []byte(csvContent), 0644); err != nil {
		t.Fatalf("Failed to write test CSV: %v", err)
	}

	store := &mockStore{}
	result, err := ImportCSVFile(store, csvPath)
	if err != nil {
		t.Fatalf("ImportCSVFile() error = %v", err)
	}

	if result.Skipped != 1 {
		t.Errorf("ImportCSVFile() Skipped = %d, want 1", result.Skipped)
	}
	if result.Imported != 0 {
		t.Errorf("ImportCSVFile() Imported = %d, want 0", result.Imported)
	}
}

func TestImportCSVFile_InvalidDate(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := tmpDir + "/test.csv"

	csvContent := "name,amount,date,category\nCoffee,5.00,not-a-date,Food\n"
	if err := os.WriteFile(csvPath, []byte(csvContent), 0644); err != nil {
		t.Fatalf("Failed to write test CSV: %v", err)
	}

	store := &mockStore{}
	result, err := ImportCSVFile(store, csvPath)
	if err != nil {
		t.Fatalf("ImportCSVFile() error = %v", err)
	}

	if result.Skipped != 1 {
		t.Errorf("ImportCSVFile() Skipped = %d, want 1", result.Skipped)
	}
	if result.Imported != 0 {
		t.Errorf("ImportCSVFile() Imported = %d, want 0", result.Imported)
	}
}

func TestImportCSVFile_ExclusionMatch(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := tmpDir + "/test.csv"

	store := &mockStore{
		config: &storage.Config{
			ExclusionList: []string{"test", "sample"},
		},
	}

	csvContent := "name,amount,date,category\nTest Expense,5.00,2024-01-15,Food\nReal Expense,3.00,2024-01-16,Food\n"
	if err := os.WriteFile(csvPath, []byte(csvContent), 0644); err != nil {
		t.Fatalf("Failed to write test CSV: %v", err)
	}

	result, err := ImportCSVFile(store, csvPath)
	if err != nil {
		t.Fatalf("ImportCSVFile() error = %v", err)
	}

	if result.Skipped != 1 {
		t.Errorf("ImportCSVFile() Skipped = %d, want 1", result.Skipped)
	}
	if result.Imported != 1 {
		t.Errorf("ImportCSVFile() Imported = %d, want 1", result.Imported)
	}
}

func TestImportCSVFile_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := tmpDir + "/test.csv"

	csvContent := "name,amount,date,category\n"
	if err := os.WriteFile(csvPath, []byte(csvContent), 0644); err != nil {
		t.Fatalf("Failed to write test CSV: %v", err)
	}

	store := &mockStore{}
	_, err := ImportCSVFile(store, csvPath)
	if err == nil {
		t.Error("ImportCSVFile() expected error for empty file, got nil")
	}
}
