package mobile

import (
	"io"

	"github.com/fealsamh/go-utils/nocopy"
	"golang.org/x/mobile/asset"
)

// GetFileData returns the contents of a resource file as a slice of bytes.
func GetFileData(filename string) ([]byte, error) {
	f, err := asset.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

// GetFileString returns the contents of a resource file as a string.
func GetFileString(filename string) (string, error) {
	b, err := GetFileData(filename)
	if err != nil {
		return "", err
	}
	return nocopy.String(b), nil
}
