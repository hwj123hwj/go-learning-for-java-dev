package main

import (
	"errors"
	"fmt"
)

// ─── 1. 哨兵错误：包级错误变量，调用方用 errors.Is 匹配 ──────────────────────
// 对比 Java：异常类型即分类（NotFoundException e）。Go 没有异常，
// 错误就是一个普通的 error 值，"类型"就是这些导出的哨兵变量。
// 你网关的 internal/handlers/errors.go 用的正是这套。
var (
	ErrNotFound = errors.New("resource not found")
	ErrForbidden = errors.New("forbidden")
)

// ─── 3. 自定义错误类型：错误也能携带结构化字段 ────────────────────────────────
// 对比 Java：自定义 Exception 带 field/getter。Go 里 errors.As 取出原始类型读字段。
type ValidationError struct {
	Field string
	Msg   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %s: %s", e.Field, e.Msg)
}

// ─── 错误产生方：包装而不是替换 ───────────────────────────────────────────────
// %w 保留底层错误（错误链还在），%v/%s 会把底层错误压扁成字符串（链断了）。
// 包装时补充"我当时在干什么"，调用方日志才有上下文。
func findUser(id int) error {
	if id <= 0 {
		return fmt.Errorf("findUser(id=%d): %w", id, ErrNotFound)
	}
	if id == 13 {
		return ErrForbidden // 不包装直接返回也合法，只是少了上下文
	}
	return nil
}

func validateAge(age int) error {
	if age < 0 || age > 150 {
		return &ValidationError{Field: "age", Msg: fmt.Sprintf("out of range: %d", age)}
	}
	return nil
}

func main() {
	fmt.Println("=== Go 错误处理：错误是值 ===")
	fmt.Println()

	// 1) errors.Is：穿透 %w 包装链判断"根因是不是它"
	//    对应 Java 的 catch (NotFoundException e) 按类型捕获
	err := findUser(-1)
	fmt.Println("err                 =", err)
	fmt.Println("errors.Is(ErrNotFound) =", errors.Is(err, ErrNotFound))
	fmt.Println()

	// 2) errors.As：从错误链里取出自定义类型，读它的字段
	//    对应 Java 的 e.getField()
	verr := validateAge(200)
	var ve *ValidationError
	if errors.As(verr, &ve) {
		fmt.Printf("errors.As 取出字段   = field:%s msg:%s\n", ve.Field, ve.Msg)
	}
	fmt.Println()

	// 3) errors.Join：一个操作产生多个错误（Go 1.20+）
	//    对应 Java 收集 ConstraintViolation 列表
	joined := errors.Join(validateAge(-5), errors.New("db connection refused"))
	fmt.Println("errors.Join          =", joined)
	fmt.Println()

	compareWithJava()
}

func compareWithJava() {
	fmt.Println("--- Java vs Go 对比 ---")
	fmt.Println()
	fmt.Println("Java 异常（_checked 逼你处理，运行时异常容易漏_）:")
	fmt.Println(`  try {
      User u = userService.find(id);          // 必须 catch 或 throws
  } catch (NotFoundException e) {              // 按类型匹配
      log.warn("user {} not found: {}", id, e.getField());
  }`)
	fmt.Println()
	fmt.Println("Go 错误是值（编译器逼你至少看一眼 err）:")
	fmt.Println(`  u, err := svc.Find(ctx, id)
  if err != nil {
      if errors.Is(err, storage.ErrNotFound) { // 按哨兵匹配
          // ... 404
      }
      return fmt.Errorf("handle request: %w", err) // 处理不了就包装上抛
  }`)
	fmt.Println()
	fmt.Println("Go惯用法:")
	fmt.Println("  - 错误是普通值：函数签名把 error 放最后一个返回值")
	fmt.Println("  - 要么处理，要么用 %w 包装上抛，绝不 _ = err 吞掉")
	fmt.Println("  - 判断用 errors.Is/As（走错误链），不要 err.Error() 字符串匹配")
	fmt.Printf("  - 包 %%w 不包 %%v：链断了调用方就没法判断根因\n")
	fmt.Println("  - 哨兵错误命名 ErrXxx，放在包级导出，文档写清语义")
}
