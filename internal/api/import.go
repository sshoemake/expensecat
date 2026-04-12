package api

import (
	"encoding/csv"
	"expensecat/internal/storage"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type ImportResult struct {
	Imported   int
	Skipped    int
	Duplicates int
}

func getImportLogger() *log.Logger {
	logFile, err := os.OpenFile("data/imports.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return log.Default()
	}
	return log.New(logFile, "", log.LstdFlags)
}

func ImportCSVFile(store storage.Storage, filePath string) (*ImportResult, error) {
	logger := getImportLogger()
	filename := filePath

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error retrieving the file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %w", err)
	}
	if len(records) < 2 {
		return nil, fmt.Errorf("CSV file must have a header and at least one data row")
	}

	header := records[0]
	colMap := make(map[string]int)
	for i, col := range header {
		colMap[strings.ToLower(strings.TrimSpace(col))] = i
	}

	requiredCols := []string{"name", "amount", "date"}
	for _, col := range requiredCols {
		if _, ok := colMap[col]; !ok {
			return nil, fmt.Errorf("missing required column: [%s]", col)
		}
	}

	currentCategories, err := store.GetCategories()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve current categories: %w", err)
	}
	categorySet := make(map[string]bool)
	for _, cat := range currentCategories {
		categorySet[strings.ToLower(cat)] = true
	}

	config, err := store.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve config: %w", err)
	}
	exclusionList := config.ExclusionList

	var importedCount, skippedCount, duplicatesCount int

	logger.Printf("[INFO] Starting import of file: %s", filename)

	for rowNum, record := range records[1:] {
		csvRowNum := rowNum + 2

		if len(record) != len(header) {
			logger.Printf("[SKIP] %s row %d: column count mismatch (expected %d, got %d)", filename, csvRowNum, len(header), len(record))
			skippedCount++
			continue
		}

		amount, err := strconv.ParseFloat(record[colMap["amount"]], 64)
		if err != nil {
			logger.Printf("[SKIP] %s row %d: invalid amount '%s'", filename, csvRowNum, record[colMap["amount"]])
			skippedCount++
			continue
		}
		if amount < 0 {
			logger.Printf("[WARN] %s row %d: negative amount '%f' converted to positive", filename, csvRowNum, amount)
			amount = -amount
		}
		date, err := ParseDate(record[colMap["date"]])
		if err != nil {
			logger.Printf("[SKIP] %s row %d: invalid date '%s': %v", filename, csvRowNum, record[colMap["date"]], err)
			skippedCount++
			continue
		}
		category := strings.TrimSpace(record[colMap["category"]])
		categorySet[strings.ToLower(category)] = true

		expense := storage.Expense{
			Name:     strings.TrimSpace(record[colMap["name"]]),
			Category: category,
			Amount:   amount,
			Date:     date,
		}
		if err := expense.Validate(); err != nil {
			logger.Printf("[SKIP] %s row %d: validation failed for '%s': %v", filename, csvRowNum, expense.Name, err)
			skippedCount++
			continue
		}
		expenseNameLower := strings.ToLower(expense.Name)
		isExcluded := false
		var matchedExclusion string
		for _, exclusion := range exclusionList {
			if strings.Contains(expenseNameLower, strings.ToLower(exclusion)) {
				isExcluded = true
				matchedExclusion = exclusion
				break
			}
		}
		if isExcluded {
			logger.Printf("[SKIP] %s row %d: '%s' matches exclusion pattern '%s'", filename, csvRowNum, expense.Name, matchedExclusion)
			skippedCount++
			continue
		}
		if duplicate, dupDesc := findDuplicate(store, expense); duplicate != nil {
			logger.Printf("[SKIP] %s row %d: duplicate of existing expense (%s)", filename, csvRowNum, dupDesc)
			duplicatesCount++
			continue
		}
		if err := store.AddExpense(expense); err != nil {
			logger.Printf("[SKIP] %s row %d: failed to add expense '%s': %v", filename, csvRowNum, expense.Name, err)
			skippedCount++
			continue
		}
		importedCount++
		time.Sleep(10 * time.Millisecond)
	}

	logger.Printf("[INFO] Import complete: imported=%d, skipped=%d, duplicates=%d", importedCount, skippedCount, duplicatesCount)

	return &ImportResult{
		Imported:   importedCount,
		Skipped:    skippedCount,
		Duplicates: duplicatesCount,
	}, nil
}

func ParseDate(dateStr string) (time.Time, error) {
	dateFormats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"2006-1-2",
		"2006/01/02",
		"2006/1/2",
		"01/02/2006",
	}
	for _, format := range dateFormats {
		if d, err := time.Parse(format, dateStr); err == nil {
			return d.UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

func findDuplicate(store storage.Storage, expense storage.Expense) (*storage.Expense, string) {
	startOfDay := time.Date(expense.Date.Year(), expense.Date.Month(), expense.Date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)
	existingExpenses, err := store.GetAllExpenses(&startOfDay, &endOfDay)
	if err != nil {
		return nil, ""
	}
	newExpenseNameLower := strings.ToLower(expense.Name)
	for i := range existingExpenses {
		existing := &existingExpenses[i]
		if strings.EqualFold(newExpenseNameLower, strings.ToLower(existing.Name)) &&
			existing.Amount == expense.Amount {
			desc := fmt.Sprintf("name=\"%s\", amount=%.2f, date=%s", existing.Name, existing.Amount, expense.Date.Format("2006-01-02"))
			return existing, desc
		}
	}
	return nil, ""
}
