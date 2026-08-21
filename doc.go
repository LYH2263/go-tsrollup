// Package tsrollup 提供进程内时序指标滚动聚合引擎。
//
// 写路径：Append 采样点 → tumbling 窗口聚合（sum/avg/max）→ 降采样层 →
// 标签倒排索引 → rollup 段文件持久化。配套 cmd/tsd 管理页查看系列与窗口结果。
//
// 主类型：Engine、Sample、WindowResult。
package tsrollup
