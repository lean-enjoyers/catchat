package main

import "bytes"

var (
	newLineByte = []byte{'\n'}
	spaceByte   = []byte{' '}
)

func trimByte(msg []byte) []byte {
	return bytes.TrimSpace(bytes.Replace(msg, newLineByte, spaceByte, -1))
}
