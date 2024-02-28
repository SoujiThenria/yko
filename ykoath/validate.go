package ykoath

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"hash"

	"github.com/charmbracelet/log"
	"golang.org/x/crypto/pbkdf2"
)

func (y *YKO) Validate(password []byte) error {
    var h func() hash.Hash
    switch (y.sd.algo[0]) {
    case ALG_HMAC_SHA1:
        h = sha1.New
    case ALG_HMAC_SHA256:
        h = sha256.New
    case ALG_HMAC_SHA512:
        h = sha512.New
    }

    derivedKey := pbkdf2.Key(password, y.sd.name, 1000, 16, sha1.New)
    hash := hmac.New(h, derivedKey)
    hash.Write(y.sd.challange)
    response := make([]byte, 2 + hash.Size())
    response[0] = TAG_RESPONSE
    response[1] = byte(hash.Size())
    copy(response[2:], hash.Sum(nil))

    // tag (1) + size (1) + challange (8) = 10
    challange := make([]byte, 10)
    challange[0] = TAG_CHALLENGE
    challange[1] = byte(8)
    rand.Read(challange[2:])

    command := make([]byte, 1 + len(response) + len(challange))
    command[0] = byte(len(response) + len(challange))
    copy(command[1:], response)
    copy(command[1+len(response):], challange)

    resp, err := y.send(&sendData{
        CLA: 0x00,
        INS: INS_VALIDATE,
        P1: 0x00,
        P2: 0x00,
        DATA: command,
    })

    // A wrong password results into the return code wrong syntax!!!!
    if err != nil {
        return err
    }

    log.Debug("validate", "response", resp)
    if resp[0].tag == TAG_RESPONSE {
        {
            hash := hmac.New(h, derivedKey)
            hash.Write(challange[2:])
            if !bytes.Equal(hash.Sum(nil), resp[0].data) {
                return errors.New("validate: Could not verify the challenge send by the yubikey")
            }
        }
    } else {
        return errors.New("validate: Wrong tag was returned")
    }

    return nil
}
