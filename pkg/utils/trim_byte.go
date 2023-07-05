package utils

import "bytes"

var (
	NEW_LINE_BYTE = []byte{'\n'}
	SPACE_BYTE    = []byte{' '}
)

func TrimByte(msg []byte) []byte {
	return bytes.TrimSpace(bytes.Replace(msg, NEW_LINE_BYTE, SPACE_BYTE, -1))
}
