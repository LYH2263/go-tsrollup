package policy

import "time"

// CompactPolicy 压缩触发策略。
type CompactPolicy struct {
	MinClosedWindows int
	MaxDelay         time.Duration
	SyncWALFirst     bool
}

func DefaultCompactPolicy() CompactPolicy {
	return CompactPolicy{
		MinClosedWindows: 1,
		MaxDelay:         time.Minute,
		SyncWALFirst:     true,
	}
}

func (p CompactPolicy) ShouldCompact(closed int, sinceLast time.Duration) bool {
	if closed >= p.MinClosedWindows {
		return true
	}
	if p.MaxDelay > 0 && sinceLast >= p.MaxDelay && closed > 0 {
		return true
	}
	return false
}
