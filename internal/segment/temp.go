package segment

import (
	"os"
	"path/filepath"
)

// AtomicReplace 将 src Sync 后 Rename 到 dst（通用工具）。
func AtomicReplace(src, dst string) error {
	f, err := os.OpenFile(src, os.O_RDWR, 0)
	if err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	return os.Rename(src, dst)
}
