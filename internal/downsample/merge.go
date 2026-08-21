package downsample

// MergeBuckets 合并同起点桶（用于跨段汇总）。
func MergeBuckets(dst *Bucket, src Bucket) {
	if dst == nil {
		return
	}
	dst.Sum += src.Sum
	dst.Count += src.Count
	if src.set && (!dst.set || src.Max > dst.Max) {
		dst.Max = src.Max
		dst.set = true
	}
	if dst.Start.IsZero() || (!src.Start.IsZero() && src.Start.Before(dst.Start)) {
		dst.Start = src.Start
	}
	if src.End.After(dst.End) {
		dst.End = src.End
	}
}
