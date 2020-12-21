package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// File helpers.

func expandPath(p string) (string, error) {
	if strings.HasPrefix(p, "~/") || strings.HasPrefix(p, "~\\") {
		if home := homeDir(); home != "" {
			p = home + p[1:]
		}
	}
	return filepath.Abs(filepath.Clean(os.ExpandEnv(p)))
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

// Bytes helpers.

func stringArrayFrom2DBytes(bytes2d [][]byte) []string {
	stringArray := make([]string, len(bytes2d))
	for i, bytes := range bytes2d {
		stringArray[i] = fmt.Sprintf("%#x", bytes)
	}
	return stringArray
}

func uint64To256BytesBigEndian(i uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return common.LeftPadBytes(buf, 32)
}