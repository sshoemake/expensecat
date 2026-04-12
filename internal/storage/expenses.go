package storage

import (
	"fmt"
	"sort"
	"time"
)

type CategoryTotal struct {
	Category string
	Total    float64
}

func GetExpensesByMonth(store Storage, year int, month int) ([]Expense, error) {
	if month < 1 || month > 12 {
		return nil, fmt.Errorf("invalid month: %d, must be between 1 and 12", month)
	}
	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return store.GetAllExpenses(&startOfMonth, &endOfMonth)
}

func GetExpensesByMonthGroupedByCategory(store Storage, year int, month int) ([]CategoryTotal, error) {
	expenses, err := GetExpensesByMonth(store, year, month)
	if err != nil {
		return nil, err
	}

	totals := make(map[string]float64)
	for _, e := range expenses {
		category := e.Category
		if category == "" {
			category = "blank"
		}
		totals[category] += e.Amount
	}

	result := make([]CategoryTotal, 0, len(totals))
	for category, total := range totals {
		result = append(result, CategoryTotal{Category: category, Total: total})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Category < result[j].Category
	})

	return result, nil
}
