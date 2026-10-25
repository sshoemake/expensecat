package cli

import (
	"fmt"

	"expensecat/internal/version"
)

// appTitle returns the application name and version shown in screen headers.
// It is evaluated when called so the version injected at build time is used.
func appTitle() string {
	return fmt.Sprintf("ExpenseCat (%s)", version.Version)
}

// renderLogo fills in the version-dependent header text for a screen logo.
func renderLogo(format string) string {
	return fmt.Sprintf(format, appTitle())
}
