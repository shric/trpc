package fileutils

import (
	"os"
	"path/filepath"
)

// IsDirectory returns true if path is a directory, false if not (or on error).
func IsDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

// RealPath makes filename an absolute path with symlinks evaluated if possible.
func RealPath(filename string) string {
	abs, err := filepath.Abs(filename)

	if err == nil {
		filename = abs
	}

	symLinked, err := filepath.EvalSymlinks(filename)
	if err == nil {
		filename = symLinked
	}

	return filename
}
