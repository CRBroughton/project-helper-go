package main

import tea "github.com/charmbracelet/bubbletea"

type Model struct {
}

func main() {
	tea.NewProgram(&Model{}).Start()
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Model) View() string {
	return ""
}
