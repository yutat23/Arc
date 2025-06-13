package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("使い方: go-arch-check file.exe [file2.exe ...]")
		return
	}

	for _, path := range os.Args[1:] {
		arch, err := detectArch(path)
		if err != nil {
			fmt.Printf("%s: エラー - %v\n", path, err)
		} else {
			fmt.Printf("%s: %s\n", path, arch)
		}
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
		return "", fmt.Errorf("PEシグネチャが見つかりません")
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
		return fmt.Sprintf("未知 (0x%X)", code), nil
	}
}
