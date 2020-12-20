package main

import (
	"encoding/binary"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

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

func powerOf2(n uint64) uint64 {
	if n >= 64 {
		panic("integer overflow")
	}
	return 1 << n
}

func toBytes32(x []byte) [32]byte {
	var y [32]byte
	copy(y[:], x)
	return y
}

func copy2dBytes(ary [][]byte) [][]byte {
	if ary != nil {
		copied := make([][]byte, len(ary))
		for i, a := range ary {
			copied[i] = safeCopyBytes(a)
		}
		return copied
	}
	return nil
}


func safeCopyBytes(cp []byte) []byte {
	if cp != nil {
		copied := make([]byte, len(cp))
		copy(copied, cp)
		return copied
	}
	return nil
}

func bytes8(x uint64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, x)
	return bytes
}