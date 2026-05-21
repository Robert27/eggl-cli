package dedash

import (
	"path/filepath"
	"strings"
	"unicode/utf8"
)

var skippedDirs = map[string]struct{}{
	".git":         {},
	"node_modules": {},
	"vendor":       {},
	"bin":          {},
	"dist":         {},
	"build":        {},
}

var skippedExtensions = map[string]struct{}{
	".png":  {},
	".jpg":  {},
	".jpeg": {},
	".gif":  {},
	".webp": {},
	".ico":  {},
	".bmp":  {},
	".exe":  {},
	".dll":  {},
	".so":   {},
	".dylib": {},
	".a":    {},
	".o":    {},
	".zip":  {},
	".tar":  {},
	".gz":   {},
	".bz2":  {},
	".xz":   {},
	".7z":   {},
	".rar":  {},
	".woff": {},
	".woff2": {},
	".ttf":  {},
	".otf":  {},
	".eot":  {},
	".pdf":  {},
	".wasm": {},
	".pyc":  {},
	".class": {},
	".bin":  {},
}

func shouldSkipDir(name string) bool {
	_, ok := skippedDirs[name]
	return ok
}

func shouldSkipExtension(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	_, ok := skippedExtensions[ext]
	return ok
}

func isBinaryContent(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	if strings.Contains(string(data), "\x00") {
		return true
	}
	return !utf8.Valid(data)
}
