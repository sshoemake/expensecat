# ExpenseCat

```
 /\_/\  
( $.$ )  
(  >💳 )  ExpenseCat
```

## About

ExpenseCat is a command-line expense tracker with a terminal user interface (TUI) for managing personal finances. Track expenses, import from CSV, and view monthly spending reports organized by category.

## Features

- Import expenses from CSV files
- View monthly spending totals by category
- Track recurring expenses
- JSON file-based storage
- Terminal-based user interface (TUI)

## Installation

### Prerequisites
- Go 1.25+

### Build from source
```
go build -o bin/expensecat
```

### Quick Start
```
./bin/expensecat
```

## Usage

### Main Menu
- **Import Expenses** - Import CSV files into the tracker
- **View Monthly Totals** - View spending reports by month and category
- **Settings** - Configure application settings
- **Quit** - Exit the application

The import process reads CSV files from the base directory and adds new expenses to the storage.

### Settings

The Settings menu provides configuration options for the application:

- **Exclusion List** - Manage patterns to exclude during CSV import
  - Add new patterns to exclude matching transactions
  - Remove existing patterns
  - Patterns are case-insensitive and matched against transaction names

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `STORAGE_TYPE` | Storage backend (`json`) | `json` |
| `STORAGE_URL` | Storage path | `data` |
| `STORAGE_USER` | Database username | (empty) |
| `STORAGE_PASS` | Database password | (empty) |
| `STORAGE_SSL` | SSL mode (`disable`) | `disable` |
| `EXPENSE_BASE_PATH` | Import directory | `~/Downloads/Expenses-Mar-2026-Exports` |

## Technology

Built with:
- **Go** - Programming language
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal user interface framework for Go

## Development

### Commands

```bash
# Build the application
go build -o bin/expensecat

# Run the application
go run .

# Run tests
go test ./...

# Run linting
go vet ./...

# Format code
go fmt ./...
```

## Future Features

- **Auto-determine Categories** - Automatically categorize expenses based on transaction patterns
- **Modify Categories** - Add, edit, or remove expense categories
- **Modify Expenses** - Edit or delete imported expenses

## License

MIT License - See LICENSE file

## Acknowledgments

This project includes code adapted from [ExpenseOwl](https://github.com/Tanq16/ExpenseOwl), licensed under the MIT License.