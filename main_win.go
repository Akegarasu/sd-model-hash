//go:build windows
// +build windows

// use go build -ldflags -H=windowsgui
package main

import (
	"fmt"
	"os"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

func showmessagebox(title, content string) {
	toHighDPI()
	_ = boxW(getConsoleWindows(), content, title, 0x00000040|0x00000000)
	return
}

// BoxW of Win32 API. Check https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-messageboxw for more detail.
func boxW(hwnd uintptr, caption, title string, flags uint) int {
	captionPtr, _ := windows.UTF16PtrFromString(caption)
	titlePtr, _ := windows.UTF16PtrFromString(title)
	u32 := windows.NewLazySystemDLL("user32.dll")
	ret, _, _ := u32.NewProc("MessageBoxW").Call(
		hwnd,
		uintptr(unsafe.Pointer(captionPtr)),
		uintptr(unsafe.Pointer(titlePtr)),
		uintptr(flags))

	return int(ret)
}

// GetConsoleWindows retrieves the window handle used by the console associated with the calling process.
func getConsoleWindows() (hWnd uintptr) {
	hWnd, _, _ = windows.NewLazySystemDLL("kernel32.dll").NewProc("GetConsoleWindow").Call()
	return
}

// toHighDPI tries to raise DPI awareness context to DPI_AWARENESS_CONTEXT_UNAWARE_GDISCALED
func toHighDPI() {
	systemAware := ^uintptr(2) + 1
	unawareGDIScaled := ^uintptr(5) + 1
	u32 := windows.NewLazySystemDLL("user32.dll")
	proc := u32.NewProc("SetThreadDpiAwarenessContext")
	if proc.Find() != nil {
		return
	}
	for i := unawareGDIScaled; i <= systemAware; i++ {
		_, _, _ = u32.NewProc("SetThreadDpiAwarenessContext").Call(i)
	}
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
				showmessagebox("错误", fmt.Sprintf("sd-model-hash (%s) err: %s", f, err))
				return
			}
			showmessagebox("成功", fmt.Sprintf("sd-model-hash (%s) = %s", f, h))
		}(f)
	}
	wg.Wait()
}
