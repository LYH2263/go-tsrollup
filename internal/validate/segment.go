package validate

import (
	"fmt"
	"os"
)

func SegmentFile(path string) error {
	st, err := os.Stat(path)
	if err != nil {
		return err
	}
	if st.IsDir() {
		return fmt.Errorf("segment path is directory")
	}
	if st.Size() < 16 {
		return fmt.Errorf("segment too small")
	}
	return nil
}
