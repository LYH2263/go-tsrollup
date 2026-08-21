package errs

import (
	"errors"
	"fmt"
)

// Wrap 若 err 非空则包装哨兵。
func Wrap(sentinel, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%w: %v", sentinel, err)
}

// AsCorrupt 将底层损坏错误提升为带哨兵错误。
func AsCorrupt(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, ErrCorrupt) {
		return err
	}
	return fmt.Errorf("%w: %v", ErrCorrupt, err)
}
