package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	projectName string
	list        list.Model
}

type scrapboxPage struct {
	Title_ string `json:"title"`
	ID     string `json:"id"`
}

// Description implements list.DefaultItem
func (p scrapboxPage) Description() string {
	return p.ID
}

// Title implements list.DefaultItem
func (p scrapboxPage) Title() string {
	return p.Title_
}

// FilterValue implements list.Item
func (p scrapboxPage) FilterValue() string {
	return p.Title_
}

var _ list.DefaultItem = (*scrapboxPage)(nil)

type scrapboxPagesLoadedMsg struct {
	pages []list.Item
}

func (m model) loadScrapboxPages() tea.Msg {
	resp, err := http.Get(
		"https://scrapbox.io/api/pages/" + m.projectName,
	)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	var pages []scrapboxPage
	result := struct{ Pages []scrapboxPage }{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		panic(err)
	}

	pages = result.Pages

	items := make([]list.Item, len(pages))
	for i, page := range pages {
		items[i] = page
	}
	return scrapboxPagesLoadedMsg{pages: items}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.loadScrapboxPages,
		m.list.StartSpinner(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	case scrapboxPagesLoadedMsg:
		m.list.StopSpinner()
		m.list.SetItems(msg.pages)
	}

	return m, cmd
}

func (m model) View() string {
	return m.list.View()
}

func initModel(projectName string) model {
	m := model{
		projectName: projectName,
		list:        list.New(nil, list.NewDefaultDelegate(), 1, 1),
	}
	// XXX required here to set showSpinner
	m.list.StartSpinner()
	return m
}

func main() {
	projectName := "help-jp"
	if len(os.Args) > 1 {
		projectName = os.Args[1]
	}

	prog := tea.NewProgram(initModel(projectName))
	err := prog.Start()
	if err != nil {
		log.Fatal(err)
	}
}
