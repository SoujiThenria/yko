package ykoath

// Constants based on: https://developers.yubico.com/OATH/YKOATH_Protocol.html

// Instructions
type yubiKeyInstruction uint8

const (
	PUT            yubiKeyInstruction = 0x01
	DELETE         yubiKeyInstruction = 0x02
	SET_CODE       yubiKeyInstruction = 0x03
	RESET          yubiKeyInstruction = 0x04
	LIST           yubiKeyInstruction = 0xa1
	CALCULATE      yubiKeyInstruction = 0xa2
	VALIDATE       yubiKeyInstruction = 0xa3
	CALCULATE_ALL  yubiKeyInstruction = 0xa4
	SELECT         yubiKeyInstruction = 0xa4 // Synthetic, does not realy exists
	SEND_REMAINING yubiKeyInstruction = 0xa5
)

// Algorithms
type yubiKeyAlgo uint8

const (
	HMAC_SHA1   yubiKeyAlgo = 0x01
	HMAC_SHA256 yubiKeyAlgo = 0x02
	HMAC_SHA512 yubiKeyAlgo = 0x03
)

// Types
type yubiKeyType uint8

const (
	HOTP yubiKeyType = 0x10
	TOTP yubiKeyType = 0x20
)

// Properties
type yubiKeyProperty uint8

const (
	ONLY_INCREASING yubiKeyProperty = 0x01
	REQUIRE_TOUCH   yubiKeyProperty = 0x02
)

// Response Codes
type yubiKeyResponse uint16

const (
	RES_SUCCESS                 yubiKeyResponse = 0x9000
	RES_NO_SPACE                yubiKeyResponse = 0x6a84
	RES_AUTH_REQUIRED           yubiKeyResponse = 0x6982
	RES_WRONG_SYNTAX            yubiKeyResponse = 0x6a80
	RES_NO_SUCH_OBJECT          yubiKeyResponse = 0x6984
	RES_RESPONSE_DOES_NOT_MATCH yubiKeyResponse = 0x6984
	RES_MORE_DATA_AVAILABLE     yubiKeyResponse = 0x61
	RES_GENERIC_ERROR           yubiKeyResponse = 0x6581
	RES_AUTH_NOT_ENABLED        yubiKeyResponse = 0x6984
)

// Tags
type yubiKeyTag uint8

const (
	VERSION   yubiKeyTag = 0x79
	NAME      yubiKeyTag = 0x71
	NAME_LIST yubiKeyTag = 0x72
	CHALLENGE yubiKeyTag = 0x74
	ALGORITHM yubiKeyTag = 0x7b
	KEY       yubiKeyTag = 0x73
	PROPERTY  yubiKeyTag = 0x78
	IMF       yubiKeyTag = 0x7a
	RESPONSE  yubiKeyTag = 0x75
)
