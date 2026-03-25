# Gin + GORM Web项目实战计划

## 项目概述

创建一个完整的用户管理REST API项目，帮助Java开发者快速掌握Go Web开发。

**技术栈：**
- Web框架：Gin（对比Spring Boot）
- ORM：GORM（对比MyBatis）
- 数据库：PostgreSQL（使用本地已有数据库服务）
- 配置管理：环境变量/配置文件

**数据库连接信息：**
- Host: localhost
- Port: 5433
- User: root
- Password: 15671040800q
- 新建数据库: go_user_api

---

## 项目结构（对比Java Spring Boot）

```
begin/
├── quickstart/              # 入门教程（已完成）
└── user-api/                # 新项目
    ├── main.go              # 入口文件
    ├── config/
    │   └── config.go        # 配置管理
    ├── model/
    │   └── user.go          # 数据模型（对比Entity）
    ├── repository/
    │   └── user_repository.go   # 数据访问层（对比Mapper/DAO）
    ├── service/
    │   └── user_service.go      # 业务逻辑层（对比Service）
    ├── controller/
    │   └── user_controller.go   # 控制器层（对比Controller）
    ├── router/
    │   └── router.go            # 路由配置（对比@RequestMapping）
    └── go.mod                    # 模块管理（对比pom.xml）
```

---

## API设计

| 方法 | 路径 | 功能 | 对比Spring |
|------|------|------|-----------|
| GET | /api/users | 获取用户列表 | @GetMapping |
| GET | /api/users/:id | 获取单个用户 | @GetMapping("/{id}") |
| POST | /api/users | 创建用户 | @PostMapping |
| PUT | /api/users/:id | 更新用户 | @PutMapping |
| DELETE | /api/users/:id | 删除用户 | @DeleteMapping |

---

## 实现步骤

### 第一步：项目初始化
1. 创建项目目录结构
2. 初始化go mod
3. 安装依赖（gin、gorm、postgres驱动）

### 第二步：数据库准备
1. 创建新数据库 go_user_api
2. 配置数据库连接

### 第三步：数据模型层
1. 定义User结构体
2. 配置GORM自动迁移（建表）

### 第四步：数据访问层
1. 实现CRUD基础操作
2. 使用GORM链式查询（对比MyBatis XML/注解）

### 第五步：业务逻辑层
1. 封装业务逻辑
2. 错误处理

### 第六步：控制器层
1. 定义HTTP处理函数
2. 请求参数绑定
3. 响应JSON

### 第七步：路由配置
1. 配置RESTful路由
2. 分组路由

### 第八步：入口与启动
1. 组装各层依赖
2. 启动服务器

### 第九步：测试验证
1. 启动服务
2. 测试各API接口

---

## Java vs Go Web开发对比

### 依赖管理
```go
// Go: go.mod
module user-api
go 1.26
require (
    github.com/gin-gonic/gin v1.9+
    gorm.io/gorm v1.25+
    gorm.io/driver/postgres v1.5+
)
```

### 控制器
```go
// Go: Gin
func (c *UserController) GetUser(ctx *gin.Context) {
    id := ctx.Param("id")
    user, err := c.service.GetByID(id)
    ctx.JSON(200, user)
}

// Java: Spring Boot
@GetMapping("/{id}")
public User getUser(@PathVariable Long id) {
    return userService.getById(id);
}
```

### 数据访问（GORM vs MyBatis）
```go
// Go: GORM（链式调用，类似MyBatis-Plus）
db.Where("name = ?", name).First(&user)
db.Create(&user)
db.Model(&user).Updates(newUser)

// Java: MyBatis XML
// <select id="selectByName" resultType="User">
//   SELECT * FROM users WHERE name = #{name}
// </select>
```

---

## 预期成果

完成后你将掌握：
1. Go项目标准目录结构
2. Gin框架REST API开发
3. GORM数据库操作（对比MyBatis）
4. PostgreSQL连接与操作
5. 分层架构设计
6. 错误处理最佳实践
