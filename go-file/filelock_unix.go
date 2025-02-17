//go:build !windows

// filelock_unix.go
package gofile

import "syscall"

func Flock(fd uintptr) error {
	return syscall.Flock(int(fd), syscall.LOCK_EX)
}

func Funlock(fd uintptr) error {
	return syscall.Flock(int(fd), syscall.LOCK_UN)
}
