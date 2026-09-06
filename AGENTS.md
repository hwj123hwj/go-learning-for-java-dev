# go-learning-for-java-dev · 仓库约定

个人 Java→Go 后端学习仓库。主分支 main，全部为可直接运行的示例代码。

## 结构

- 一级目录 = 学习阶段，序号即顺序（`01-quickstart` … `09-streaming-proxy`）；`interview-prep/` 是面试准备支线，与主线解耦。
- 每个一级目录是独立 Go module（`01-quickstart` 除外，无 go.mod，单文件 `go run xx.go` 运行）。
- 部分旧目录的 module 名与目录名不一致（如 `06-concurrency` 的 module 名仍是 `concurrency-demo`），历史原因，能正常构建，勿"顺手统一"，改动它会连带全量 import 更新。

## 代码风格（新代码必须遵守）

- 示例代码风格：每个 demo 是一个可独立运行的 `main.go`，用 `fmt.Println("--- 小节 ---")` 分节，结尾保留 `compareWithJava()` 小节（Java 代码片段 + Go 对应写法 + "Go惯用法"要点）。
- 教学注释只写"为什么"，不写"下一行做什么"。
- 有练习的模块：`EXERCISES.md` 出题（预测输出/改造/徒手实现/找 bug 四类），答案在 `SOLUTIONS.md`，两者目录一一对应。

## 改动规则

- 不破坏其他模块的构建；提交前在改动过的 module 内跑 `go build ./...`，有测试的跑 `go test ./...`。
- 教学代码里允许为演示而生的简化（如 sleep 模拟延迟），不要"修复"它们，除非练习题答案需要。
- `06-concurrency`、`07-backend-patterns`、`08-engineering`、`09-streaming-proxy` 是含金量最高的模块，改动需保持练习题与代码同步。
- `.trae/specs/` 为 2026-06 规划存档，只读，不更新；学习路线的现行事实来源是根 README.md。
