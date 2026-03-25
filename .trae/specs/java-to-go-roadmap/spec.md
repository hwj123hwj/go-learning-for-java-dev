# Java转Go完整学习路线规划

## Why

帮助有Java后端经验的开发者快速掌握Go语言技术栈，通过对比学习和实战项目，系统性地完成技术栈迁移。

## What Changes

### 已完成内容
- ✅ Go基础语法入门（7个核心脚本）
- ✅ Gin + GORM Web项目实战（用户管理CRUD）
- ✅ 前端测试界面
- ✅ Java vs Go对比文档

### 待补充内容
- 参数校验（对比Spring Validation）
- JWT认证（对比Spring Security）
- 日志中间件
- 配置管理（对比Spring配置）
- 单元测试（对比JUnit）
- 并发进阶（context、sync包）
- 缓存集成（Redis）
- Docker部署

## Impact

- 学习路径：quickstart → user-api → 进阶功能
- 目标：完成Java到Go的完整技术栈迁移

---

## ADDED Requirements

### Requirement: 参数校验
系统应提供请求参数校验功能，对比Java的Spring Validation。

#### Scenario: 创建用户时校验
- **WHEN** 用户提交空姓名或无效邮箱
- **THEN** 返回400错误和具体校验信息

### Requirement: JWT认证
系统应提供JWT认证功能，对比Java的Spring Security。

#### Scenario: 访问受保护接口
- **WHEN** 用户未登录访问需要认证的接口
- **THEN** 返回401未授权错误

#### Scenario: 登录获取Token
- **WHEN** 用户使用正确凭证登录
- **THEN** 返回有效JWT Token

### Requirement: 日志中间件
系统应提供请求日志记录功能，对比Java的日志切面。

#### Scenario: 记录请求日志
- **WHEN** 每个HTTP请求到达
- **THEN** 记录请求方法、路径、耗时、状态码

### Requirement: 配置管理
系统应支持从配置文件读取配置，对比Spring的application.yml。

#### Scenario: 加载配置
- **WHEN** 应用启动
- **THEN** 从config.yaml读取数据库、服务器等配置

### Requirement: 单元测试
系统应提供单元测试，对比Java的JUnit。

#### Scenario: 测试Service层
- **WHEN** 运行测试
- **THEN** 验证业务逻辑正确性

### Requirement: 并发进阶
系统应展示Go并发高级特性，对比Java并发工具。

#### Scenario: 使用context超时控制
- **WHEN** 数据库查询超时
- **THEN** 自动取消请求并返回超时错误

### Requirement: Redis缓存
系统应集成Redis缓存，对比Java的Redis客户端。

#### Scenario: 缓存用户数据
- **WHEN** 查询用户
- **THEN** 优先从缓存获取，缓存未命中再查数据库

### Requirement: Docker部署
系统应支持Docker容器化部署。

#### Scenario: 构建镜像
- **WHEN** 执行docker build
- **THEN** 生成可运行的Docker镜像
