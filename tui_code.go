package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"golang.design/x/clipboard"
)

type tickMsg time.Time

func UpdateCode(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

        switch {
        case key.Matches(msg, keys.Quit):
			return m, tea.Quit
        case key.Matches(msg, keys.Help):
            m.help.ShowAll = !m.help.ShowAll
        case key.Matches(msg, keys.Copy):
            clipboard.Write(clipboard.FmtText, []byte(m.code))
        case key.Matches(msg, keys.List):
            m.state = LIST_ACCOUNTS
            return m, nil
        }
	case tea.WindowSizeMsg:
		h, _ := m.docStyle.GetFrameSize()
		m.progress.Width = msg.Width - h
        if (m.progress.Width > 40) {
            m.progress.Width = 40
        }

		m.help.Width = msg.Width - h
        return m, nil
	case tickMsg:
		m.progressPerc += 0.0333333333
		if m.progressPerc > 1.0 {
			m.progressPerc = 1.0

			m.state = LIST_ACCOUNTS
			return m, nil
		}
		return m, tickCmd()
	}

	return m, nil
}

func ViewCode(m model) string {
	account := m.list.SelectedItem().(accounts)
	return m.docStyle.Render("Service: " + account.Title() + "\n" +
		"User: " + account.Description() + "\n" +
		"Code: " + m.code + "\n\n" +
		"Seconds remaining: " + fmt.Sprintf("%d", 30-int(m.progressPerc*30)) + "\n" +
		m.progress.ViewAs(m.progressPerc) + "\n\n\n" +
		m.help.View(keys))
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}



// TODO: Put this somewhere else...
type keyMap struct {
	Quit key.Binding
	Help key.Binding
	Copy key.Binding
	List key.Binding
}

var keys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "Quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "Toggle Help"),
	),
	Copy: key.NewBinding(
		key.WithKeys("enter", "c"),
		key.WithHelp("enter/c", "Copy"),
	),
	List: key.NewBinding(
		key.WithKeys("l", "ctrl+o"),
		key.WithHelp("l", "Back to List"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Copy, k.List, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help, k.Copy, k.List, k.Quit}, // second column
	}
}
