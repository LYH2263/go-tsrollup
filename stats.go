package tsrollup

// Stats 返回运行统计快照。
func (e *Engine) Stats() Stats {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return Stats{
		Series:        len(e.series),
		Samples:       e.samples.Load(),
		WindowsClosed: e.closedW.Load(),
		Compactions:   e.compacts.Load(),
		Segments:      len(e.readonly),
		WALRecords:    e.walRecs.Load(),
		Closed:        e.closed.Load(),
	}
}
