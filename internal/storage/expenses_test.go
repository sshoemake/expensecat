package storage

import (
	"testing"
	"time"
)

func TestGetExpensesByMonth_ValidMonth(t *testing.T) {
	mock := &mockStoreForExpenses{
		expenses: []Expense{
			{ID: "1", Name: "Coffee", Amount: 5.00, Date: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)},
			{ID: "2", Name: "Lunch", Amount: 12.00, Date: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC)},
		},
	}

	expenses, err := GetExpensesByMonth(mock, 2024, 1)
	if err != nil {
		t.Fatalf("GetExpensesByMonth() error = %v", err)
	}
	if len(expenses) != 2 {
		t.Errorf("GetExpensesByMonth() returned %d expenses, want 2", len(expenses))
	}
}

func TestGetExpensesByMonth_InvalidMonth_Zero(t *testing.T) {
	mock := &mockStoreForExpenses{expenses: []Expense{}}
	_, err := GetExpensesByMonth(mock, 2024, 0)
	if err == nil {
		t.Error("GetExpensesByMonth() expected error for month 0, got nil")
	}
}

func TestGetExpensesByMonth_InvalidMonth_Thirteen(t *testing.T) {
	mock := &mockStoreForExpenses{expenses: []Expense{}}
	_, err := GetExpensesByMonth(mock, 2024, 13)
	if err == nil {
		t.Error("GetExpensesByMonth() expected error for month 13, got nil")
	}
}

func TestGetExpensesByMonth_EmptyResult(t *testing.T) {
	mock := &mockStoreForExpenses{expenses: []Expense{}}

	expenses, err := GetExpensesByMonth(mock, 2024, 1)
	if err != nil {
		t.Fatalf("GetExpensesByMonth() error = %v", err)
	}
	if len(expenses) != 0 {
		t.Errorf("GetExpensesByMonth() returned %d expenses, want 0", len(expenses))
	}
}

func TestGetExpensesByMonth_FiltersCorrectMonth(t *testing.T) {
	mock := &mockStoreForExpenses{
		expenses: []Expense{
			{ID: "1", Name: "January Expense", Amount: 10.00, Date: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)},
			{ID: "2", Name: "February Expense", Amount: 20.00, Date: time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC)},
			{ID: "3", Name: "January Expense 2", Amount: 30.00, Date: time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC)},
		},
	}

	expenses, err := GetExpensesByMonth(mock, 2024, 1)
	if err != nil {
		t.Fatalf("GetExpensesByMonth() error = %v", err)
	}
	if len(expenses) != 2 {
		t.Errorf("GetExpensesByMonth() returned %d expenses, want 2", len(expenses))
	}
	for _, e := range expenses {
		if e.Date.Month() != time.January {
			t.Errorf("GetExpensesByMonth() returned expense with month %v, want January", e.Date.Month())
		}
	}
}

type mockStoreForExpenses struct {
	expenses []Expense
}

func (m *mockStoreForExpenses) Close() error                     { return nil }
func (m *mockStoreForExpenses) GetConfig() (*Config, error)      { return &Config{}, nil }
func (m *mockStoreForExpenses) GetCategories() ([]string, error) { return []string{}, nil }
func (m *mockStoreForExpenses) UpdateCategories([]string) error  { return nil }
func (m *mockStoreForExpenses) GetCurrency() (string, error)     { return "usd", nil }
func (m *mockStoreForExpenses) UpdateCurrency(string) error      { return nil }
func (m *mockStoreForExpenses) GetStartDate() (int, error)       { return 1, nil }
func (m *mockStoreForExpenses) UpdateStartDate(int) error        { return nil }
func (m *mockStoreForExpenses) GetRecurringExpenses() ([]RecurringExpense, error) {
	return []RecurringExpense{}, nil
}
func (m *mockStoreForExpenses) GetRecurringExpense(string) (RecurringExpense, error) {
	return RecurringExpense{}, nil
}
func (m *mockStoreForExpenses) AddRecurringExpense(RecurringExpense) error { return nil }
func (m *mockStoreForExpenses) RemoveRecurringExpense(string, bool) error  { return nil }
func (m *mockStoreForExpenses) UpdateRecurringExpense(string, RecurringExpense, bool) error {
	return nil
}
func (m *mockStoreForExpenses) GetAllExpenses(startDate, endDate *time.Time) ([]Expense, error) {
	var result []Expense
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
func (m *mockStoreForExpenses) GetExpense(string) (Expense, error)    { return Expense{}, nil }
func (m *mockStoreForExpenses) AddExpense(Expense) error              { return nil }
func (m *mockStoreForExpenses) RemoveExpense(string) error            { return nil }
func (m *mockStoreForExpenses) AddMultipleExpenses([]Expense) error   { return nil }
func (m *mockStoreForExpenses) RemoveMultipleExpenses([]string) error { return nil }
func (m *mockStoreForExpenses) UpdateExpense(string, Expense) error   { return nil }

func TestGetExpensesByMonthGroupedByCategory_Basic(t *testing.T) {
	mock := &mockStoreForExpenses{
		expenses: []Expense{
			{ID: "1", Name: "Coffee", Amount: 5.00, Category: "Food", Date: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)},
			{ID: "2", Name: "Lunch", Amount: 12.00, Category: "Food", Date: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC)},
			{ID: "3", Name: "Groceries", Amount: 100.00, Category: "Groceries", Date: time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC)},
		},
	}

	result, err := GetExpensesByMonthGroupedByCategory(mock, 2024, 1)
	if err != nil {
		t.Fatalf("GetExpensesByMonthGroupedByCategory() error = %v", err)
	}
	if len(result) != 2 {
		t.Errorf("GetExpensesByMonthGroupedByCategory() returned %d categories, want 2", len(result))
	}

	categoryTotals := make(map[string]float64)
	for _, ct := range result {
		categoryTotals[ct.Category] = ct.Total
	}
	if categoryTotals["Food"] != 17.00 {
		t.Errorf("Food total = %.2f, want 17.00", categoryTotals["Food"])
	}
	if categoryTotals["Groceries"] != 100.00 {
		t.Errorf("Groceries total = %.2f, want 100.00", categoryTotals["Groceries"])
	}
}

func TestGetExpensesByMonthGroupedByCategory_EmptyCategory(t *testing.T) {
	mock := &mockStoreForExpenses{
		expenses: []Expense{
			{ID: "1", Name: "Misc Item", Amount: 25.00, Category: "", Date: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)},
		},
	}

	result, err := GetExpensesByMonthGroupedByCategory(mock, 2024, 1)
	if err != nil {
		t.Fatalf("GetExpensesByMonthGroupedByCategory() error = %v", err)
	}
	if len(result) != 1 {
		t.Errorf("GetExpensesByMonthGroupedByCategory() returned %d categories, want 1", len(result))
	}
	if result[0].Category != "blank" {
		t.Errorf("Category = %q, want 'blank'", result[0].Category)
	}
	if result[0].Total != 25.00 {
		t.Errorf("Total = %.2f, want 25.00", result[0].Total)
	}
}

func TestGetExpensesByMonthGroupedByCategory_Sorted(t *testing.T) {
	mock := &mockStoreForExpenses{
		expenses: []Expense{
			{ID: "1", Name: "Zebra", Amount: 10.00, Category: "Zoo", Date: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)},
			{ID: "2", Name: "Apple", Amount: 5.00, Category: "Food", Date: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC)},
			{ID: "3", Name: "Bank Fee", Amount: 15.00, Category: "Banking", Date: time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC)},
		},
	}

	result, err := GetExpensesByMonthGroupedByCategory(mock, 2024, 1)
	if err != nil {
		t.Fatalf("GetExpensesByMonthGroupedByCategory() error = %v", err)
	}
	if len(result) != 3 {
		t.Errorf("GetExpensesByMonthGroupedByCategory() returned %d categories, want 3", len(result))
	}

	categories := []string{"Banking", "Food", "Zoo"}
	for i, expected := range categories {
		if result[i].Category != expected {
			t.Errorf("result[%d].Category = %q, want %q", i, result[i].Category, expected)
		}
	}
}

func TestGetExpensesByMonthGroupedByCategory_EmptyMonth(t *testing.T) {
	mock := &mockStoreForExpenses{
		expenses: []Expense{
			{ID: "1", Name: "February Expense", Amount: 20.00, Date: time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC)},
		},
	}

	result, err := GetExpensesByMonthGroupedByCategory(mock, 2024, 1)
	if err != nil {
		t.Fatalf("GetExpensesByMonthGroupedByCategory() error = %v", err)
	}
	if len(result) != 0 {
		t.Errorf("GetExpensesByMonthGroupedByCategory() returned %d categories, want 0", len(result))
	}
}

func TestExpense_Validate(t *testing.T) {
	testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		expense Expense
		wantErr bool
	}{
		{
			name:    "valid expense",
			expense: Expense{Name: "Coffee", Category: "Food", Amount: 5.00, Date: testDate},
			wantErr: false,
		},
		{
			name:    "valid expense with tags",
			expense: Expense{Name: "Coffee", Category: "Food", Amount: 5.00, Date: testDate, Tags: []string{"morning", "caffeine"}},
			wantErr: false,
		},
		{
			name:    "empty name",
			expense: Expense{Name: "", Category: "Food", Amount: 5.00, Date: testDate},
			wantErr: true,
		},
		{
			name:    "whitespace name",
			expense: Expense{Name: "   ", Category: "Food", Amount: 5.00, Date: testDate},
			wantErr: true,
		},
		{
			name:    "empty category - defaults to Uncategorized",
			expense: Expense{Name: "Coffee", Category: "", Amount: 5.00, Date: testDate},
			wantErr: false,
		},
		{
			name:    "zero amount",
			expense: Expense{Name: "Coffee", Category: "Food", Amount: 0, Date: testDate},
			wantErr: true,
		},
		{
			name:    "negative amount",
			expense: Expense{Name: "Coffee", Category: "Food", Amount: -5.00, Date: testDate},
			wantErr: false,
		},
		{
			name:    "empty date",
			expense: Expense{Name: "Coffee", Category: "Food", Amount: 5.00, Date: time.Time{}},
			wantErr: true,
		},
		{
			name:    "tags with invalid chars sanitized",
			expense: Expense{Name: "Coffee", Category: "Food", Amount: 5.00, Date: testDate, Tags: []string{"morning\t\n", "caf<fe>ine"}},
			wantErr: false,
		},
		{
			name:    "empty tags array",
			expense: Expense{Name: "Coffee", Category: "Food", Amount: 5.00, Date: testDate, Tags: []string{}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.expense.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Expense.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.name == "empty category - defaults to Uncategorized" && tt.expense.Category != "Uncategorized" {
				t.Errorf("Category = %q, want 'Uncategorized'", tt.expense.Category)
			}
		})
	}
}

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal string",
			input:    "Coffee Shop",
			expected: "Coffee Shop",
		},
		{
			name:     "multiple spaces collapsed",
			input:    "Coffee    Shop",
			expected: "Coffee Shop",
		},
		{
			name:     "leading/trailing spaces trimmed",
			input:    "  Coffee  ",
			expected: "Coffee",
		},
		{
			name:     "tabs and newlines replaced",
			input:    "Coffee\t\nShop",
			expected: "Coffee Shop",
		},
		{
			name:     "special chars replaced",
			input:    "Coffee<script>Shop",
			expected: "Coffee script Shop",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only invalid chars",
			input:    "!@#$%",
			expected: "!",
		},
		{
			name:     "SQL injection attempt",
			input:    "Coffee'; DROP TABLE expenses;--",
			expected: "Coffee' DROP TABLE expenses --",
		},
		{
			name:     "unicode letters preserved",
			input:    "Café Shop",
			expected: "Café Shop",
		},
		{
			name:     "emoji removed",
			input:    "Café ☕ Shop",
			expected: "Café Shop",
		},
		{
			name:     "ampersand removed",
			input:    "Food & Dining",
			expected: "Food Dining",
		},
		{
			name:     "angle brackets removed",
			input:    "Test <tag>",
			expected: "Test tag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeString(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeString(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateCategory(t *testing.T) {
	tests := []struct {
		name       string
		category   string
		wantErr    bool
		wantResult string
	}{
		{
			name:       "valid category",
			category:   "Food",
			wantErr:    false,
			wantResult: "Food",
		},
		{
			name:       "valid with spaces",
			category:   "  Groceries  ",
			wantErr:    false,
			wantResult: "Groceries",
		},
		{
			name:       "empty category",
			category:   "",
			wantErr:    true,
			wantResult: "",
		},
		{
			name:       "only spaces",
			category:   "   ",
			wantErr:    true,
			wantResult: "",
		},
		{
			name:       "only invalid chars",
			category:   "@#$%",
			wantErr:    true,
			wantResult: "",
		},
		{
			name:       "valid with special chars sanitized",
			category:   "Food&Dining",
			wantErr:    false,
			wantResult: "Food Dining",
		},
		{
			name:       "valid with exclamation",
			category:   "Food!",
			wantErr:    false,
			wantResult: "Food!",
		},
		{
			name:       "valid with quotes",
			category:   "Food\"",
			wantErr:    false,
			wantResult: "Food\"",
		},
		{
			name:       "valid with period",
			category:   "Food.",
			wantErr:    false,
			wantResult: "Food.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidateCategory(tt.category)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCategory() error = %v, wantErr %v", err, tt.wantErr)
			}
			if result != tt.wantResult {
				t.Errorf("ValidateCategory() = %q, want %q", result, tt.wantResult)
			}
		})
	}
}
