package main

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
)

func int64ToByteArray(number int64) []byte {
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

func byteArrayToUint64(byteArray []byte) uint64 {
	// Ensure this byte array was produced in a way that matches how you interpret it (big endian or little endian)
	missingBytes := 8 - len(byteArray)

	for i := 0; i < missingBytes; i++ {
		byteArray = append([]byte{0x00}, byteArray...)
	}

	return binary.BigEndian.Uint64(byteArray)
}

func idToKey(id int64) string {
	v := int64ToByteArray(id)
	fmt.Println(v)
	s := base64.RawURLEncoding.EncodeToString(v)
	return s
}

func keyToID(key string) (uint64, error) {
	b, err := base64.RawURLEncoding.DecodeString(key)
	u := byteArrayToUint64(b)

	if err != nil {
		return 0, err
	}
	return u, nil
}
