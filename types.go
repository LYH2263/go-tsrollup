package tsrollup

import "time"

// Sample 是一次指标采样。
type Sample struct {
	Series  string
	Labels  map[string]string
	Ts      time.Time
	Value   float64
}

// Point 是查询返回的聚合点。
type Point struct {
	Start time.Time
	End   time.Time
	Value float64
	Count int64
	Agg   string
}

// WindowResult 是一次窗口查询的结果。
type WindowResult struct {
	Series string
	Labels map[string]string
	Agg    string
	Points []Point
}

// SeriesInfo 描述已注册系列。
type SeriesInfo struct {
	ID     string
	Labels map[string]string
	Samples int64
	Windows int64
}

// Stats 引擎运行统计。
type Stats struct {
	Series       int
	Samples      int64
	WindowsClosed int64
	Compactions  int64
	Segments     int
	WALRecords   int64
	Closed       bool
}
