package main

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
)

func hash(fileName string) (string, error) {
	var buf [0x10000]byte
	var r [256 / 8]byte
	f, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = f.ReadAt(buf[:], 0x100000)
	if err != nil {
		return "", err
	}
	h := sha256.New()
	_, err = h.Write(buf[:])
	if err != nil {
		return "", err
	}
	h.Sum(r[:0])
	return hex.EncodeToString(r[:4]), nil
}
