package storage

import (
	"os"
	"testing"
	"time"
)

func TestInitializeJsonStore(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	if store == nil {
		t.Fatal("store should not be nil")
	}
}

func TestJsonStore_GetConfig(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	cfg, err := store.GetConfig()
	if err != nil {
		t.Fatalf("GetConfig() error = %v", err)
	}
	if cfg == nil {
		t.Fatal("config should not be nil")
	}
	if len(cfg.Categories) == 0 {
		t.Error("categories should not be empty")
	}
}

func TestJsonStore_GetCategories(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	categories, err := store.GetCategories()
	if err != nil {
		t.Fatalf("GetCategories() error = %v", err)
	}
	if len(categories) == 0 {
		t.Error("categories should not be empty")
	}
}

func TestJsonStore_UpdateCategories(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	newCategories := []string{"Food", "Travel", "NewCategory"}
	err = store.UpdateCategories(newCategories)
	if err != nil {
		t.Fatalf("UpdateCategories() error = %v", err)
	}

	categories, err := store.GetCategories()
	if err != nil {
		t.Fatalf("GetCategories() error = %v", err)
	}
	if len(categories) != 3 {
		t.Errorf("categories length = %d, want 3", len(categories))
	}
}

func TestJsonStore_GetCurrency(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	currency, err := store.GetCurrency()
	if err != nil {
		t.Fatalf("GetCurrency() error = %v", err)
	}
	if currency != "usd" {
		t.Errorf("currency = %q, want 'usd'", currency)
	}
}

func TestJsonStore_UpdateCurrency(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	err = store.UpdateCurrency("eur")
	if err != nil {
		t.Fatalf("UpdateCurrency() error = %v", err)
	}

	currency, err := store.GetCurrency()
	if err != nil {
		t.Fatalf("GetCurrency() error = %v", err)
	}
	if currency != "eur" {
		t.Errorf("currency = %q, want 'eur'", currency)
	}
}

func TestJsonStore_GetStartDate(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	startDate, err := store.GetStartDate()
	if err != nil {
		t.Fatalf("GetStartDate() error = %v", err)
	}
	if startDate != 1 {
		t.Errorf("startDate = %d, want 1", startDate)
	}
}

func TestJsonStore_UpdateStartDate(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	err = store.UpdateStartDate(15)
	if err != nil {
		t.Fatalf("UpdateStartDate() error = %v", err)
	}

	startDate, err := store.GetStartDate()
	if err != nil {
		t.Fatalf("GetStartDate() error = %v", err)
	}
	if startDate != 15 {
		t.Errorf("startDate = %d, want 15", startDate)
	}
}

func TestJsonStore_GetExclusionList(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	list, err := store.GetExclusionList()
	if err != nil {
		t.Fatalf("GetExclusionList() error = %v", err)
	}
	if len(list) == 0 {
		t.Error("exclusion list should not be empty by default")
	}
}

func TestJsonStore_AddExclusion(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	err = store.AddExclusion("TEST PATTERN")
	if err != nil {
		t.Fatalf("AddExclusion() error = %v", err)
	}

	list, err := store.GetExclusionList()
	if err != nil {
		t.Fatalf("GetExclusionList() error = %v", err)
	}
	found := false
	for _, p := range list {
		if p == "TEST PATTERN" {
			found = true
			break
		}
	}
	if !found {
		t.Error("exclusion pattern should have been added")
	}
}

func TestJsonStore_AddExclusion_Duplicate(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	err = store.AddExclusion("DUPLICATE")
	if err != nil {
		t.Fatalf("AddExclusion() error = %v", err)
	}

	err = store.AddExclusion("DUPLICATE")
	if err == nil {
		t.Error("AddExclusion() should return error for duplicate")
	}
}

func TestJsonStore_AddExclusion_Empty(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	err = store.AddExclusion("   ")
	if err == nil {
		t.Error("AddExclusion() should return error for empty pattern")
	}
}

func TestJsonStore_RemoveExclusion(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	err = store.AddExclusion("REMOVE ME")
	if err != nil {
		t.Fatalf("AddExclusion() error = %v", err)
	}

	err = store.RemoveExclusion("REMOVE ME")
	if err != nil {
		t.Fatalf("RemoveExclusion() error = %v", err)
	}

	list, err := store.GetExclusionList()
	if err != nil {
		t.Fatalf("GetExclusionList() error = %v", err)
	}
	for _, p := range list {
		if p == "REMOVE ME" {
			t.Error("exclusion pattern should have been removed")
		}
	}
}

func TestJsonStore_RemoveExclusion_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	err = store.RemoveExclusion("NONEXISTENT")
	if err == nil {
		t.Error("RemoveExclusion() should return error for non-existent pattern")
	}
}

func TestJsonStore_UpdateExclusionList(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	newList := []string{"PATTERN1", "PATTERN2", "PATTERN3"}
	err = store.UpdateExclusionList(newList)
	if err != nil {
		t.Fatalf("UpdateExclusionList() error = %v", err)
	}

	list, err := store.GetExclusionList()
	if err != nil {
		t.Fatalf("GetExclusionList() error = %v", err)
	}
	if len(list) != 3 {
		t.Errorf("exclusion list length = %d, want 3", len(list))
	}
	if list[0] != "PATTERN1" || list[1] != "PATTERN2" || list[2] != "PATTERN3" {
		t.Error("exclusion list contents do not match expected values")
	}
}

func TestJsonStore_UpdateExclusionList_EmptyItems(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	err = store.UpdateExclusionList([]string{"VALID", "", "ALSO VALID"})
	if err == nil {
		t.Error("UpdateExclusionList() should return error for empty items")
	}
}

func TestJsonStore_AddExpense(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	expense := Expense{
		ID:       "1",
		Name:     "Coffee",
		Amount:   5.00,
		Category: "Food",
		Date:     time.Now(),
	}

	err = store.AddExpense(expense)
	if err != nil {
		t.Fatalf("AddExpense() error = %v", err)
	}
}

func TestJsonStore_GetExpense(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	expense := Expense{
		ID:       "test-123",
		Name:     "Coffee",
		Amount:   5.00,
		Category: "Food",
		Date:     time.Now(),
	}

	err = store.AddExpense(expense)
	if err != nil {
		t.Fatalf("AddExpense() error = %v", err)
	}

	result, err := store.GetExpense("test-123")
	if err != nil {
		t.Fatalf("GetExpense() error = %v", err)
	}
	if result.ID != "test-123" {
		t.Errorf("expense ID = %q, want 'test-123'", result.ID)
	}
}

func TestJsonStore_RemoveExpense(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	expense := Expense{
		ID:       "to-remove",
		Name:     "Coffee",
		Amount:   5.00,
		Category: "Food",
		Date:     time.Now(),
	}

	err = store.AddExpense(expense)
	if err != nil {
		t.Fatalf("AddExpense() error = %v", err)
	}

	err = store.RemoveExpense("to-remove")
	if err != nil {
		t.Fatalf("RemoveExpense() error = %v", err)
	}
}

func TestJsonStore_UpdateExpense(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	expense := Expense{
		ID:       "update-test",
		Name:     "Coffee",
		Amount:   5.00,
		Category: "Food",
		Date:     time.Now(),
	}

	err = store.AddExpense(expense)
	if err != nil {
		t.Fatalf("AddExpense() error = %v", err)
	}

	expense.Amount = 10.00
	err = store.UpdateExpense("update-test", expense)
	if err != nil {
		t.Fatalf("UpdateExpense() error = %v", err)
	}

	result, err := store.GetExpense("update-test")
	if err != nil {
		t.Fatalf("GetExpense() error = %v", err)
	}
	if result.Amount != 10.00 {
		t.Errorf("amount = %.2f, want 10.00", result.Amount)
	}
}

func TestJsonStore_GetAllExpenses(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	err = store.AddExpense(Expense{ID: "1", Name: "Coffee", Amount: 5.00, Date: time.Now()})
	if err != nil {
		t.Fatalf("AddExpense() error = %v", err)
	}
	err = store.AddExpense(Expense{ID: "2", Name: "Lunch", Amount: 10.00, Date: time.Now()})
	if err != nil {
		t.Fatalf("AddExpense() error = %v", err)
	}

	expenses, err := store.GetAllExpenses(nil, nil)
	if err != nil {
		t.Fatalf("GetAllExpenses() error = %v", err)
	}
	if len(expenses) != 2 {
		t.Errorf("expenses length = %d, want 2", len(expenses))
	}
}

func TestJsonStore_AddMultipleExpenses(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	expenses := []Expense{
		{ID: "1", Name: "Coffee", Amount: 5.00, Date: time.Now()},
		{ID: "2", Name: "Lunch", Amount: 10.00, Date: time.Now()},
	}

	err = store.AddMultipleExpenses(expenses)
	if err != nil {
		t.Fatalf("AddMultipleExpenses() error = %v", err)
	}

	result, err := store.GetAllExpenses(nil, nil)
	if err != nil {
		t.Fatalf("GetAllExpenses() error = %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expenses length = %d, want 2", len(result))
	}
}

func TestJsonStore_RemoveMultipleExpenses(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	err = store.AddExpense(Expense{ID: "1", Name: "Coffee", Amount: 5.00, Date: time.Now()})
	if err != nil {
		t.Fatalf("AddExpense() error = %v", err)
	}
	err = store.AddExpense(Expense{ID: "2", Name: "Lunch", Amount: 10.00, Date: time.Now()})
	if err != nil {
		t.Fatalf("AddExpense() error = %v", err)
	}

	err = store.RemoveMultipleExpenses([]string{"1", "2"})
	if err != nil {
		t.Fatalf("RemoveMultipleExpenses() error = %v", err)
	}

	result, err := store.GetAllExpenses(nil, nil)
	if err != nil {
		t.Fatalf("GetAllExpenses() error = %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expenses length = %d, want 0", len(result))
	}
}

func TestJsonStore_GetRecurringExpenses(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	result, err := store.GetRecurringExpenses()
	if err != nil {
		t.Fatalf("GetRecurringExpenses() error = %v", err)
	}
	if len(result) != 0 {
		t.Errorf("recurring expenses length = %d, want 0", len(result))
	}
}

func TestJsonStore_AddRecurringExpense(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	recurring := RecurringExpense{
		ID:          "rec-1",
		Name:        "Netflix",
		Amount:      15.99,
		Category:    "Entertainment",
		Interval:    "monthly",
		Occurrences: 12,
		StartDate:   time.Now(),
	}

	err = store.AddRecurringExpense(recurring)
	if err != nil {
		t.Fatalf("AddRecurringExpense() error = %v", err)
	}

	result, err := store.GetRecurringExpenses()
	if err != nil {
		t.Fatalf("GetRecurringExpenses() error = %v", err)
	}
	if len(result) != 1 {
		t.Errorf("recurring expenses length = %d, want 1", len(result))
	}
}

func TestJsonStore_GetRecurringExpense(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	recurring := RecurringExpense{
		ID:          "rec-get-test",
		Name:        "Netflix",
		Amount:      15.99,
		Category:    "Entertainment",
		Interval:    "monthly",
		Occurrences: 12,
		StartDate:   time.Now(),
	}

	err = store.AddRecurringExpense(recurring)
	if err != nil {
		t.Fatalf("AddRecurringExpense() error = %v", err)
	}

	result, err := store.GetRecurringExpense("rec-get-test")
	if err != nil {
		t.Fatalf("GetRecurringExpense() error = %v", err)
	}
	if result.ID != "rec-get-test" {
		t.Errorf("recurring expense ID = %q, want 'rec-get-test'", result.ID)
	}
}

func TestJsonStore_RemoveRecurringExpense(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	recurring := RecurringExpense{
		ID:          "rec-remove",
		Name:        "Netflix",
		Amount:      15.99,
		Category:    "Entertainment",
		Interval:    "monthly",
		Occurrences: 12,
		StartDate:   time.Now(),
	}

	err = store.AddRecurringExpense(recurring)
	if err != nil {
		t.Fatalf("AddRecurringExpense() error = %v", err)
	}

	err = store.RemoveRecurringExpense("rec-remove", false)
	if err != nil {
		t.Fatalf("RemoveRecurringExpense() error = %v", err)
	}

	result, err := store.GetRecurringExpenses()
	if err != nil {
		t.Fatalf("GetRecurringExpenses() error = %v", err)
	}
	if len(result) != 0 {
		t.Errorf("recurring expenses length = %d, want 0", len(result))
	}
}

func TestJsonStore_UpdateRecurringExpense(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	recurring := RecurringExpense{
		ID:          "rec-update",
		Name:        "Netflix",
		Amount:      15.99,
		Category:    "Entertainment",
		Interval:    "monthly",
		Occurrences: 12,
		StartDate:   time.Now(),
	}

	err = store.AddRecurringExpense(recurring)
	if err != nil {
		t.Fatalf("AddRecurringExpense() error = %v", err)
	}

	recurring.Amount = 19.99
	err = store.UpdateRecurringExpense("rec-update", recurring, false)
	if err != nil {
		t.Fatalf("UpdateRecurringExpense() error = %v", err)
	}

	result, err := store.GetRecurringExpense("rec-update")
	if err != nil {
		t.Fatalf("GetRecurringExpense() error = %v", err)
	}
	if result.Amount != 19.99 {
		t.Errorf("amount = %.2f, want 19.99", result.Amount)
	}
}

func TestJsonStore_Close(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}

	err = store.Close()
	if err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func TestInitializeStorage_Default(t *testing.T) {
	os.Unsetenv("STORAGE_TYPE")
	os.Unsetenv("STORAGE_URL")

	store, err := InitializeStorage()
	if err != nil {
		t.Fatalf("InitializeStorage() error = %v", err)
	}
	if store == nil {
		t.Fatal("store should not be nil")
	}
	defer store.Close()
}

func TestInitializeStorage_Json(t *testing.T) {
	os.Setenv("STORAGE_TYPE", "json")
	tmpDir := t.TempDir()
	os.Setenv("STORAGE_URL", tmpDir)
	defer os.Unsetenv("STORAGE_TYPE")
	defer os.Unsetenv("STORAGE_URL")

	store, err := InitializeStorage()
	if err != nil {
		t.Fatalf("InitializeStorage() error = %v", err)
	}
	if store == nil {
		t.Fatal("store should not be nil")
	}
	defer store.Close()
}

func TestInitializeStorage_InvalidType(t *testing.T) {
	os.Setenv("STORAGE_TYPE", "postgres")
	os.Unsetenv("STORAGE_URL")
	defer os.Unsetenv("STORAGE_TYPE")

	store, err := InitializeStorage()
	if err == nil {
		t.Fatal("InitializeStorage() expected error for invalid type")
	}
	if store != nil {
		t.Error("store should be nil for invalid type")
	}
}

func TestJsonStore_GetAllExpenses_WithDates(t *testing.T) {
	tmpDir := t.TempDir()
	config := SystemConfig{
		StorageURL:  tmpDir,
		StorageType: BackendTypeJSON,
	}

	store, err := InitializeJsonStore(config)
	if err != nil {
		t.Fatalf("InitializeJsonStore() error = %v", err)
	}
	defer store.Close()

	now := time.Now()
	startDate := now.AddDate(0, -1, 0)
	endDate := now.AddDate(0, 1, 0)

	err = store.AddExpense(Expense{ID: "1", Name: "Old", Amount: 5.00, Date: now.AddDate(0, -2, 0)})
	if err != nil {
		t.Fatalf("AddExpense() error = %v", err)
	}
	err = store.AddExpense(Expense{ID: "2", Name: "Current", Amount: 10.00, Date: now})
	if err != nil {
		t.Fatalf("AddExpense() error = %v", err)
	}

	expenses, err := store.GetAllExpenses(&startDate, &endDate)
	if err != nil {
		t.Fatalf("GetAllExpenses() error = %v", err)
	}
	if len(expenses) != 1 {
		t.Errorf("expenses length = %d, want 1", len(expenses))
	}
}

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

func (m *mockStoreForExpenses) Close() error                        { return nil }
func (m *mockStoreForExpenses) GetConfig() (*Config, error)         { return &Config{}, nil }
func (m *mockStoreForExpenses) GetCategories() ([]string, error)    { return []string{}, nil }
func (m *mockStoreForExpenses) UpdateCategories([]string) error     { return nil }
func (m *mockStoreForExpenses) GetCurrency() (string, error)        { return "usd", nil }
func (m *mockStoreForExpenses) UpdateCurrency(string) error         { return nil }
func (m *mockStoreForExpenses) GetStartDate() (int, error)          { return 1, nil }
func (m *mockStoreForExpenses) UpdateStartDate(int) error           { return nil }
func (m *mockStoreForExpenses) GetExclusionList() ([]string, error) { return []string{}, nil }
func (m *mockStoreForExpenses) AddExclusion(string) error           { return nil }
func (m *mockStoreForExpenses) RemoveExclusion(string) error        { return nil }
func (m *mockStoreForExpenses) UpdateExclusionList([]string) error  { return nil }
func (m *mockStoreForExpenses) GetImportPath() (string, error) {
	return "~/Downloads/Expense-Files", nil
}
func (m *mockStoreForExpenses) UpdateImportPath(string) error { return nil }
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

func TestConfig_SetBaseConfig(t *testing.T) {
	config := &Config{}
	config.SetBaseConfig()

	if config.Currency != "usd" {
		t.Errorf("Currency = %q, want 'usd'", config.Currency)
	}
	if config.StartDate != 1 {
		t.Errorf("StartDate = %d, want 1", config.StartDate)
	}
	if len(config.Categories) == 0 {
		t.Error("Categories should not be empty")
	}
	if config.RecurringExpenses == nil {
		t.Error("RecurringExpenses should not be nil")
	}
	if config.ExclusionList == nil {
		t.Error("ExclusionList should not be nil")
	}
}

func TestSystemConfig_SetStorageConfig(t *testing.T) {
	os.Setenv("STORAGE_TYPE", "json")
	os.Setenv("STORAGE_URL", "/custom/path")
	os.Setenv("STORAGE_SSL", "require")
	os.Setenv("STORAGE_USER", "user")
	os.Setenv("STORAGE_PASS", "pass")
	defer func() {
		os.Unsetenv("STORAGE_TYPE")
		os.Unsetenv("STORAGE_URL")
		os.Unsetenv("STORAGE_SSL")
		os.Unsetenv("STORAGE_USER")
		os.Unsetenv("STORAGE_PASS")
	}()

	config := &SystemConfig{}
	config.SetStorageConfig()

	if config.StorageType != BackendTypeJSON {
		t.Errorf("StorageType = %v, want %v", config.StorageType, BackendTypeJSON)
	}
	if config.StorageURL != "/custom/path" {
		t.Errorf("StorageURL = %q, want '/custom/path'", config.StorageURL)
	}
	if config.StorageSSL != "require" {
		t.Errorf("StorageSSL = %q, want 'require'", config.StorageSSL)
	}
	if config.StorageUser != "user" {
		t.Errorf("StorageUser = %q, want 'user'", config.StorageUser)
	}
	if config.StoragePass != "pass" {
		t.Errorf("StoragePass = %q, want 'pass'", config.StoragePass)
	}
}

func TestBackendTypeFromEnv(t *testing.T) {
	tests := []struct {
		input string
		want  BackendType
	}{
		{"json", BackendTypeJSON},
		{"postgres", BackendTypePostgres},
		{"sqlite", BackendTypeSQLite},
		{"unknown", BackendTypeJSON},
		{"", BackendTypeJSON},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := backendTypeFromEnv(tt.input)
			if got != tt.want {
				t.Errorf("backendTypeFromEnv(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestBackendURLFromEnv(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", "data"},
		{"custom", "custom"},
		{"/path/to/data", "/path/to/data"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := backendURLFromEnv(tt.input)
			if got != tt.want {
				t.Errorf("backendURLFromEnv(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestBackendSSLFromEnv(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"disable", "disable"},
		{"require", "require"},
		{"verify-full", "verify-full"},
		{"verify-ca", "verify-ca"},
		{"invalid", "disable"},
		{"", "disable"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := backendSSLFromEnv(tt.input)
			if got != tt.want {
				t.Errorf("backendSSLFromEnv(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSanitizeString_Additional(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello world", "hello world"},
		{"Test<Invalid>Chars", "Test Invalid Chars"},
		{"multiple   spaces", "multiple spaces"},
		{"", ""},
		{"  trimmed  ", "trimmed"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := SanitizeString(tt.input)
			if got != tt.want {
				t.Errorf("SanitizeString(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestExpense_Validate_AdditionalCases(t *testing.T) {
	tests := []struct {
		name    string
		expense Expense
		wantErr bool
	}{
		{
			name:    "valid expense",
			expense: Expense{Name: "Coffee", Category: "Food", Amount: 5.00, Date: time.Now()},
			wantErr: false,
		},
		{
			name:    "empty name",
			expense: Expense{Name: "", Amount: 5.00, Date: time.Now()},
			wantErr: true,
		},
		{
			name:    "zero amount",
			expense: Expense{Name: "Coffee", Amount: 0, Date: time.Now()},
			wantErr: true,
		},
		{
			name:    "zero date",
			expense: Expense{Name: "Coffee", Amount: 5.00, Date: time.Time{}},
			wantErr: true,
		},
		{
			name:    "empty category defaults to uncategorized",
			expense: Expense{Name: "Test", Amount: 5.00, Date: time.Now()},
			wantErr: false,
		},
		{
			name:    "with tags",
			expense: Expense{Name: "Test", Amount: 5.00, Date: time.Now(), Tags: []string{"tag1", "tag2"}},
			wantErr: false,
		},
		{
			name:    "with invalid tag chars",
			expense: Expense{Name: "Test", Amount: 5.00, Date: time.Now(), Tags: []string{"<tag>"}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.expense.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Expense.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRecurringExpense_Validate(t *testing.T) {
	tests := []struct {
		name    string
		expense RecurringExpense
		wantErr bool
	}{
		{
			name: "valid recurring",
			expense: RecurringExpense{
				Name:        "Netflix",
				Category:    "Entertainment",
				Amount:      15.99,
				Interval:    "monthly",
				Occurrences: 12,
				StartDate:   time.Now(),
			},
			wantErr: false,
		},
		{
			name:    "empty name",
			expense: RecurringExpense{Name: "", Category: "Food", Amount: 5.00},
			wantErr: true,
		},
		{
			name:    "empty category",
			expense: RecurringExpense{Name: "Test", Category: ""},
			wantErr: true,
		},
		{
			name:    "invalid interval",
			expense: RecurringExpense{Name: "Test", Category: "Food", Interval: "invalid"},
			wantErr: true,
		},
		{
			name:    "less than 2 occurrences",
			expense: RecurringExpense{Name: "Test", Category: "Food", Occurrences: 1},
			wantErr: true,
		},
		{
			name:    "zero start date",
			expense: RecurringExpense{Name: "Test", Category: "Food", StartDate: time.Time{}},
			wantErr: true,
		},
		{
			name:    "valid intervals",
			expense: RecurringExpense{Name: "Test", Category: "Food", Interval: "daily", Occurrences: 10, StartDate: time.Now()},
			wantErr: false,
		},
		{
			name:    "valid weekly",
			expense: RecurringExpense{Name: "Test", Category: "Food", Interval: "weekly", Occurrences: 10, StartDate: time.Now()},
			wantErr: false,
		},
		{
			name:    "valid yearly",
			expense: RecurringExpense{Name: "Test", Category: "Food", Interval: "yearly", Occurrences: 10, StartDate: time.Now()},
			wantErr: false,
		},
		{
			name:    "with tags",
			expense: RecurringExpense{Name: "Test", Category: "Food", Interval: "monthly", Occurrences: 10, StartDate: time.Now(), Tags: []string{"tag1"}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.expense.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("RecurringExpense.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
