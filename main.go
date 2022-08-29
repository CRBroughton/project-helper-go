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
)

type item struct {
	title   string
	desc    string
	command string
	args    []string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) Command() string     { return i.command }
func (i item) FilterValue() string { return i.title }

type model struct {
	list    list.Model
	spinner spinner.Model
	// items   []item
	choice  string
	loading bool
}

func (m model) Init() tea.Cmd {
	return nil
}

type CommandType struct {
	name string
}

func (m model) runCommand(i item) tea.Cmd {
	return func() tea.Msg {
		m.loading = false
		out, err := exec.Command("npm", i.args...).Output()

		if err != nil {
			log.Fatal("there was an error ", err, out)
		}

		return CommandType{name: "npm-install"}
	}
}

func main() {
	items := []list.Item{
		item{title: "Front-end - NPM Init", desc: "Initalise NPM", command: "npm", args: []string{"init", "-y"}},
		item{title: "Front-end - NPM Install Lodash", desc: "Installs Lodash", command: "npm", args: []string{"install", "lodash"}},
	}

	// Create a new default delegate
	d := list.NewDefaultDelegate()

	// Change colors
	c := lipgloss.Color("#39d884")
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(c).BorderLeftForeground(c)
	d.Styles.SelectedDesc = d.Styles.SelectedTitle.Copy() // reuse the title style here

	m := model{list: list.New(items, d, 0, 0)}
	m.list.Title = "Project Helper v0.1"
	m.list.Styles.Title = titleStyle
	m.spinner = spinner.NewModel()
	m.spinner.Spinner = spinner.Dot
	m.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#39d884"))

	p := tea.NewProgram(m, tea.WithAltScreen())

	if err := p.Start(); err != nil {
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
					m.choice = string(i.title)
					m.loading = true
					var k = tea.Batch(
						m.runCommand(i),
						spinner.Tick,
					)
					return m, k
				}
			}

		}

	case CommandType:
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
	if !m.loading {
		return appStyle.Render(m.list.View())

	}
	if m.loading {
		return fmt.Sprintf("%sRunning %s", m.spinner.View(), m.choice)
	}
	return appStyle.Render(m.list.View())

}
