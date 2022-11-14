package main

import (
	"crypto/sha256"
	"fmt"
	"os"
)

func hash(fileName string) string {
	f, err := os.Open(fileName)
	if err != nil {
		return ""
	}
	_, err = f.Seek(0x100000, 0)
	if err != nil {
		return ""
	}
	defer f.Close()
	buf := make([]byte, 0x10000)
	_, err = f.Read(buf)
	if err != nil {
		return ""
	}
	h := sha256.New()
	h.Write(buf)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func main() {
	for _, arg := range os.Args[1:] {
		h := hash(arg)
		if h == "" {
			h = "error"
		} else {
			h = h[:8]
		}
		fmt.Printf("[%s] %s\n", h, arg)
	}
	fmt.Println("Press enter to exit...")
	b := make([]byte, 1)
	os.Stdin.Read(b)
}
