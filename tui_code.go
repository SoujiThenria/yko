package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"golang.design/x/clipboard"
)

var testCode string

func codeUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "enter" {
            m.selected = false
            return m, frame()
        }
        if msg.String() == "c" {
            clipboard.Write(clipboard.FmtText, []byte(testCode))
        }
    case tea.WindowSizeMsg:
        h, v := docStyle.GetFrameSize()
        m.list.SetSize(msg.Width-h, msg.Height-v)
    }

    var cmd tea.Cmd
    m.list, cmd = m.list.Update(msg)
    return m, cmd
}

func codeView(m model) string {
    if testCode == "" {
        code, err := yko.Calculate(m.data[m.sel_item])
        if err != nil {
            code = ""
        }
        testCode = code
    }
    return fmt.Sprintf("Account: %s\nCode: %s\n\n\nc - copy, q - quit", m.data[m.sel_item], testCode)
}
