package mobile

import (
	"io"

	"golang.org/x/mobile/asset"
)

// GetFileData returns the contents of a resource file.
func GetFileData(filename string) ([]byte, error) {
	f, err := asset.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}
