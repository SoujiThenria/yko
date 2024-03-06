package ykoath

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"fmt"
	"hash"

	"github.com/ebfe/scard"
)

type YKO struct {
	card       *scard.Card
	selectData selectData
}

type sendData struct {
	CLA  byte
	INS  yubiKeyInstruction
	P1   byte
	P2   byte
	DATA []byte
}

type resultData struct {
	tag  yubiKeyTag
	data []byte
}

type yubiKeyError struct {
	msg  string
	code yubiKeyResponse
}

func (e *yubiKeyError) Error() string { return e.msg }

func New(c *scard.Card) *YKO {
	return &YKO{card: c}
}

func (y *YKO) AuthRequired() bool { return len(y.selectData.challange) > 0 }

func (y *YKO) send(d *sendData) ([]resultData, error) {
	command := make([]byte, 5+len(d.DATA))
	command[0] = d.CLA
	command[1] = byte(d.INS)
	command[2] = d.P1
	command[3] = d.P2
	command[4] = byte(len(d.DATA))
	copy(command[5:], d.DATA)

	res, err := y.card.Transmit(command)
	if err != nil {
		return nil, &yubiKeyError{msg: err.Error(), code: 0}
	}

	for yubiKeyResponse(res[len(res)-2]) == RES_MORE_DATA_AVAILABLE {
		moreRes, err := y.card.Transmit([]byte{0x00, byte(SEND_REMAINING), 0x00, 0x00})
		if err != nil {
			return nil, &yubiKeyError{msg: err.Error(), code: 0}
		}
		res = append(res[:len(res)-2], moreRes...)
	}

	resp_code := binary.BigEndian.Uint16(res[len(res)-2:])
	switch yubiKeyResponse(resp_code) {
	case RES_SUCCESS:
		return parse(res[:len(res)-2]), nil
	case RES_NO_SPACE:
		return nil, &yubiKeyError{msg: "The YubiKey returned an error", code: RES_NO_SPACE}
	case RES_AUTH_REQUIRED:
		return nil, &yubiKeyError{msg: "The YubiKey returned an error", code: RES_AUTH_REQUIRED}
	case RES_WRONG_SYNTAX:
		return nil, &yubiKeyError{msg: "The YubiKey returned an error", code: RES_WRONG_SYNTAX}
	case RES_NO_SUCH_OBJECT: // can also be RES_AUTH_NOT_ENABLED or RES_RESPONSE_DOES_NOT_MATCH
		return nil, &yubiKeyError{msg: "The YubiKey returned an error", code: RES_NO_SUCH_OBJECT}
	case RES_GENERIC_ERROR:
		return nil, &yubiKeyError{msg: "The YubiKey returned an error", code: RES_GENERIC_ERROR}
	}

	return nil, &yubiKeyError{msg: "Some thing went wrong, unidentified response code", code: yubiKeyResponse(resp_code)}
}

func parse(data []byte) []resultData {
	result := []resultData{}
	for i := 0; i < len(data); {
		tag := data[i]
		i++
		dataLen := int(data[i])
		i++
		dataData := data[i : i+dataLen]
		result = append(result, resultData{
			tag:  yubiKeyTag(tag),
			data: dataData,
		})
		i += dataLen
	}
	return result
}

func buildError(inst yubiKeyInstruction, err error) error {
	// TODO: make this global ???
	var ErrorMessages = map[yubiKeyResponse]string{
		RES_SUCCESS:        "SUCCESS",
		RES_NO_SPACE:       "NO SPACE",
		RES_AUTH_REQUIRED:  "AUTH REQUIRED",
		RES_WRONG_SYNTAX:   "WRONG SYNTAX",
		RES_NO_SUCH_OBJECT: "NO SUCH OBJECT",
		// RES_RESPONSE_DOES_NOT_MATCH: "RESPONSE DOES NOT MATCH",
		// RES_AUTH_NOT_ENABLED:        "AUTH NOT ENABLED",
		RES_MORE_DATA_AVAILABLE: "MORE DATA AVAILABLE",
		RES_GENERIC_ERROR:       "GENERIC ERROR",
	}
	switch err.(type) {
	case *yubiKeyError:
		// TODO: fix
		test := err.(*yubiKeyError)
		errMsg := ErrorMessages[test.code]
		if errMsg == "" {
			errMsg = fmt.Sprintf("% 02X", test.code)
		}

		fnName := ""
		if test.code == RES_NO_SUCH_OBJECT {
			switch inst {
			case VALIDATE:
				errMsg = "AUTH NOT ENABLED"
			case SET_CODE:
				errMsg = "RESPONSE DOES NOT MATCH"
			}
		}
		if test.code == RES_WRONG_SYNTAX && inst == VALIDATE {
			errMsg = "WRONG PASSWORD"
		}

		switch inst {
		case PUT:
			fnName = "PUT"
		case DELETE:
			fnName = "DELETE"
		case SET_CODE:
			fnName = "SET_CODE"
		case RESET:
			fnName = "RESET"
		case LIST:
			fnName = "LIST"
		case CALCULATE:
			fnName = "CALCULATE (or SELECT)"
		case VALIDATE:
			fnName = "VALIDATE"
		case CALCULATE_ALL:
			fnName = "CALCULATE_ALL"
		case SEND_REMAINING:
			fnName = "SEND_REMAINING"
		}
		return fmt.Errorf("%s: %s: %s", fnName, test.msg, errMsg)
	case nil:
		return nil
	}
	return err
}

func (a yubiKeyAlgo) getHashFunc() func() hash.Hash {
	switch a {
	case HMAC_SHA1:
		return sha1.New
	case HMAC_SHA256:
		return sha256.New
	case HMAC_SHA512:
		return sha512.New
	default:
		return sha1.New
	}
}
