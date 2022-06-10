package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	count   int
	api     countAPIConfig
	loading bool
	spinner spinner.Model
}

type countLoadedMsg struct {
	count int
}

func (m model) hitCounter() tea.Msg {
	count, err := m.api.Hit()
	if err != nil {
		panic(err)
	}

	return countLoadedMsg{count: count}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.hitCounter,
		m.spinner.Tick,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case " ":
			if !m.loading {
				m.loading = true
				cmd = tea.Batch(cmd, m.hitCounter)
				return m, cmd
			}
		}

	case countLoadedMsg:
		m.count = msg.count
		m.loading = false
		return m, cmd
	}

	return m, cmd
}

func (m model) View() string {
	if m.loading {
		return m.spinner.View()
	} else {
		return fmt.Sprintf("count: %v", m.count)
	}
}

var _ tea.Model = (*model)(nil)

func initModel() tea.Model {
	return &model{
		api:     countAPIConfig{Key: "example"},
		loading: true,
		spinner: spinner.New(),
	}
}

func main() {
	prog := tea.NewProgram(initModel())
	err := prog.Start()
	if err != nil {
		log.Fatal(err)
	}
}
