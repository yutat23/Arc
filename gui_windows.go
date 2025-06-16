//go:build windows

package main

import (
	"syscall"
	"unsafe"
)

// Windows用のメッセージボックス表示関数
func showWindowsMessageBox(title, message string, uType uint32) {
	user32 := syscall.NewLazyDLL("user32.dll")
	messageBoxW := user32.NewProc("MessageBoxW")

	titlePtr, _ := syscall.UTF16PtrFromString(title)
	messagePtr, _ := syscall.UTF16PtrFromString(message)

	messageBoxW.Call(
		0, // hWnd (親ウィンドウなし)
		uintptr(unsafe.Pointer(messagePtr)),
		uintptr(unsafe.Pointer(titlePtr)),
		uintptr(uType),
	)
}
