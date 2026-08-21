package audit

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Event 审计事件。
type Event struct {
	Time   time.Time `json:"time"`
	Action string    `json:"action"`
	Series string    `json:"series,omitempty"`
	Detail string    `json:"detail,omitempty"`
}

// Logger 追加写审计日志。
type Logger struct {
	mu sync.Mutex
	f  *os.File
}

func Open(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	return &Logger{f: f}, nil
}

func (l *Logger) Log(action, series, detail string) error {
	if l == nil || l.f == nil {
		return nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	b, err := json.Marshal(Event{Time: time.Now().UTC(), Action: action, Series: series, Detail: detail})
	if err != nil {
		return err
	}
	if _, err := l.f.Write(append(b, '\n')); err != nil {
		return err
	}
	return nil
}

func (l *Logger) Close() error {
	if l == nil || l.f == nil {
		return nil
	}
	err := l.f.Close()
	l.f = nil
	return err
}
