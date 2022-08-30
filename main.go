package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#5853e5")).
			Padding(1, 2).
			Bold(true)

	selectedItem = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#39d884")).
			BorderLeftForeground(lipgloss.Color("#39d884")).
			Padding(0, 2)

	spinnerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#39d884"))
)

type item struct {
	title   string
	desc    string
	command string
	args    []string
}

type model struct {
	list    list.Model
	spinner spinner.Model
	// items   []item
	choice  string
	loading bool
}

type commandEnd bool

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) Command() string     { return i.command }
func (i item) FilterValue() string { return i.title }

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) runCommand(i item) tea.Cmd {
	return func() tea.Msg {
		m.loading = false
		out, err := exec.Command("npm", i.args...).Output()

		if err != nil {
			log.Fatal("there was an error ", err, out)
		}

		return commandEnd(true)
	}
}

func items() []list.Item {
	return []list.Item{
		item{title: "Front-end - NPM Init", desc: "Initalise NPM", command: "npm", args: []string{"init", "-y"}},
		item{title: "Front-end - NPM Install Lodash", desc: "Installs Lodash", command: "npm", args: []string{"install", "lodash"}},
	}
}

func main() {
	// want to spread out multiple types; front end, back end, misc
	items := items()

	// Create a new default delegate
	delegate := list.NewDefaultDelegate()

	delegate.Styles.SelectedTitle = selectedItem
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedTitle.Copy()

	m := model{list: list.New(items, delegate, 0, 0)}
	m.list.Title = "Project Helper v0.1"
	m.list.Styles.Title = titleStyle

	m.spinner = spinner.NewModel()
	m.spinner.Spinner = spinner.Dot
	m.spinner.Style = spinnerStyle

	program := tea.NewProgram(m, tea.WithAltScreen())

	if err := program.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+w":
			m.loading = false
		case "ctrl+c", "ctrl+q":
			return m, tea.Quit
		case "enter":
			if !m.loading {
				i, ok := m.list.SelectedItem().(item)
				if ok {
					m.loading = true
					m.choice = string(i.title)
					var batch = tea.Batch(
						spinner.Tick,
						m.runCommand(i),
					)
					return m, batch
				}
			}

		}

	case commandEnd:
		m.loading = false

	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	if m.loading {
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.loading {
		return fmt.Sprintf("%sRunning %s", m.spinner.View(), m.choice)
	}
	return appStyle.Render(m.list.View())

}
