package validate

import (
	"fmt"
	"time"
)

func Sample(series string, ts time.Time, value float64) error {
	if series == "" {
		return fmt.Errorf("empty series")
	}
	if ts.IsZero() {
		return fmt.Errorf("zero timestamp")
	}
	if value != value { // NaN
		return fmt.Errorf("nan value")
	}
	return nil
}
