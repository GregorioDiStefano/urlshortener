package main

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// int64ToByteArray converts an int64 to a byte array and removed leading 0x00s
// For example, a number such as 1020, when converted to a byte array, will be [3 252] instead of [0 0 0 0 0 0 3 252]
func uint64ToByteArray(number int64) []byte {
	byteArray := make([]byte, 8)
	binary.BigEndian.PutUint64(byteArray, uint64(number))

	var idx int
	var v byte
	for idx, v = range byteArray {
		if v != 0x00 {
			break
		}
	}

	return byteArray[idx:]
}

// byteArrayToUint64 converts a byte array to a uint64 and add restores leading 0x00s
func byteArrayToUint64(byteArray []byte) uint64 {
	// Ensure this byte array was produced in a way that matches how you interpret it (big endian or little endian)
	missingBytes := 8 - len(byteArray)

	for i := 0; i < missingBytes; i++ {
		byteArray = append([]byte{0x00}, byteArray...)
	}

	return binary.BigEndian.Uint64(byteArray)
}

// uint64ToBase64 converts a int64 to a base64 encoded string
func uint64ToBase64(id int64) string {
	v := uint64ToByteArray(id)
	s := base64.RawURLEncoding.EncodeToString(v)
	return s
}

// base64StringToUint64 converts a base64 encoded string to a uint64
func base64StringToUint64(key string) (uint64, error) {
	b, err := base64.RawURLEncoding.DecodeString(key)
	u := byteArrayToUint64(b)

	if err != nil {
		return 0, err
	}
	return u, nil
}

func shortURLKeyToIDAndNonce(key string) (uint64, string, error) {
	// split short url into two parts: actual id component based on row id, and nonce
	if len(key) < 4 {
		return 0, "", fmt.Errorf("invalid key, too short")
	}

	id := key[0 : len(key)-2]
	nonce := key[len(key)-2:]

	if dbID, err := base64StringToUint64(id); err != nil {
		log.WithError(err).WithField("id", id).Error("error converting id")
		return 0, "", err
	} else {
		return dbID, nonce, nil
	}
}
