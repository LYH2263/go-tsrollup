package fsutil

import (
	"fmt"
	"os"
	"path/filepath"
)

// LockFile 简易锁文件（单进程提示，非跨主机强锁）。
type LockFile struct {
	path string
	f    *os.File
}

func AcquireLock(dir string) (*LockFile, error) {
	if err := EnsureDir(dir); err != nil {
		return nil, err
	}
	path := filepath.Join(dir, "LOCK")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("lock busy: %w", err)
	}
	_, _ = f.WriteString(fmt.Sprintf("pid=%d\n", os.Getpid()))
	_ = f.Sync()
	return &LockFile{path: path, f: f}, nil
}

func (l *LockFile) Release() error {
	if l == nil || l.f == nil {
		return nil
	}
	_ = l.f.Close()
	l.f = nil
	return os.Remove(l.path)
}
