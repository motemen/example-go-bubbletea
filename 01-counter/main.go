package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	count int
}

func (model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case " ":
			m.count = m.count + 1
			return m, nil
		}
	}

	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("count: %v", m.count)
}

var _ tea.Model = (*model)(nil)

func main() {
	prog := tea.NewProgram(model{count: 0})
	err := prog.Start()
	if err != nil {
		log.Fatal(err)
	}
}
