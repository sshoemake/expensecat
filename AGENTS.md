# Agent Guidelines for expensecat

This document provides guidance for agentic coding agents working in this repository.

## Project Overview

- **Module**: `expensecat`
- **Language**: Go 1.25.0
- **Type**: Expense tracking/CLI application
- **Entry Point**: `main.go`

## Build/Lint/Test Commands

### Build Commands
```bash
go build -o bin/expensecat # Build the application
go run .              # Run the application
go mod tidy           # Clean up dependencies
```

### Linting/Analysis
```bash
go fmt ./...          # Format code (standard go fmt)
go vet ./...          # Run Go vet for static analysis
golangci-lint run    # Run comprehensive linting (optional but recommended)
```

### Testing
```bash
go test ./...                    # Run all tests
go test ./... -v                 # Run all tests with verbose output
go test ./internal/storage       # Run tests in specific package
go test -run TestFunctionName   # Run single test by name
go test -cover                   # Run tests with coverage report
```

**Note**: Currently no tests exist in the codebase. Agents should add tests when modifying code.

## Code Style Guidelines

### Formatting
- Use `go fmt` for formatting (tabs for indentation)
- Run `go fmt` before committing changes

### Import Organization
Organize imports in three groups separated by blank lines:
```go
import (
    // Standard library
    "encoding/csv"
    "fmt"
    "os"
    "time"

    // Third-party imports
    "github.com/google/uuid"
    "github.com/lib/pq"

    // Local/internal imports
    "expensecat/internal/storage"
)
```

### Naming Conventions

| Element | Convention | Example |
|---------|------------|---------|
| Exported types/functions | PascalCase | `InitializeStorage`, `GetCategories` |
| Unexported types/functions | PascalCase | `jsonStore`, `getConfig` |
| Local variables/parameters | camelCase | `baseConfig`, `importedCount` |
| Constants | PascalCase/Screaming | `BackendTypeJSON` |
| Package names | lowercase, single word | `storage`, `api` |
| JSON fields | snake_case in tags | `json:"recurring_id"` |
| Go struct fields | PascalCase | `RecurringID string` |

### Type System

- Use Go's native type system with struct types
- Use JSON tags for serialization: `fieldName string `json:"json_name"`
- Use interfaces for abstraction (e.g., `Storage` interface)
- Use `time.Time` for date/time fields

Example struct:
```go
type Expense struct {
    ID          string    `json:"id"`
    RecurringID string    `json:"recurring_id"`
    Name        string    `json:"name"`
    Tags        []string  `json:"tags"`
    Category    string    `json:"category"`
    Amount      float64   `json:"amount"`
    Currency    string    `json:"currency"`
    Date        time.Time `json:"date"`
}
```

### Error Handling

Use descriptive error messages with `fmt.Errorf`:
```go
// Return errors from functions
func (s *jsonStore) GetCategories() ([]string, error) {
    config, err := s.GetConfig()
    if err != nil {
        return nil, err
    }
    return config.Categories, nil
}

// Create descriptive errors
return nil, fmt.Errorf("invalid storage type: %s", storageType)

// Validate input with clear messages
if e.Name == "" {
    return fmt.Errorf("expense 'name' cannot be empty")
}
```

### Validation Pattern

Use `Validate()` methods on types for input validation:
```go
func (e *Expense) Validate() error {
    e.Name = SanitizeString(e.Name)
    if e.Name == "" {
        return fmt.Errorf("expense 'name' cannot be empty")
    }
    if e.Category == "" {
        return fmt.Errorf("expense 'category' cannot be empty")
    }
    // ... additional validation
    return nil
}
```

### Concurrency

The JSON store uses `sync.RWMutex` for thread safety. Always acquire appropriate locks:
```go
func (s *jsonStore) GetExpenses() ([]Expense, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    // read operations
}

func (s *jsonStore) AddExpense(expense Expense) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    // write operations
}
```

## Architecture

### Storage Layer

The application uses an interface-based storage abstraction (`Storage` interface in `internal/storage/storage.go`):

- **JSON Backend** (`jsonStore`): File-based storage using `data/expenses.json` and `data/config.json`
- **PostgreSQL Backend** (`postgresStore`): Database storage (partially implemented)

### API Layer

The `internal/api` package contains shared API logic:

- **import.go**: CSV import functionality
  - `ImportCSVFile(store storage.Storage, filePath string) (*ImportResult, error)` - imports expenses from CSV
  - `ParseDate(dateStr string) (time.Time, error)` - parses dates with multiple format support

Example usage:
```go
result, err := api.ImportCSVFile(store, "path/to/file.csv")
if err != nil {
    log.Fatalf("Failed to import: %v", err)
}
log.Printf("Imported %d, Skipped %d", result.Imported, result.Skipped)
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `STORAGE_TYPE` | Storage backend: `json`, `postgres`, `sqlite` | `json` |
| `STORAGE_URL` | Storage path or database connection string | `data` |
| `STORAGE_USER` | Database username | (empty) |
| `STORAGE_PASS` | Database password | (empty) |
| `STORAGE_SSL` | SSL mode: `disable`, `require`, `verify-full`, `verify-ca` | `disable` |

## Testing Guidelines

- Test files should be named `*_test.go`
- Test functions should be named `TestFunctionName` or `TestType_Method`
- Use Go's standard `testing` package
- Follow table-driven test patterns for multiple test cases

Example test:
```go
func TestExpense_Validate(t *testing.T) {
    tests := []struct {
        name    string
        expense Expense
        wantErr bool
    }{
        {
            name:    "valid expense",
            expense: Expense{Name: "Coffee", Category: "Food", Amount: 5.00},
            wantErr: false,
        },
        {
            name:    "missing name",
            expense: Expense{Category: "Food", Amount: 5.00},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.expense.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Coverage Requirements

- **Target**: Maintain >80% test coverage across all packages
- **New functionality**: When adding new features or methods, write unit tests to cover them
- **Check coverage**: Run `go test -cover ./...` to verify coverage before completing changes
- **If coverage drops**: Add tests to restore >80% coverage before finishing the task
```

## General Guidelines

1. Always run `go fmt`, `go vet`, and optionally `golangci-lint` before completing changes
2. Add tests for new functionality
3. Use descriptive commit messages
4. Follow existing patterns in the codebase
5. Use `SanitizeString()` for user input to prevent injection issues
6. Return errors rather than panicking except in fatal initialization scenarios
7. Run `go build -o bin/expensecat` before completing changes to ensure the executable is up to date
