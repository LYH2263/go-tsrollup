package agg

// MustNew 创建聚合器，未知 kind 时 panic（仅测试辅助；生产用 Factory.New）。
func MustNew(kind string) Aggregator {
	a, err := DefaultFactory{}.New(kind)
	if err != nil {
		panic(err)
	}
	return a
}

// Kinds 返回内置聚合名。
func Kinds() []string {
	return []string{"sum", "avg", "max"}
}

// IsKind 判断是否支持。
func IsKind(kind string) bool {
	switch kind {
	case "sum", "avg", "max", "":
		return true
	default:
		return false
	}
}
