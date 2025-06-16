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

// macOS用のダイアログ表示関数
func showMacDialog(title, message string) {
	// osascriptを使ってAppleScriptでダイアログを表示
	cmd := fmt.Sprintf(`osascript -e 'display dialog "%s" with title "%s" buttons {"OK"} default button "OK"'`, 
		strings.ReplaceAll(message, `"`, `\"`), 
		strings.ReplaceAll(title, `"`, `\"`))
	
	if err := runCommand(cmd); err != nil {
		// フォールバック：標準出力に表示
		fmt.Printf("\n=== %s ===\n%s\n", title, message)
	}
}

// Linux用のダイアログ表示関数
func showLinuxDialog(title, message string) {
	// zenityまたはkdialogを試す
	commands := []string{
		fmt.Sprintf(`zenity --info --title="%s" --text="%s"`, title, message),
		fmt.Sprintf(`kdialog --msgbox "%s" --title "%s"`, message, title),
		fmt.Sprintf(`xmessage -title "%s" "%s"`, title, message),
	}
	
	for _, cmd := range commands {
		if err := runCommand(cmd); err == nil {
			return // 成功したら終了
		}
	}
	
	// フォールバック：標準出力に表示
	fmt.Printf("\n=== %s ===\n%s\n", title, message)
}

// コマンド実行のヘルパー関数
func runCommand(cmdLine string) error {
	parts := strings.Fields(cmdLine)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}
	
	cmd := exec.Command(parts[0], parts[1:]...)
	return cmd.Run()
}
