package main

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type model struct {
    list list.Model
    selected bool
    sel_item int
    data []string
    code string
}

type (
	frameMsg struct{}
)

func frame() tea.Cmd {
	return tea.Tick(time.Second/60, func(time.Time) tea.Msg {
		return frameMsg{}
	})
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
        if msg.String() == "q" {
			return m, tea.Quit
		}
	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	if m.selected {
        return codeUpdate(msg, m)
	}
    return listUpdate(msg, m)
}

func (m model) View() string {
	if m.selected {
        return codeView(m)
	}
    return listView(m)
}

func startTUI(dataList []string) {
    items := []list.Item{}
    for _,d := range dataList {
        items = append(items, tuiListData{title: d[:strings.Index(d, ":")], description: d[strings.Index(d, ":")+1:]})
    }
    m := model{
        data: dataList,
        list: list.New(items, list.NewDefaultDelegate(), 0, 0),
    }
    m.list.Title = "OATH Accounts"

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal("could not start program:", err)
	}
}
