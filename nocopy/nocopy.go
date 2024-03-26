package nocopy

import "unsafe"

// String converts a slice of bytes into a string.
func String(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// Bytes converts a string into a slice of bytes.
func Bytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
