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
- [x] Task 2: Gin + GORM Web项目
  - [x] 项目结构搭建
  - [x] 数据模型定义
  - [x] CRUD接口实现
  - [x] 前端测试界面
  - [x] 项目文档

## 阶段三：进阶功能（已完成）
- [x] Task 3: 参数校验
  - [x] 集成validator库
  - [x] 定义校验规则
  - [x] 统一错误响应
  - [x] 对比Spring Validation文档

- [x] Task 4: JWT认证
  - [x] 用户登录接口
  - [x] JWT Token生成与验证
  - [x] 认证中间件
  - [x] 密码加密（bcrypt）
  - [x] 对比Spring Security文档

- [x] Task 5: 日志中间件
  - [x] 请求日志记录
  - [x] 日志格式化
  - [x] 对比Java日志切面

- [x] Task 6: 配置管理
  - [x] 创建config.yaml
  - [x] 集成viper读取配置
  - [x] 环境变量支持
  - [x] 对比Spring配置

- [ ] Task 7: 单元测试
  - [ ] Service层测试
  - [ ] Controller层测试
  - [ ] 测试覆盖率
  - [ ] 对比JUnit

## 阶段四：高级特性（待实现）
- [ ] Task 8: 并发进阶
  - [ ] Context超时控制
  - [ ] sync.WaitGroup使用
  - [ ] 并发安全Map
  - [ ] 对比Java并发工具

- [ ] Task 9: Redis缓存
  - [ ] 集成go-redis
  - [ ] 缓存用户数据
  - [ ] 缓存过期策略
  - [ ] 对比Java Redis客户端

- [ ] Task 10: Docker部署
  - [ ] 编写Dockerfile
  - [ ] Docker Compose配置
  - [ ] 多阶段构建优化
  - [ ] 部署文档

## 阶段五：项目完善
- [ ] Task 11: 完善文档
  - [ ] 更新README
  - [ ] 添加API文档
  - [ ] 添加部署说明

# Task Dependencies
- Task 4 依赖 Task 3（JWT需要校验支持）
- Task 7 依赖 Task 3-6（测试需要功能完成）
- Task 9 依赖 Task 6（Redis配置需要配置管理）
- Task 10 依赖 Task 3-9（部署需要功能完整）
