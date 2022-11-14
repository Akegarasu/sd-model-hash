package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"sync"
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

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("usage: sd-model-hash file1.ckpt file2.ckpt ...")
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(len(os.Args[1:]))
	for _, f := range os.Args[1:] {
		go func(f string) {
			defer wg.Done()
			h, err := hash(f)
			if err != nil {
				fmt.Printf("sd-model-hash (%s) err: %s\n", f, err)
				return
			}
			fmt.Printf("sd-model-hash (%s) = %s\n", f, h)
		}(f)
	}
	wg.Wait()
}
