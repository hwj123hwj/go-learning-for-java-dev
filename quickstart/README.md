# Go语言快速入门（Java开发者版）

> 专为Java后端开发者设计的Go语言入门教程，通过对比学习快速掌握Go核心特性。

## 快速开始

```bash
cd quickstart
go run 01_hello.go
```

## 教程目录

| 序号 | 文件 | 主题 | 核心知识点 |
|------|------|------|-----------|
| 01 | [01_hello.go](01_hello.go) | Go基础 | 程序结构、包导入 |
| 02 | [02_variable.go](02_variable.go) | 变量声明 | `:=`短声明、常量、基本类型 |
| 03 | [03_function.go](03_function.go) | 函数与错误 | 多返回值、error处理、可变参数 |
| 04 | [04_struct.go](04_struct.go) | 结构体 | struct定义、方法、嵌入（组合） |
| 05 | [05_interface.go](05_interface.go) | 接口 | duck typing、空接口 |
| 06 | [06_concurrency.go](06_concurrency.go) | 并发 | goroutine、channel、select |
| 07 | [07_collection.go](07_collection.go) | 集合 | slice、map、for-range |

## Java vs Go 核心差异速查

### 1. 变量声明

```go
// Go: 短声明（推荐）
name := "Alice"
age := 25

// Java
String name = "Alice";
int age = 25;
```

### 2. 错误处理

```go
// Go: 多返回值 + 显式错误
result, err := divide(10, 2)
if err != nil {
    // 处理错误
}

// Java: try-catch
try {
    int result = divide(10, 2);
} catch (Exception e) {
    // 处理异常
}
```

### 3. 面向对象

```go
// Go: struct + 方法
type User struct {
    Name string
    Age  int
}

func (u User) SayHello() {
    fmt.Println("Hello, " + u.Name)
}

// Java: class
public class User {
    private String name;
    private int age;
    
    public void sayHello() {
        System.out.println("Hello, " + name);
    }
}
```

### 4. 接口（duck typing）

```go
// Go: 隐式实现，不需要显式声明implements
type Animal interface {
    Speak() string
}

type Dog struct{}

func (d Dog) Speak() string { return "汪汪" }
// Dog自动实现了Animal接口！

// Java: 显式实现
public class Dog implements Animal {
    @Override
    public String speak() { return "汪汪"; }
}
```

### 5. 并发

```go
// Go: goroutine + channel
go func() {
    ch <- "hello"
}()
msg := <-ch

// Java: ExecutorService
executor.submit(() -> {
    // ...
});
```

### 6. 集合

```go
// Go: slice（动态数组）
nums := []int{1, 2, 3}
nums = append(nums, 4)

// Go: map
m := map[string]int{"a": 1, "b": 2}

// Java: List/Map
List<Integer> nums = new ArrayList<>(Arrays.asList(1, 2, 3));
Map<String, Integer> m = new HashMap<>();
```

## 学习建议

1. **按顺序运行每个脚本**，观察输出
2. **修改代码**，加深理解
3. **对比Java**，理解设计差异
4. **实践项目**，巩固知识

## 下一步学习方向

- Web框架：Gin（类似Spring Boot）
- ORM：GORM（类似MyBatis/Hibernate）
- 微服务：gRPC、go-micro
- 工具链：go mod、go test
