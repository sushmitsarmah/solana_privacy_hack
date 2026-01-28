package cli

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// Run starts the CLI application
func Run() error {
	// Get API key from environment
	apiKey := os.Getenv("SHADOWPAY_API_KEY")

	// Create the model
	m := NewModel(apiKey)

	// Create the program
	p := tea.NewProgram(m, tea.WithAltScreen())

	// Run the program
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	return nil
}
