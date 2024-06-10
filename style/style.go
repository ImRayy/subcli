package style

import "github.com/charmbracelet/lipgloss"

var (
	ErrBg = lipgloss.NewStyle().
		Background(lipgloss.Color("#ef4444")).
		Foreground(lipgloss.Color("#e7e5e4")).
		PaddingLeft(1).
		PaddingRight(1).
		Bold(true)
	ErrColor     = lipgloss.NewStyle().Foreground(lipgloss.Color("#ef4444"))
	WarnColor    = lipgloss.NewStyle().Foreground(lipgloss.Color("#fbbf24"))
	InfoColor    = lipgloss.NewStyle().Foreground(lipgloss.Color("#06b6d4"))
	SuccessColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
)
