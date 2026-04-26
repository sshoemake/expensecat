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
- Browse, search, and manage expenses
- View monthly spending totals by category
- **Recurring expense management** - Add, edit, and delete recurring payment schedules
- **Recurring expense status tracking** - Check if recurring payments have been made each month
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

| Option | Description |
|--------|-------------|
| **Import Expenses** | Import CSV files into the tracker |
| **Manage Expenses** | Browse, search, edit, and delete expenses |
| **View Monthly Totals** | View spending reports by month and category (←→ to navigate months) |
| **View Recurring Expenses** | Check status of recurring payments (paid/unpaid for current month) |
| **Settings** | Configure application settings |
| **Quit** | Exit the application |

### View Recurring Expenses

From the main menu, select "View Recurring Expenses" to see a list of your recurring payments and their status for the current month.

- Use **←→** keys to navigate between months
- Use **↑↓** to scroll through the list
- Each expense shows: name, amount, and paid status
- Press **Enter** or **Esc** to return to the main menu

### Settings

The Settings menu provides configuration options:

| Option | Description |
|--------|-------------|
| **Manage Categories** | Add, edit, and delete expense categories |
| **Recurring Expenses** | Add, edit, and delete recurring payment schedules |
| **Exclusion List** | Manage patterns to exclude during CSV import |

#### Recurring Expenses

Configure recurring payments from **Settings → Recurring Expenses**:

- **Name** - Description of the payment
- **Category** - Which category to assign matching expenses
- **Interval** - How often (daily, weekly, monthly, yearly)
- **Occurrences** - How many times to repeat
- **Enabled** - Whether the payment is active

Matching is done by comparing recurring expense names against expense names (case-insensitive substring match).

#### Exclusion List

Manage patterns to exclude during CSV import:

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
| `EXPENSE_BASE_PATH` | Import directory | `~/Downloads/Expense-Files` |

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
- [x] **Re-occurring Expenses cross-check** - List and report monthly expenses and if they have been paid or not (✓ Done)
- [x] **Modify Categories** - Add, edit, or remove expense categories (✓ Done)
- [x] **Modify Expenses** - Edit or delete imported expenses (✓ Done)

## License

MIT License - See LICENSE file

## Acknowledgments

This project includes code adapted from [ExpenseOwl](https://github.com/Tanq16/ExpenseOwl), licensed under the MIT License.