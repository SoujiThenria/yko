package ykoath

import (
	"fmt"

	"github.com/charmbracelet/log"
)

type ListData struct {
    Name string
    algo byte
}

func (y *YKO) List() ([]ListData, error) {
    resp, err := y.send(&sendData{
        CLA: 0x00,
        INS: INS_LIST,
        P1: 0x00,
        P2: 0x00,
        DATA: nil,
    })

    if err != nil {
        return nil, err
    }
    log.Debug("list", "response", resp)

    data := make([]ListData, len(resp))
    i := 0
    for _,e := range resp {
        switch (e.tag) {
        case TAG_NAME_LIST:
            data[i].Name = string(e.data[1:])
            i++
        default:
            return nil, fmt.Errorf("list: invalid tag: % 02X", e.tag)
        }
    }

    return data, nil
}
