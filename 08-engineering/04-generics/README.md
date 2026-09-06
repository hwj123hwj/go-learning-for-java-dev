# 04 · 泛型

Go 1.18 才有泛型，比 Java 晚 17 年，所以 Go 的泛型更克制：**解决"容器和工具函数的类型安全"，不追求抽象一切**。

## 跑

```bash
go run ./04-generics
```

## 练习（答案见 [../SOLUTIONS.md](../SOLUTIONS.md)）

**4.1〔徒手〕** 不看源码实现 `Keys[M ~map[K]V, K comparable, V any](m M) []K`，然后用它取出 `map[string]int` 的所有 key。

**4.2〔预测〕** `Sum([]bool{true, false})` 能编译吗？为什么？把约束换成 `any` 呢？

**4.3〔思考〕** 你网关里有 `map[string]*Provider`、`[]Handler` 这样的容器——什么时候值得为它们写泛型（比如 `Find[T any](xs []T, pred func(T) bool) (T, bool)`），什么时候 `for` 循环直接写更清楚？给自己的结论。
