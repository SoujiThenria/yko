package ykoath

type selectData struct {
	version   []byte
	name      []byte
	challange []byte
	algo      yubiKeyAlgo
}

func (y *YKO) Select() error {
	resp, err := y.send(&sendData{
		CLA:  0x00,
		INS:  SELECT,
		P1:   0x04,
		P2:   0x00,
		DATA: []byte{0xa0, 0x00, 0x00, 0x05, 0x27, 0x21, 0x01},
	})
	if err != nil {
		return buildError(SELECT, err)
	}

	for _, e := range resp {
		switch e.tag {
		case VERSION:
			y.selectData.version = e.data
		case NAME:
			y.selectData.name = e.data
		case CHALLENGE:
			y.selectData.challange = e.data
		case ALGORITHM:
			y.selectData.algo = yubiKeyAlgo(e.data[0])
		}
	}

	return nil
}
