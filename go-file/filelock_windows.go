//go:build windows

// filelock_windows.go
package gofile

func Flock(fd uintptr) error {
	return nil
}

func Funlock(fd uintptr) error {
	return nil
}
