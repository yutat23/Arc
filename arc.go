package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	args := os.Args[1:]
	useGUI := false
	var filePaths []string

	// コマンドライン引数を解析
	for _, arg := range args {
		if arg == "-g" || arg == "--gui" {
			useGUI = true
		} else {
			filePaths = append(filePaths, arg)
		}
	}
	if len(filePaths) == 0 {
		if useGUI {
			showGUIMessage("arc", "Usage: arc [-g] file.exe [file2.exe ...]\n\nDrag and drop executable files or specify them as command line arguments.")
		} else {
			fmt.Println("Usage: arc [-g|--gui] file.exe [file2.exe ...]")
		}
		return
	}

	var results []string
	for _, path := range filePaths {
		arch, err := detectArch(path)
		if err != nil {
			result := fmt.Sprintf("%s: Error - %v", path, err)
			results = append(results, result)
			if !useGUI {
				fmt.Println(result)
			}
		} else {
			result := fmt.Sprintf("%s: %s", path, arch)
			results = append(results, result)
			if !useGUI {
				fmt.Println(result)
			}
		}
	}
	// GUIモードの場合、結果をダイアログで表示
	if useGUI {
		title := "arc"
		message := strings.Join(results, "\n")
		showGUIMessage(title, message)
	}
}

func detectArch(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// DOSヘッダーの最後にPEヘッダーのオフセットがある（0x3C）
	dosHeader := make([]byte, 64)
	if _, err := file.Read(dosHeader); err != nil {
		return "", err
	}
	peOffset := binary.LittleEndian.Uint32(dosHeader[0x3C:])

	// PEヘッダーの位置へ
	if _, err := file.Seek(int64(peOffset), 0); err != nil {
		return "", err
	}

	// "PE\0\0" シグネチャ確認
	signature := make([]byte, 4)
	if _, err := file.Read(signature); err != nil {
		return "", err
	}
	if string(signature) != "PE\x00\x00" {
		return "", fmt.Errorf("PE signature not found")
	}

	// Machine フィールド（2バイト）
	machine := make([]byte, 2)
	if _, err := file.Read(machine); err != nil {
		return "", err
	}
	code := binary.LittleEndian.Uint16(machine)

	switch code {
	case 0x014c:
		return "x86", nil
	case 0x8664:
		return "x64", nil
	case 0x01c0, 0xaa64:
		return "ARM", nil
	default:
		return fmt.Sprintf("Unknown (0x%X)", code), nil
	}
}

// クロスプラットフォーム対応のGUIメッセージ表示関数
func showGUIMessage(title, message string) {
	switch runtime.GOOS {
	case "windows":
		showWindowsMessageBox(title, message, 0x40) // MB_ICONINFORMATION
	case "darwin": // macOS
		showMacDialog(title, message)
	case "linux":
		showLinuxDialog(title, message)
	default:
		// フォールバック：標準出力に表示
		fmt.Printf("\n=== %s ===\n%s\n", title, message)
	}
}

// コマンド実行のヘルパー関数
func runCommandArgs(cmdName string, args ...string) error {
	cmd := exec.Command(cmdName, args...)
	return cmd.Run()
}

// macOS用のダイアログ表示関数
func showMacDialog(title, message string) {
	// AppleScriptでダイアログを表示（text型で渡す）
	script := fmt.Sprintf(`display dialog %q with title %q buttons {"OK"} default button "OK"`, message, title)
	if err := runCommandArgs("osascript", "-e", script); err != nil {
		// フォールバック：標準出力に表示
		fmt.Printf("\n=== %s ===\n%s\n", title, message)
	}
}

// Linux用のダイアログ表示関数
func showLinuxDialog(title, message string) {
	// 利用可能なダイアログツールを自動検出
	candidates := []struct {
		cmd  string
		args []string
	}{
		{"zenity", []string{"--info", "--title", title, "--text", message}},
		{"kdialog", []string{"--msgbox", message, "--title", title}},
		{"xmessage", []string{"-title", title, message}},
	}
	for _, c := range candidates {
		if _, err := exec.LookPath(c.cmd); err == nil {
			if err := runCommandArgs(c.cmd, c.args...); err == nil {
				return // 成功したら終了
			}
		}
	}
	// フォールバック：標準出力に表示
	fmt.Printf("\n=== %s ===\n%s\n", title, message)
}
