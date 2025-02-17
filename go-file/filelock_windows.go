// filelock_windows.go
//go:build windows

package gofile

func Flock(fd uintptr) error {
	h := syscall.Handle(fd)
	var ol syscall.Overlapped

	err := syscall.LockFileEx(h,
		syscall.LOCKFILE_EXCLUSIVE_LOCK|syscall.LOCKFILE_FAIL_IMMEDIATELY,
		0, 0xffffffff, 0xffffffff, &ol)
	return err
}

func Funlock(fd uintptr) error {
	h := syscall.Handle(fd)
	var ol syscall.Overlapped

	err := syscall.UnlockFileEx(h, 0, 0xffffffff, 0xffffffff, &ol)
	return err
}
