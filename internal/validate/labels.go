package validate

import (
	"fmt"
	"unicode/utf8"
)

const maxLabelPairs = 64
const maxLabelLen = 256

func Labels(labels map[string]string) error {
	if len(labels) > maxLabelPairs {
		return fmt.Errorf("too many labels")
	}
	for k, v := range labels {
		if k == "" {
			return fmt.Errorf("empty label key")
		}
		if !utf8.ValidString(k) || !utf8.ValidString(v) {
			return fmt.Errorf("invalid utf8 in labels")
		}
		if len(k) > maxLabelLen || len(v) > maxLabelLen {
			return fmt.Errorf("label too long")
		}
	}
	return nil
}
