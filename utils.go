package main

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func getRequiredEnvVar(key string) string {
	value := os.Getenv(key)
	if value == "" {
		logrus.Fatalf("Environment variable %s is required", key)
	}

	return value
}

// int64ToByteArray converts an int64 to a byte array and removed leading 0x00s
// For example, a number such as 1020, when converted to a byte array, will be [3 252] instead of [0 0 0 0 0 0 3 252]
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
	v := int64ToByteArray(id)
	fmt.Println(v)
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
