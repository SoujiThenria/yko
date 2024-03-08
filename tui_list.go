package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type accounts struct {
	acc string
}

func (a accounts) Title() string       { return a.acc[:strings.Index(a.acc, ":")] }
func (a accounts) Description() string { return a.acc[strings.Index(a.acc, ":")+1:] }
func (a accounts) FilterValue() string { return a.acc[:strings.Index(a.acc, ":")] }

func UpdateList(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "/" {
			m.list.FilterInput.SetValue("")
			break
		}
		if msg.String() == "enter" {
			if m.list.SettingFilter() {
				break
			}
			m.state = TOUCH_KEY
			return m, func() tea.Msg { return nextStep() }
		}
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := m.docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func ViewList(m model) string {
	return m.docStyle.Render(m.list.View())
}
