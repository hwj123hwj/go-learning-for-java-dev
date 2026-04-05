# Tasks

## 模块一：gRPC基础

- [x] Task 1: 创建grpc-demo项目结构
  - [x] 创建目录和go.mod
  - [x] 安装protoc和Go插件
  - [x] 创建README.md说明文档

- [x] Task 2: Protobuf基础与buf工具
  - [x] 安装buf CLI，创建buf.yaml和buf.gen.yaml配置文件
  - [x] 创建proto文件定义message（对比Java POJO：Go用struct literal，无需Builder模式）
  - [x] 创建proto文件定义service
  - [x] 使用`buf generate`生成Go代码（对比Java Maven插件生成方式）
  - [x] 在README中用注释对比说明：Java生成Builder风格 vs Go生成struct直接赋值风格

- [x] Task 3: Unary RPC
  - [x] 实现服务端handler
  - [x] 实现客户端调用
  - [x] 演示请求-响应流程
  - [x] 对比Java gRPC调用方式

- [x] Task 4: Server Streaming
  - [x] 定义流式proto
  - [x] 实现服务端流式返回
  - [x] 实现客户端流式接收
  - [x] 场景：分页查询演示

- [x] Task 5: Client Streaming
  - [x] 定义客户端流式proto
  - [x] 实现客户端流式发送
  - [x] 实现服务端批量处理
  - [x] 场景：批量上传演示

- [x] Task 6: Bidirectional Streaming
  - [x] 定义双向流proto
  - [x] 实现双向实时通信
  - [x] 场景：简单聊天演示

## 模块二：gRPC进阶

- [x] Task 7: 拦截器
  - [x] 实现日志拦截器（记录方法名、耗时）
  - [x] 实现认证拦截器（从metadata中读取token，对比Java Filter从Header读取）
  - [x] 实现Recovery拦截器（panic恢复，对比Java @ControllerAdvice）
  - [x] 使用`grpc.ChainUnaryInterceptor`将三个拦截器组合成chain，演示执行顺序
  - [x] 对比Java Spring AOP执行顺序（Around/Before/After）

- [x] Task 8: 错误处理
  - [x] 使用status包返回标准错误码
  - [x] 自定义错误详情
  - [x] 客户端错误处理

- [x] Task 9: 超时与取消
  - [x] 使用context.WithTimeout
  - [x] 演示超时取消传播
  - [x] 演示客户端取消请求

## 模块三：微服务实战

- [x] Task 10: 创建grpc-microservice项目结构
  - [x] 创建user-service目录
  - [x] 创建order-service目录
  - [x] 创建共享proto目录

- [x] Task 11: 用户服务
  - [x] 定义用户服务proto
  - [x] 实现用户CRUD
  - [x] 启动用户服务

- [x] Task 12: 订单服务
  - [x] 定义订单服务proto
  - [x] 实现订单CRUD
  - [x] 使用`grpc.NewClient`创建用户服务连接（全局单例，程序启动时初始化）
  - [x] 调用用户服务获取用户信息，演示服务间通信
  - [x] 演示`defer conn.Close()`正确关闭连接（对比Java框架自动管理连接生命周期）

- [x] Task 13: 整合演示
  - [x] 同时启动两个服务
  - [x] 演示完整调用链
  - [x] 对比Java微服务架构

# Task Dependencies
- Task 3-6 依赖 Task 2（需要先生成代码）
- Task 7-9 依赖 Task 3（需要基础RPC才能演示拦截器）
- Task 11-12 依赖 Task 10（需要项目结构）
- Task 12 依赖 Task 11（订单服务调用用户服务）
