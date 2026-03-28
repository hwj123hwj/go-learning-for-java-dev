# Java 转 Go 完整学习路线规划

## Why

帮助有 Java 后端经验的开发者快速掌握 Go 语言技术栈，通过对比学习和实战项目，系统性地完成技术栈迁移。

## What Changes

### 已完成内容
- ✅ Go 基础语法入门（7 个核心脚本）
- ✅ Gin + GORM Web 项目实战（user-api，基础 CRUD）
- ✅ 前端测试界面
- ✅ Java vs Go 对比文档
- ✅ 参数校验（validator，对比 Spring Validation）
- ✅ JWT 认证（golang-jwt，对比 Spring Security）
- ✅ 日志中间件（对比 Java AOP 日志切面）
- ✅ 配置管理（viper + 环境变量，对比 Spring application.yml）
- ✅ 单元测试（testify mock，对比 JUnit + Mockito）
- ✅ 代码质量全面修复（两个项目）：
  - 密码/Secret 禁止硬编码，改用环境变量
  - Repository 依赖注入，消除全局变量依赖
  - Service / Controller 依赖接口，支持 mock 测试
  - 修复 CreateUser 忽略 FindByEmail 错误的 bug
  - HTTP 状态码规范（404/500 分类正确）
  - Age 字段使用 *int 指针，区分"未填"和"填 0"

### 待补充内容
- 并发进阶（context、sync 包，对比 Java 并发工具）
- Redis 缓存（go-redis，对比 Java Redis 客户端）
- Docker 部署（多阶段构建）
- Controller 层测试 + 覆盖率报告

## Impact

- 学习路径：quickstart → user-api → user-api-advanced → 高级特性
- 目标：完成 Java 到 Go 的完整技术栈迁移

---

## ADDED Requirements

### Requirement: 参数校验
系统应提供请求参数校验功能，对比 Java 的 Spring Validation。

#### Scenario: 创建用户时校验
- **WHEN** 用户提交空姓名或无效邮箱
- **THEN** 返回 400 错误和具体校验信息

### Requirement: JWT 认证
系统应提供 JWT 认证功能，对比 Java 的 Spring Security。

#### Scenario: 访问受保护接口
- **WHEN** 用户未登录访问需要认证的接口
- **THEN** 返回 401 未授权错误

#### Scenario: 登录获取 Token
- **WHEN** 用户使用正确凭证登录
- **THEN** 返回有效 JWT Token

### Requirement: 日志中间件
系统应提供请求日志记录功能，对比 Java 的日志切面。

#### Scenario: 记录请求日志
- **WHEN** 每个 HTTP 请求到达
- **THEN** 记录请求方法、路径、耗时、状态码

### Requirement: 配置管理
系统应支持从配置文件读取配置，对比 Spring 的 application.yml。
敏感信息（密码、JWT secret）必须通过环境变量注入，禁止明文写入配置文件。

#### Scenario: 加载配置
- **WHEN** 应用启动
- **THEN** 从 config.yaml 读取基础配置，环境变量覆盖敏感字段

### Requirement: 单元测试
系统应提供单元测试，对比 Java 的 JUnit + Mockito。

#### Scenario: 测试 Service 层
- **WHEN** 运行 go test ./service/...
- **THEN** 4 个用例全部通过，业务逻辑通过 mock 验证

### Requirement: 并发进阶
系统应展示 Go 并发高级特性，对比 Java 并发工具。

#### Scenario: 使用 context 超时控制
- **WHEN** 数据库查询超时
- **THEN** 自动取消请求并返回超时错误

### Requirement: Redis 缓存
系统应集成 Redis 缓存，对比 Java 的 Redis 客户端。

#### Scenario: 缓存用户数据
- **WHEN** 查询用户
- **THEN** 优先从缓存获取，缓存未命中再查数据库

### Requirement: Docker 部署
系统应支持 Docker 容器化部署。

#### Scenario: 构建镜像
- **WHEN** 执行 docker build
- **THEN** 生成可运行的 Docker 镜像
