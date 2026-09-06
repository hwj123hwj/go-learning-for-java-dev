# gRPC微服务实战（Java开发者版）

> 通过对比Java gRPC，掌握Go的gRPC开发最佳实践。

## 模块学习目标

- **模块一：gRPC基础** — 学完后能独立定义proto文件、生成Go代码，并实现四种通信模式的服务端与客户端
- **模块二：gRPC进阶** — 学完后能为gRPC服务添加日志、认证、错误处理和超时控制，具备生产可用的基础能力
- **模块三：微服务实战** — 学完后能拆分两个独立服务并实现服务间gRPC调用，理解Go微服务与Java Spring Cloud的架构差异

## 目录结构

```
grpc-demo/                    # 模块一：gRPC基础
├── proto/                    # Protobuf定义
├── server/                   # 基础服务端
├── client/                   # 基础客户端
├── server_interceptor/       # 模块二：拦截器服务端
├── client_interceptor/       # 模块二：拦截器客户端
├── server_error/             # 模块二：错误处理服务端
├── client_error/             # 模块二：错误处理客户端
└── gen/go/demo/              # 生成的Go代码

grpc-microservice/            # 模块三：微服务实战
├── proto/user/               # 用户服务proto
├── proto/order/              # 订单服务proto
├── user-service/             # 用户服务（端口50061）
├── order-service/            # 订单服务（端口50062，调用用户服务）
└── client/                   # 测试客户端
```

## 快速开始

### 模块一：四种通信模式

```bash
# 启动服务端
cd grpc-demo
go run server/main.go

# 新终端运行客户端
go run client/main.go
```

### 模块二：拦截器与错误处理

```bash
# 启动拦截器服务端
go run server_interceptor/main.go

# 新终端运行客户端
go run client_interceptor/main.go

# 错误处理演示（另一端口50052）
go run server_error/main.go &
go run client_error/main.go
```

### 模块三：微服务调用

```bash
# 终端1：启动用户服务
cd grpc-microservice
go run user-service/main.go

# 终端2：启动订单服务
go run order-service/main.go

# 终端3：运行客户端
go run client/main.go
```

## Java vs Go 对比

### 拦截器对比

| 方面 | Java | Go |
|------|------|-----|
| 概念 | ServerInterceptor | grpc.UnaryServerInterceptor |
| 注册 | intercept()方法 | grpc.ChainUnaryInterceptor() |
| 执行顺序 | 按注册顺序 | 洋葱模型，先入后出 |

```java
// Java Spring AOP
@Around("execution(* com.example.service.*.*(..))")
public Object logMethod(ProceedingJoinPoint pjp) {
    log.info("Before: {}", pjp.getSignature());
    Object result = pjp.proceed();
    log.info("After: {}", pjp.getSignature());
    return result;
}
```

```go
// Go拦截器
func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    log.Printf("Before: %s", info.FullMethod)
    resp, err := handler(ctx, req)
    log.Printf("After: %s", info.FullMethod)
    return resp, err
}
```

### 错误处理对比

| 方面 | Java | Go |
|------|------|-----|
| 错误类型 | StatusRuntimeException | status.Status |
| 创建错误 | Status.NOT_FOUND.asRuntimeException() | status.Error(codes.NotFound, "msg") |
| 获取错误码 | e.getStatus().getCode() | status.Convert(err).Code() |

### 超时控制对比

```java
// Java: 使用Deadline
ManagedChannel channel = ManagedChannelBuilder.forAddress("localhost", 50051)
    .build();
UserGrpc.UserBlockingStub stub = UserGrpc.newBlockingStub(channel)
    .withDeadline(Deadline.after(1, TimeUnit.SECONDS));
```

```go
// Go: 使用context
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
defer cancel()
resp, err := client.GetUser(ctx, req)
```

### 微服务调用对比

| 方面 | Java Spring Cloud | Go gRPC |
|------|-------------------|---------|
| 服务发现 | Eureka/Nacos | 需自己集成consul/etcd |
| 负载均衡 | Ribbon | grpc内置round_robin |
| 声明式客户端 | @FeignClient | 手动创建grpc.Client |
| 连接管理 | Spring管理 | 需自己管理生命周期 |

```java
// Java: @FeignClient自动生成客户端
@FeignClient(name = "user-service")
public interface UserClient {
    @GetMapping("/users/{id}")
    User getUser(@PathVariable Long id);
}
```

```go
// Go: 手动创建和管理连接
conn, err := grpc.NewClient("localhost:50061", grpc.WithTransportCredentials(insecure.NewCredentials()))
if err != nil {
    log.Fatal(err)
}
defer conn.Close()  // 程序退出时关闭

client := userpb.NewUserServiceClient(conn)
resp, err := client.GetUser(ctx, &userpb.GetUserRequest{Id: 1})
```

## 四种通信模式

| 模式 | 说明 | 场景 |
|------|------|------|
| Unary | 请求-响应 | 普通API调用 |
| Server Streaming | 服务端流式返回 | 分页查询、日志流 |
| Client Streaming | 客户端流式上传 | 文件上传、批量写入 |
| Bidirectional | 双向实时通信 | 聊天、实时协作 |

## Go惯用法

1. **使用buf工具** — 统一管理proto文件和代码生成
2. **grpc.NewClient** — 使用新API（grpc.Dial已废弃）
3. **连接单例复用** — 全局初始化，defer关闭
4. **context传递** — 作为第一个参数
5. **错误处理** — 函数返回error，调用方处理
6. **拦截器链** — ChainUnaryInterceptor按洋葱模型执行
7. **panic恢复** — RecoveryInterceptor捕获panic并返回Internal错误
