package ykoath

import (
	"encoding/binary"
	"fmt"
	"time"
)

func (y *YKO) Calculate(name string) (string, error) {
    name_data := make([]byte, len(name) + 2)
    name_data[0] = byte(NAME)
    name_data[1] = byte(len(name))
    copy(name_data[2:], []byte(name))

    challange := make([]byte, CHALLANG_LEN + 2)
    challange[0] = byte(CHALLENGE)
    challange[1] = CHALLANG_LEN
    binary.BigEndian.PutUint64(challange[2:], uint64(time.Now().Unix() / 30))

    command := make([]byte, len(name_data) + len(challange))
    copy(command, name_data)
    copy(command[len(name_data):], challange)

	res, err := y.send(&sendData{
		CLA:  0x00,
		INS:  CALCULATE,
		P1:   0x00,
		P2:   0x01,
		DATA: command,
	})
	if err != nil {
		return "", buildError(CALCULATE, err)
	}

    for _,e := range res {
        switch (e.tag) {
        case 0x76:
            dig := int(e.data[0])
            code := binary.BigEndian.Uint32(e.data[1:])
            return fmt.Sprintf(fmt.Sprintf("%%0%dd", dig), code), nil
        default:
            return "", fmt.Errorf("list: invalid tag: % 02X", e.tag)
        }
    }

    return "", fmt.Errorf("calgulate: something went wrong")
}
