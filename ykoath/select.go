package ykoath

import "github.com/charmbracelet/log"

type selectData struct {
    version []byte
    name []byte
    challange []byte
    algo []byte
}

func (y *YKO) Select() error {
    resp, err := y.send(&sendData{
        CLA: 0x00,
        INS: INS_CALCULATE_ALL, // For what ever reasons its the same
        P1: 0x04,
        P2: 0x00,
        DATA: []byte{0x07, 0xa0, 0x00, 0x00, 0x05, 0x27, 0x21, 0x01},
    })
    if err != nil {
        return err
    }

    log.Debug("Select", "response", resp)
    sd := selectData{}
    y.sd = &sd
    for _,e := range resp {
        switch (e.tag) {
        case TAG_VERSION:
            sd.version = e.data
        case TAG_NAME:
            sd.name = e.data
        case TAG_CHALLENGE:
            sd.challange = e.data
            y.AuthRequired = true
        case TAG_ALGORITHM:
            sd.algo = e.data
        }
    }

    return nil
}
