package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	count int
}

func main() {
	err := tea.NewProgram(&Model{}, tea.WithAltScreen()).Start()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+q", "esc":
			return m, tea.Quit
		case "up":
			m.count++
		case "down":
			m.count--
		}

	}
	return m, nil
}

func (m *Model) View() string {
	return fmt.Sprintf("count: %d\n\n↑ increase  ↓ decrease", m.count)
}
