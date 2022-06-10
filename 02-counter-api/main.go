package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	count   int
	api     countAPIConfig
	loading bool
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
	return m.hitCounter
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case " ":
			if !m.loading {
				m.loading = true
				return m, m.hitCounter
			}
		}

	case countLoadedMsg:
		m.count = msg.count
		m.loading = false
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	if m.loading {
		return "..."
	} else {
		return fmt.Sprintf("count: %v", m.count)
	}
}

var _ tea.Model = (*model)(nil)

func initModel() tea.Model {
	return &model{
		api:     countAPIConfig{Key: "example"},
		loading: true,
	}
}

func main() {
	prog := tea.NewProgram(initModel())
	err := prog.Start()
	if err != nil {
		log.Fatal(err)
	}
}
