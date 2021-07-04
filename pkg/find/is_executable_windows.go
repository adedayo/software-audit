// +build windows

package find

import (
	"debug/pe"
	"encoding/binary"
	"io/fs"
	"os"
	"path/filepath"
)

func isExecutable(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	var dosheader [96]byte
	if _, err := f.ReadAt(dosheader[0:], 0); err != nil {
		return false
	}

	var base int64
	exe := "NOT-EXE"
	dll := "NOT-DLL"
	// see https://docs.microsoft.com/en-us/windows/win32/debug/pe-format
	if dosheader[0] == 'M' && dosheader[1] == 'Z' {
		signatureOffset := int64(binary.LittleEndian.Uint32(dosheader[0x3c:]))
		var signature [4]byte
		f.ReadAt(signature[:], signatureOffset)

		if signature[0] == 'P' && signature[1] == 'E' && signature[2] == 0 && signature[3] == 0 {
			base = signatureOffset + 4 //4 bytes signature
			characteristicsOffset := base + 18
			var characteristics [2]byte
			f.ReadAt(characteristics[:], characteristicsOffset)
			if (binary.LittleEndian.Uint16(characteristics[:]) & pe.IMAGE_FILE_DLL) == pe.IMAGE_FILE_DLL {
				dll = "DLL"
			}
			if (binary.LittleEndian.Uint16(characteristics[:]) & pe.IMAGE_FILE_EXECUTABLE_IMAGE) == pe.IMAGE_FILE_EXECUTABLE_IMAGE {
				exe = "EXE"
			}

		}

	}

	if dll == "NOT-DLL" && exe == "EXE" {
		return true
	}

	return false
}
