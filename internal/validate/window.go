package validate

import (
	"fmt"
	"time"
)

func WindowSize(d time.Duration) error {
	if d <= 0 {
		return fmt.Errorf("non-positive window size")
	}
	if d < time.Millisecond {
		return fmt.Errorf("window too small")
	}
	return nil
}
