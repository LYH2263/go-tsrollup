package wal

import (
	"bufio"
	"encoding/json"
	"os"
)

// Replay 重放全部记录。
func (l *Log) Replay() ([]Record, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if err := l.bw.Flush(); err != nil {
		return nil, err
	}
	f, err := os.Open(l.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	var out []Record
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		var rec Record
		if err := json.Unmarshal(line, &rec); err != nil {
			return nil, err
		}
		out = append(out, rec)
	}
	return out, sc.Err()
}
