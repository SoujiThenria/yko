package main

import (
	"fmt"
	"os"
	"strings"
	"yko/ykoath"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/ebfe/scard"
	"golang.design/x/clipboard"
	"golang.org/x/term"
)

type model struct {
	state        tui_state
	key          *ykoath.YKO
	docStyle     lipgloss.Style
	list         list.Model
	listName     string
	progress     progress.Model
	progressPerc float64
	help         help.Model
	code         string
}

type executeNextUpdate struct{}

func nextStep() executeNextUpdate {
	return executeNextUpdate{}
}

func main() {
	if err := clipboard.Init(); err != nil {
		log.Fatal("Cannot init the clipboard", "error", err)
	}

	ctx, card, err := setupYubiKey()
	if err != nil {
		log.Fatal("Failed to setup YubiKey", "error", err)
	}
	defer func() {
		if err := card.Disconnect(scard.ResetCard); err != nil {
			log.Warn("Failed to disconnect YubiKey", "error", err)
		}
		if err := ctx.Release(); err != nil {
			log.Warn("Failed to release scard context", "error", err)
		}
	}()

	yKey := ykoath.New(card)
	if err := yKey.Select(); err != nil {
		log.Fatal("Failed to select", "error", err)
	}
	if yKey.AuthRequired() {
		fmt.Print("Password: ")
		pass, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println()
		if err := yKey.Validate(pass); err != nil {
			log.Fatal(err)
		}
	}
	keyList, err := yKey.List()
	if err != nil {
		log.Fatal("Failed to list", "error", err)
	}

	a := make([]list.Item, len(keyList))
	for i := 0; i < len(keyList); i++ {
		a[i] = accounts{acc: keyList[i].Name}
	}

	m := model{
		state: LIST_ACCOUNTS,
		key:   yKey,
		docStyle: lipgloss.NewStyle().Margin(1, 2),
		list:         list.New(a, list.NewDefaultDelegate(), 0, 0),
		listName:     "YubiKey Accounts",
		progress:     progress.New(),
		progressPerc: 0,
		help:         help.New(),
		code:         "",
	}
	m.list.Title = m.listName
	m.progress.ShowPercentage = false

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal("Failed to run the TUI", "error", err)
	}
}

func setupYubiKey() (*scard.Context, *scard.Card, error) {
	ctx, err := scard.EstablishContext()
	if err != nil {
		return nil, nil, err
	}

	readers, err := ctx.ListReaders()
	if err != nil {
		return nil, nil, err
	}

	for _, reader := range readers {
		if strings.Contains(reader, "Yubico YubiKey") {
			card, err := ctx.Connect(reader, scard.ShareExclusive, scard.ProtocolAny)
			if err != nil {
				return nil, nil, err
			}
			return ctx, card, nil
		}
	}

	return nil, nil, fmt.Errorf("No YoubiKey was found")
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case LIST_ACCOUNTS:
		return UpdateList(msg, m)
	case TOUCH_KEY:
		switch msg.(type) {
		case executeNextUpdate:
			code, err := m.key.Calculate(m.list.SelectedItem().(accounts).acc)
			if err != nil {
				log.Fatal("Failed to get code", "error", err)
			}
			m.code = code
			m.state = SHOW_CODE
			m.progressPerc = 0
			return m, tickCmd()
		}
		return m, nil
	case SHOW_CODE:
		return UpdateCode(msg, m)
	}
	return nil, nil
}

func (m model) View() string {
	switch m.state {
	case LIST_ACCOUNTS:
		return ViewList(m)
	case TOUCH_KEY:
		return "Touch your YubiKey!!!"
	case SHOW_CODE:
		return ViewCode(m)
	}
	return ""
}
