package wal

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Log 简易写前日志（非 LSM）。
type Log struct {
	mu        sync.Mutex
	dir       string
	path      string
	f         *os.File
	bw        *bufio.Writer
	n         int
	syncEvery int
}

func Open(dir string) (*Log, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	path := filepath.Join(dir, "current.wal")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}
	return &Log{
		dir:       dir,
		path:      path,
		f:         f,
		bw:        bufio.NewWriter(f),
		syncEvery: 32,
	}, nil
}

// Append 写入记录；尊重调用方 ctx，序列化、刷盘前先做取消检查，
// 避免取消后仍在 fsync 上空转。
func (l *Log) Append(ctx context.Context, rec Record) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.f == nil {
		return fmt.Errorf("wal: closed")
	}
	// 拿到锁后再次确认 ctx 仍在有效期内——排队期间可能已被取消。
	if err := ctx.Err(); err != nil {
		return err
	}
	b, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	if _, err := l.bw.Write(b); err != nil {
		return err
	}
	if err := l.bw.WriteByte('\n'); err != nil {
		return err
	}
	l.n++
	if l.n%l.syncEvery == 0 {
		// fsync 最为耗时，刷盘前再确认一次取消，避免取消后仍空转。
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := l.bw.Flush(); err != nil {
			return err
		}
		if err := l.f.Sync(); err != nil {
			return err
		}
	}
	return nil
}

func (l *Log) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.f == nil {
		return nil
	}
	_ = l.bw.Flush()
	err := l.f.Close()
	l.f = nil
	return err
}

func (l *Log) Path() string { return l.path }
