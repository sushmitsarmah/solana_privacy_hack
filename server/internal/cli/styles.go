package cli

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Color palette
	primaryColor   = lipgloss.Color("#7D56F4")
	secondaryColor = lipgloss.Color("#00D9FF")
	successColor   = lipgloss.Color("#00FF87")
	errorColor     = lipgloss.Color("#FF5F87")
	warningColor   = lipgloss.Color("#FFD700")
	textColor      = lipgloss.Color("#FFFFFF")
	subtleColor    = lipgloss.Color("#626262")

	// Title style
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryColor).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(0, 1).
		MarginBottom(1)

	// Header style
	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(secondaryColor).
		MarginBottom(1)

	// Menu item styles
	menuItemStyle = lipgloss.NewStyle().
		PaddingLeft(2).
		Foreground(textColor)

	selectedMenuItemStyle = lipgloss.NewStyle().
		PaddingLeft(1).
		Foreground(primaryColor).
		Bold(true)

	// Info box style
	infoBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		Padding(1, 2).
		Margin(1, 0)

	// Success message style
	successStyle = lipgloss.NewStyle().
		Foreground(successColor).
		Bold(true)

	// Error message style
	errorStyle = lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true)

	// Help style
	helpStyle = lipgloss.NewStyle().
		Foreground(subtleColor).
		Italic(true).
		MarginTop(1)

	// Input field style
	inputStyle = lipgloss.NewStyle().
		Foreground(textColor).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(primaryColor).
		Padding(0, 1)

	// Label style
	labelStyle = lipgloss.NewStyle().
		Foreground(secondaryColor).
		Bold(true)
)
