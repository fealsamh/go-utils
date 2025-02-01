//go:build unix
// +build unix

package mem

import "syscall"

// Alloc allocates the required number of bytes on the heap.
func Alloc(size int) ([]byte, error) {
	return syscall.Mmap(-1, 0, size, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_ANON|syscall.MAP_PRIVATE)
}

// Free frees the slice allocated by [Alloc].
func Free(b []byte) error {
	return syscall.Munmap(b)
}
