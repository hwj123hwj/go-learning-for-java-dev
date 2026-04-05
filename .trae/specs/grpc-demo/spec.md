# gRPC微服务实战 Spec

## Why

gRPC是Go微服务开发的核心通信协议，相比REST API有更高性能和强类型约束。通过对比Java的gRPC/Dubbo，掌握Go的gRPC开发最佳实践。

## 模块学习目标

- **模块一：gRPC基础** — 学完后能独立定义proto文件、生成Go代码，并实现四种通信模式的服务端与客户端
- **模块二：gRPC进阶** — 学完后能为gRPC服务添加日志、认证、错误处理和超时控制，具备生产可用的基础能力
- **模块三：微服务实战** — 学完后能拆分两个独立服务并实现服务间gRPC调用，理解Go微服务与Java Spring Cloud的架构差异

## What Changes

### 新增模块
- `grpc-demo/` - gRPC基础与进阶
  - Protobuf定义与代码生成
  - Unary RPC（一元调用）
  - Server Streaming（服务端流）
  - Client Streaming（客户端流）
  - Bidirectional Streaming（双向流）
  - 拦截器（Interceptor）
  - 错误处理与超时控制

- `grpc-microservice/` - 微服务实战
  - 用户服务（user-service）
  - 订单服务（order-service）
  - 服务间调用演示

## Impact

- 学习路径：quickstart → user-api → concurrency-demo → backend-patterns → grpc-demo → grpc-microservice
- 目标：掌握gRPC开发，理解微服务通信

---

## ADDED Requirements

### Requirement: Protobuf基础
系统应提供Protobuf定义与代码生成的学习材料。

#### Scenario: proto文件定义
- **WHEN** 定义服务接口
- **THEN** 编写.proto文件，定义message和service
- **对比Java**：proto文件相同，生成的代码风格不同

#### Scenario: 代码生成
- **WHEN** 需要生成Go代码
- **THEN** 使用protoc + protoc-gen-go + protoc-gen-go-grpc
- **Go惯用法**：使用buf工具统一管理proto文件和代码生成（buf.yaml + buf.gen.yaml），替代手写protoc命令

### Requirement: gRPC四种通信模式
系统应演示gRPC的四种通信模式。

#### Scenario: Unary RPC
- **WHEN** 普通请求-响应模式
- **THEN** 客户端发送一个请求，服务端返回一个响应
- **对比Java**：类似普通方法调用

#### Scenario: Server Streaming
- **WHEN** 服务端需要返回大量数据
- **THEN** 服务端分批发送，客户端逐个接收
- **场景**：分页查询、日志流、事件推送

#### Scenario: Client Streaming
- **WHEN** 客户端需要上传大量数据
- **THEN** 客户端分批发送，服务端统一处理
- **场景**：文件上传、批量写入

#### Scenario: Bidirectional Streaming
- **WHEN** 双向实时通信
- **THEN** 客户端和服务端同时发送流
- **场景**：聊天、实时协作、游戏

### Requirement: 拦截器（Interceptor）
系统应演示gRPC拦截器的使用。

#### Scenario: 一元拦截器
- **WHEN** 需要在请求前后处理
- **THEN** 使用UnaryInterceptor实现日志、认证、监控
- **对比Java**：类似Spring AOP / Filter

#### Scenario: 流拦截器
- **WHEN** 需要拦截流式请求
- **THEN** 使用StreamInterceptor

### Requirement: 错误处理与超时
系统应演示gRPC的错误处理机制。

#### Scenario: 错误码
- **WHEN** 服务端返回错误
- **THEN** 使用status包返回标准错误码
- **对比Java**：类似自定义异常

#### Scenario: 超时控制
- **WHEN** 需要控制请求超时
- **THEN** 使用context.WithTimeout
- **Go惯用法**：context作为第一个参数传递

### Requirement: 微服务实战
系统应提供微服务间调用的实战案例。

#### Scenario: 服务定义
- **WHEN** 定义微服务
- **THEN** 每个服务独立运行，通过gRPC通信
- **对比Java**：类似Spring Cloud微服务

#### Scenario: 服务间调用
- **WHEN** 订单服务需要调用用户服务
- **THEN** 通过gRPC客户端调用
- **Go惯用法**：使用`grpc.NewClient`创建连接（`grpc.Dial`已在v1.56+废弃），连接应全局单例复用，程序退出时`defer conn.Close()`
- **注意**：Java框架自动管理连接生命周期，Go需要手动负责
