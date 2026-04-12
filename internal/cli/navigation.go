package cli

type NavigationMsg struct {
	Destination int
}

const (
	ScreenMainMenu = iota
	ScreenImport
	ScreenReport
)
