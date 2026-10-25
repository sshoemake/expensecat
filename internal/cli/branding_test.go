package cli

import (
	"strings"
	"testing"

	"expensecat/internal/version"
)

func setTestVersion(t *testing.T, v string) {
	t.Helper()
	original := version.Version
	t.Cleanup(func() {
		version.Version = original
	})
	version.Version = v
}

func TestAppTitle(t *testing.T) {
	setTestVersion(t, "v1.2.3")

	if got := appTitle(); got != "ExpenseCat (v1.2.3)" {
		t.Errorf("appTitle() = %q, want %q", got, "ExpenseCat (v1.2.3)")
	}
}

func TestRenderLogo_IncludesVersion(t *testing.T) {
	setTestVersion(t, "v9.9.9")

	logos := map[string]string{
		"main menu":        catFrames[0],
		"import":           importLogo,
		"expense list":     expenseListLogo,
		"report":           reportLogo,
		"recurring status": recurringStatusLogo,
	}

	for name, logo := range logos {
		t.Run(name, func(t *testing.T) {
			got := renderLogo(logo)
			if !strings.Contains(got, "ExpenseCat (v9.9.9)") {
				t.Errorf("rendered logo should contain version, got %q", got)
			}
		})
	}
}

func TestMainMenuModel_View_ShowsVersionInHeader(t *testing.T) {
	setTestVersion(t, "v2.0.0")

	view := NewMainMenuModel().View()
	if !strings.Contains(view, "ExpenseCat (v2.0.0)") {
		t.Error("View() should show the version in the header")
	}
	if strings.Contains(view, "expensecat v2.0.0") {
		t.Error("View() should not show the version in the footer")
	}
}

func TestScreenViews_ShowVersionInHeader(t *testing.T) {
	setTestVersion(t, "v4.5.6")

	store := &mockStoreForCLI{}
	views := map[string]string{
		"import":           NewImportModel(store, ".").View(),
		"expense list":     NewExpenseListModel(store).View(),
		"report":           NewReportModel(store).View(),
		"recurring status": NewRecurringStatusModel(store).View(),
	}

	for name, view := range views {
		t.Run(name, func(t *testing.T) {
			if !strings.Contains(view, "ExpenseCat (v4.5.6)") {
				t.Error("View() should show the version in the header")
			}
		})
	}
}
