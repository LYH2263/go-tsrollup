package wal

// Flush 强制刷盘。
func (l *Log) Flush() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.f == nil {
		return nil
	}
	if err := l.bw.Flush(); err != nil {
		return err
	}
	return l.f.Sync()
}

// SetSyncEvery 调整同步频率。
func (l *Log) SetSyncEvery(n int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if n > 0 {
		l.syncEvery = n
	}
}
