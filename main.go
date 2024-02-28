package main

import (
	"fmt"
	"os"
	"yko/ykoath"

	"github.com/charmbracelet/log"
	"github.com/ebfe/scard"
	"golang.design/x/clipboard"
	"golang.org/x/term"
)

var yko *ykoath.YKO

func main() {
	ctx, err := scard.EstablishContext()
	if err != nil {
		log.Fatal("Failed to create scard context", "error", err)
	}
	defer ctx.Release()

	readers, err := ctx.ListReaders()
	if err != nil {
		log.Fatal("Failed to list readers", "error", err)
	}
	log.Debug("Found readers", "readers", readers)

	if len(readers) > 1 {
		log.Warn("Multiple readers found, I'll use the first one", "reader", readers[0])
	}

	card, err := ctx.Connect(readers[0], scard.ShareExclusive, scard.ProtocolAny)
	if err != nil {
		log.Fatal("Failed to connect to reader", "error", err)
	}
	defer card.Disconnect(scard.ResetCard)

    if err := clipboard.Init(); err != nil {
        log.Fatal("clipboard", "error", err)
    }

    // YoubiKey CODE
    yko = ykoath.New(card)
    if err := yko.Select(); err != nil {
        log.Fatal(err)
    }
    if yko.AuthRequired {
        fmt.Print("Password: ")
        pass, err := term.ReadPassword(int(os.Stdin.Fd()))
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println()
        if err := yko.Validate(pass); err != nil {
            log.Fatal(err)
        }
    }
    list, err := yko.List()
    if err != nil {
        log.Fatal(err)
    }
    accounts := make([]string, len(list))
    for i := 0; i < len(list); i++ {
        accounts[i] = list[i].Name
    }

    startTUI(accounts)
}
