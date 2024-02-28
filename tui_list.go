package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type tuiListData struct {
    title string
    description string
}

func (d tuiListData ) Title() string       { return d.title }
func (d tuiListData ) Description() string { return d.description }
func (d tuiListData ) FilterValue() string { return d.title }

func listUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "enter" {
			if m.list.SettingFilter() {
				break
			}
            m.selected = true
            test := m.list.SelectedItem().(tuiListData)
            for i := 0; i < len(m.data); i++ {
                if strings.Contains(m.data[i], test.title) {
                    m.sel_item = i
                    testCode = ""
                    return m, frame()
                }
            }
            // m.sel_item = m.list.Index()
            return m, frame()
        }
    case tea.WindowSizeMsg:
        h, v := docStyle.GetFrameSize()
        m.list.SetSize(msg.Width-h, msg.Height-v)
    }

    var cmd tea.Cmd
    m.list, cmd = m.list.Update(msg)
    return m, cmd
}

func listView(m model) string {
    return docStyle.Render(m.list.View())
}
