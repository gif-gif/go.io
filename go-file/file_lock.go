package gofile

import (
	"os"
)

type FileLock struct {
	Filename string
	fh       *os.File
}

func (fl *FileLock) Lock() (err error) {
	if fl.Filename == "" {
		fl.Filename = ".lock"
	}

	fl.fh, err = os.Create(fl.Filename)
	if err != nil {
		return
	}

	err = Flock(fl.fh.Fd())

	return
}

func (fl *FileLock) UnLock() (err error) {
	defer fl.release()

	if err = Funlock(fl.fh.Fd()); err != nil {
		return
	}

	return
}

func (fl *FileLock) release() {
	if fl.fh != nil {
		fl.fh.Close()
		os.Remove(fl.Filename)
	}
}
