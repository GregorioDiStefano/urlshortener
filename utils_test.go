package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUInt64ToByteArray(t *testing.T) {
	var i uint64 = 1234567890
	expected := []byte{73, 150, 2, 210}
	assert.Equal(t, expected, uint64ToByteArray(i))
}

func TestByteArrayToUint64(t *testing.T) {
	var i uint64 = 1234567890
	expected := uint64(1234567890)
	assert.Equal(t, expected, byteArrayToUint64(uint64ToByteArray(i)))
}
