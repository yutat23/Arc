//go:build !windows

package main

import "fmt"

// Windows以外のプラットフォーム用のスタブ関数
func showWindowsMessageBox(title, message string, uType uint32) {
	// Windows以外では何もしない（フォールバックが使われる）
	fmt.Printf("\n=== %s ===\n%s\n", title, message)
}
