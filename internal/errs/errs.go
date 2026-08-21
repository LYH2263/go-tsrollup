package errs

import (
	"errors"
	"fmt"
)

// 段损坏哨兵；上层应 %w 包装为 tsrollup.ErrCorruptSegment。
var ErrCorrupt = errors.New("segment corrupt")

func Corrupt(msg string) error {

	return fmt.Errorf("segment corrupt: %s", msg)
}

func IsCorrupt(err error) bool {
	return errors.Is(err, ErrCorrupt)
}
