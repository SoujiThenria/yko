package ykoath

import (
	"fmt"
)

type ListData struct {
	Name string
	Algo yubiKeyAlgo
	Type yubiKeyType
}

func (y *YKO) List() ([]ListData, error) {
	resp, err := y.send(&sendData{
		CLA:  0x00,
		INS:  LIST,
		P1:   0x00,
		P2:   0x00,
		DATA: nil,
	})

	if err != nil {
		return nil, buildError(LIST, err)
	}

	data := make([]ListData, len(resp))
	i := 0
	for _, e := range resp {
		switch e.tag {
		case NAME_LIST:
			// High 4 bits is type, low 4 bits is algorithm
			data[i].Algo = yubiKeyAlgo((e.data[0] & 0x0F))
			data[i].Type = yubiKeyType((e.data[0] & 0xF0))
			data[i].Name = string(e.data[1:])
		default:
			return nil, fmt.Errorf("list: invalid tag: % 02X", e.tag)
		}
		i++
	}

	return data, nil
}
