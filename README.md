# arc - PE Architecture Detector

A simple command-line tool written in Go that detects the architecture of Windows PE (Portable Executable) files.

## Features

- Detects architecture of Windows executable files (.exe, .dll, etc.)
- Supports multiple architectures: x86, x64, ARM
- Cross-platform: runs on Windows, macOS, and Linux
- Lightweight and fast
- No external dependencies

## Supported Architectures

- **x86** (32-bit Intel/AMD)
- **x64** (64-bit Intel/AMD)
- **ARM** (ARM and ARM64)
- Reports unknown architectures with their machine code

## Installation

### Download Pre-built Binaries

Pre-built binaries are available in the `build/` directory for the following platforms:

- Windows (AMD64, ARM64)
- macOS (AMD64, ARM64)
- Linux (AMD64, ARM64)

### Build from Source

1. Make sure you have Go 1.24.4 or later installed
2. Clone this repository:
   ```bash
   git clone <repository-url>
   cd arc
   ```
3. Build the binary:
   ```bash
   go build -o arc arc.go
   ```

## Usage

```bash
arc file.exe [file2.exe ...]
```

### Examples

Check a single file:
```bash
arc myapp.exe
```

Check multiple files:
```bash
arc app1.exe app2.dll system32/notepad.exe
```

### Sample Output

```
myapp.exe: x64
legacy.exe: x86
driver.sys: ARM
unknown.exe: Unknown (0x1234)
error.txt: Error - PE signature not found
```

## How It Works

The tool reads the PE header structure of Windows executable files:

1. Reads the DOS header to locate the PE header offset
2. Verifies the PE signature ("PE\0\0")
3. Extracts the Machine field from the COFF header
4. Maps the machine code to the corresponding architecture

## Technical Details

The tool recognizes the following machine codes:
- `0x014c` - IMAGE_FILE_MACHINE_I386 (x86)
- `0x8664` - IMAGE_FILE_MACHINE_AMD64 (x64)
- `0x01c0` - IMAGE_FILE_MACHINE_ARM (ARM)
- `0xaa64` - IMAGE_FILE_MACHINE_ARM64 (ARM64)

## Error Handling

The tool provides descriptive error messages for:
- File access errors
- Invalid PE files
- Corrupted headers

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source. Please check the LICENSE file for details.

## Requirements

- Go 1.24.4 or later (for building from source)
- No runtime dependencies

## Compatibility

- **Input files**: Windows PE format (exe, dll, sys, etc.)
- **Host platforms**: Windows, macOS, Linux (any Go-supported platform)
- **Architectures**: Works on any architecture that Go supports
