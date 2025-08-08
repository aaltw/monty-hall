package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/westhuis/monty-hall/pkg/config"
	"github.com/westhuis/monty-hall/pkg/ui"
)

func main() {
	// Initialize configuration manager
	configManager, err := config.NewManager()
	if err != nil {
		fmt.Printf("Error initializing configuration: %v\n", err)
		os.Exit(1)
	}

	// Create model with configuration
	model := ui.NewModelWithConfig(configManager)

	// Configure tea program based on config
	cfg := configManager.Get()
	var options []tea.ProgramOption

	// Always use alt screen for better experience
	options = append(options, tea.WithAltScreen())

	// Add mouse support if not in reduced motion mode
	if !cfg.UI.ReducedMotion {
		options = append(options, tea.WithMouseCellMotion())
	}

	p := tea.NewProgram(model, options...)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
