package ykoath

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/charmbracelet/log"
	"github.com/ebfe/scard"
)

const (
    INS_PUT = 0x01
    INS_DELETE = 0x02
    INS_SET_CODE = 0x03
    INS_RESET = 0x04
    INS_LIST = 0xa1
    INS_CALCULATE = 0xa2
    INS_VALIDATE = 0xa3
    INS_CALCULATE_ALL = 0xa4
    INS_SEND_REMAINING = 0xa5

    ALG_HMAC_SHA1 = 0x01
    ALG_HMAC_SHA256 = 0x02
    ALG_HMAC_SHA512 =0x03

    TYP_HOTP = 0x10
    TYP_TOTP = 0x20

    PRO_ONLY_INCREASING = 0x01
    PRO_REQUIRE_TOUCH = 0x02

    RES_SUCCESS = 0x9000
    RES_NO_SPACE = 0x6a84
    RES_AUTH_REQUIRED = 0x6982
    RES_WRONG_SYNTAX = 0x6a80
    RES_MORE_DATA_AVAILABLE = 0x61
    RES_NO_SUCH_OBJECT = 0x6984
    RES_RESPONSE_DOES_NOT_MATCH = 0x6984
    RES_GENERIC_ERROR = 0x6581
    RES_AUTH_NOT_ENABLED = 0x6984

    TAG_VERSION = 0x79
    TAG_NAME = 0x71
    TAG_NAME_LIST = 0x72
    TAG_CHALLENGE = 0x74
    TAG_ALGORITHM = 0x7b
    TAG_KEY = 0x73
    TAG_PROPERTY = 0x78
    TAG_IMF = 0x7a
    TAG_RESPONSE = 0x75
)

type YKO struct {
    card *scard.Card
    AuthRequired bool
    sd *selectData
}

type sendData struct {
    CLA byte
    INS byte
    P1 byte
    P2 byte
    DATA []byte
}

type returnData struct {
    tag byte
    data []byte
}

func New(card *scard.Card) (*YKO) {
    return &YKO{card: card}
}

func (y *YKO) send(d *sendData) ([]returnData, error) {
    command := make([]byte, 4+len(d.DATA))
    command[0] = d.CLA
    command[1] = d.INS
    command[2] = d.P1
    command[3] = d.P2
    copy(command[4:], d.DATA)
    log.Debugf("command: % 02X", command)
    resp, err := y.card.Transmit(command)
    if err != nil {
        return nil, err
    }

    for resp[len(resp)-2] == RES_MORE_DATA_AVAILABLE {
        nextResp, err := y.card.Transmit([]byte{0x00, INS_SEND_REMAINING, 0x00, 0x00})
        if err != nil {
            return nil, err
        }
        resp = append(resp[:len(resp)-2], nextResp...)
    }

    log.Debugf("resp: % 02X", resp)
    return parse(resp)
}

func parse(data []byte) ([]returnData, error) {
    test := make([]byte, 2)
    binary.BigEndian.PutUint16(test, RES_SUCCESS)
    if !bytes.Equal(data[len(data) - 2:], test) {
        return nil, errors.New("Some dump return error (TODO)") // TODO: return some error
    }

    result := []returnData{}
    for i := 0; i < len(data[:len(data) - 2]); {
        tag := data[i]
        i++
        dataLen := int(data[i])
        i++
        dataData := data[i:i+dataLen]
        result = append(result, returnData{
            tag: tag,
            data: dataData,
        })
        i+=dataLen
    }
    return result, nil
}
