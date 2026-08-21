# go-tsrollup

进程内时序指标滚动聚合：Append 采样点 → tumbling 窗口（sum/avg/max）→ 降采样层 → 标签倒排 → rollup 段持久化。配套 `tsd` 管理页。

## 构建与测试

```bash
go build ./...
go test ./... -count=1
```

## 管理守护进程

```bash
go run ./cmd/tsd -addr :8113 -web web
```

打开 http://127.0.0.1:8113/ 可写入试点并查询窗口聚合。

## 库能力摘要

- `Engine.Append` / `AppendContext`：写入采样（Labels 深拷贝）
- `Engine.Query`：窗口聚合查询（Points 独立拷贝）
- `Engine.Compact`：关闭到期窗口并刷段
- `Engine.Close`：先刷末窗再释放内存
- 内部包：`agg` / `downsample` / `labelidx` / `segment` / `wal`
