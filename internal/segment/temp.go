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

	// 先 Sync 刷盘，再 Close 句柄，最后 Rename。
	// Close 之后再 Sync 作用在已关闭句柄上无效，
	// Windows 上 Rename 出去的可能是未刷盘内容。
	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	return os.Rename(src, dst)
}
