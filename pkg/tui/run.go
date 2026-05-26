package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jakej985-rgb/m3tal-core/pkg/client"
)

// Run starts the Bubble Tea program loop.
func Run(c *client.Client) error {
	model := NewModel(c)
	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
