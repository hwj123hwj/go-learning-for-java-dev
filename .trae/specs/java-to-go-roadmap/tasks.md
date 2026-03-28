# Tasks

## 阶段一：基础巩固（已完成）
- [x] Task 1: Go基础语法入门
  - [x] 变量声明、数据类型
  - [x] 函数与错误处理
  - [x] 结构体与方法
  - [x] 接口与多态
  - [x] 并发基础（goroutine + channel）
  - [x] 集合操作（slice/map）

## 阶段二：Web开发入门（已完成）
- [x] Task 2: Gin + GORM Web项目（user-api）
  - [x] 项目结构搭建（分层架构）
  - [x] 数据模型定义
  - [x] CRUD接口实现
  - [x] 前端测试界面
  - [x] 项目文档（含 Java 对比）
  - [x] 代码质量修复：接口解耦、依赖注入、环境变量配置、错误处理

## 阶段三：进阶功能（已完成）
- [x] Task 3: 参数校验
  - [x] 集成 validator 库（gin binding 标签）
  - [x] DTO 定义与校验规则（Age 用 *int 指针）
  - [x] 统一错误响应（utils/response.go）
  - [x] 对比 Spring Validation 文档

- [x] Task 4: JWT 认证
  - [x] 用户注册 / 登录接口
  - [x] JWT Token 生成与验证（utils/jwt.go）
  - [x] JWT secret 从配置/环境变量读取，不硬编码
  - [x] 认证中间件（middleware/auth.go）
  - [x] 密码加密（bcrypt）
  - [x] 对比 Spring Security 文档

- [x] Task 5: 日志中间件
  - [x] 请求日志记录（middleware/logger.go）
  - [x] 日志格式化（方法、路径、IP、状态码、耗时）
  - [x] 对比 Java 日志切面

- [x] Task 6: 配置管理
  - [x] 创建 config.yaml（敏感字段留空）
  - [x] 集成 viper 读取配置
  - [x] 环境变量覆盖敏感配置（DB_PASSWORD、JWT_SECRET）
  - [x] .env.example 模板 + .gitignore 保护
  - [x] 对比 Spring 配置

- [x] Task 7: 单元测试
  - [x] Service 层测试（mock repository 接口）
  - [x] CreateUser 成功/邮箱重复场景
  - [x] GetAllUsers / DeleteUser 场景
  - [x] 测试与业务逻辑解耦（依赖接口而非实现）
  - [ ] Controller 层测试（待扩展）
  - [ ] 测试覆盖率报告

- [x] Task 8: 代码质量全面修复
  - [x] Repository 层依赖注入（不依赖 config.DB 全局变量）
  - [x] Service / Controller 依赖接口而非具体结构体
  - [x] 修复 CreateUser 忽略 FindByEmail 错误的 bug
  - [x] Delete/Update 返回正确 HTTP 状态码（404 vs 500）
  - [x] main.go 使用 config.InitDB()，消除重复初始化

## 阶段四：高级特性（待实现）
- [ ] Task 9: 并发进阶
  - [ ] Context 超时控制
  - [ ] sync.WaitGroup 使用
  - [ ] 并发安全 Map
  - [ ] 对比 Java 并发工具

- [ ] Task 10: Redis 缓存
  - [ ] 集成 go-redis
  - [ ] 缓存用户数据
  - [ ] 缓存过期策略
  - [ ] 对比 Java Redis 客户端

- [ ] Task 11: Docker 部署
  - [ ] 编写 Dockerfile
  - [ ] Docker Compose 配置
  - [ ] 多阶段构建优化
  - [ ] 部署文档

## 阶段五：项目完善
- [ ] Task 12: 完善文档
  - [ ] 更新 README（含本次修复说明）
  - [ ] 添加 API 文档
  - [ ] 添加部署说明

# Task Dependencies
- Task 4 依赖 Task 3（JWT 需要校验支持）
- Task 7 依赖 Task 3-6（测试需要功能完成）
- Task 10 依赖 Task 6（Redis 配置需要配置管理）
- Task 11 依赖 Task 3-10（部署需要功能完整）
