//go:build !windows
// +build !windows

package main

import (
	"fmt"
	"os"
	"sync"
)

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
