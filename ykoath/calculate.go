package ykoath

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
)

func (y *YKO) Calculate(name string) (string, error) {

    name_part := make([]byte, 2+len(name))
    name_part[0] = TAG_NAME
    name_part[1] = byte(len(name))
    copy(name_part[2:], name)

    challange := make([]byte, 2+8)
    challange[0] = TAG_CHALLENGE
    challange[1] = byte(8)
    timestamp := time.Now().Unix() / 30
    binary.BigEndian.PutUint64(challange[2:], uint64(timestamp))

    command := make([]byte, 1 + len(name_part) + len(challange))
    command[0] = byte(len(name_part) + len(challange))
    copy(command[1:], name_part)
    copy(command[1+len(name_part):], challange)

    resp, err := y.send(&sendData{
        CLA: 0x00,
        INS: INS_CALCULATE,
        P1: 0x00,
        P2: 0x01,
        DATA: command,
    })
    if err != nil {
        return "", err
    }

    log.Debug("Calculate", "response", resp)

    for _,e := range resp {
        switch (e.tag) {
        case 0x76:
            dig := int(e.data[0])
            code := binary.BigEndian.Uint32(e.data[1:])
            return fmt.Sprintf(fmt.Sprintf("%%0%dd", dig), code), nil
        default:
            return "", fmt.Errorf("list: invalid tag: % 02X", e.tag)
        }
    }
    return "", nil
}
