package ykoath

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"errors"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

const (
	CHALLANG_LEN = 8
)

func (y *YKO) Validate(password []byte) error {
	derivedKey := pbkdf2.Key(password, y.selectData.name, 1000, 16, sha1.New)

	hash := hmac.New(y.selectData.algo.getHashFunc(), derivedKey)
	hash.Write(y.selectData.challange)
	response := make([]byte, hash.Size()+2)
	response[0] = byte(RESPONSE)
	response[1] = byte(hash.Size())
	copy(response[2:], hash.Sum(nil))

	challange := make([]byte, CHALLANG_LEN+2)
	challange[0] = byte(CHALLENGE)
	challange[1] = CHALLANG_LEN
	rand.Read(challange[2:])

	command := make([]byte, len(response)+len(challange))
	copy(command, response)
	copy(command[len(response):], challange)

	res, err := y.send(&sendData{
		CLA:  0x00,
		INS:  VALIDATE,
		P1:   0x00,
		P2:   0x00,
		DATA: command,
	})
	if err != nil {
		return buildError(VALIDATE, err)
	}

	if res[0].tag == RESPONSE {
		hash := hmac.New(y.selectData.algo.getHashFunc(), derivedKey)
		hash.Write(challange[2:])
		if !bytes.Equal(res[0].data, hash.Sum(nil)) {
			return errors.New("validate: failed to validate returned challange")
		}
	} else {
		return fmt.Errorf("validate: Expected response tag")
	}

	return nil
}
