# 用户管理系统 - Go Web项目实战

> 专为Java开发者设计的Go Web开发入门项目，通过对比学习快速掌握 Gin + GORM 技术栈。

## 快速开始

### 1. 配置环境变量

```bash
cp .env.example .env
```

编辑 `.env` 文件，填写真实的数据库连接信息：

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=go_user_api
SERVER_PORT=8080
```

> **注意**：`.env` 文件包含敏感信息，已被 `.gitignore` 忽略，**永远不要提交到git**。

### 2. 加载环境变量并启动

```bash
# 方式一：使用 export 逐个设置（临时生效）
export DB_PASSWORD=your_password && go run main.go

# 方式二：使用 source 加载 .env（推荐）
export $(cat .env | grep -v '#' | xargs) && go run main.go
```

访问 http://localhost:8080 即可看到前端界面。

---

## 项目结构

```
user-api/
├── main.go                    # 入口文件：组装依赖、启动服务
├── go.mod                     # 模块管理（对比 pom.xml）
├── go.sum                     # 依赖校验
├── .env.example               # 环境变量模板（可提交）
├── .env                       # 本地环境变量（禁止提交）
├── .gitignore
├── config/
│   ├── config.go              # 从环境变量读取配置
│   └── database.go            # 数据库连接与迁移
├── model/
│   └── user.go                # 数据模型（对比 Entity）
├── repository/
│   └── user_repository.go     # 数据访问层接口 + 实现（对比 DAO/Mapper）
├── service/
│   └── user_service.go        # 业务逻辑层接口 + 实现（对比 Service）
├── controller/
│   └── user_controller.go     # 控制器层（对比 Controller）
├── router/
│   └── router.go              # 路由配置（对比 @RequestMapping）
└── static/
    └── index.html             # 前端页面
```

---

## Java vs Go 核心差异速查

### 1. 项目初始化

| Java Spring Boot | Go |
|---|---|
| `spring init` 创建项目 | `go mod init user-api` |
| `pom.xml` 管理依赖 | `go.mod` 管理依赖 |
| `mvn install` 安装依赖 | `go mod tidy` |

### 2. 配置管理：环境变量替代硬编码

**Java** 用 `application.properties` 或 `@Value`：
```java
// application.properties
spring.datasource.password=your_password

// Java 读取
@Value("${spring.datasource.password}")
private String dbPassword;
```

**Go** 用 `os.Getenv` 读取环境变量（禁止硬编码密码）：
```go
// config/config.go
func LoadConfig() *Config {
    return &Config{
        DBHost:     getEnv("DB_HOST", "localhost"), // 有默认值
        DBPassword: os.Getenv("DB_PASSWORD"),        // 敏感信息，无默认值
    }
}

func getEnv(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}
```

### 3. 依赖注入：手动组装 vs Spring 容器

**Java** 由 Spring 容器自动注入：
```java
@Service
public class UserService {
    @Autowired
    private UserRepository userRepository; // Spring 自动注入
}
```

**Go** 在 `main.go` 中手动显式组装依赖链：
```go
// main.go：依赖注入全貌一目了然
userRepo    := repository.NewUserRepository(config.DB) // 注入 DB
userService := service.NewUserService(userRepo)        // 注入 repo
userCtrl    := controller.NewUserController(userService) // 注入 service
```

### 4. 接口：隐式实现 vs 显式 implements

这是 Go 与 Java **最大的设计差异之一**。

**Java** 接口需要显式声明实现：
```java
public class UserServiceImpl implements UserService { ... }
```

**Go** 接口是隐式的——只要方法签名匹配，就自动满足接口，无需声明：
```go
// 定义接口（通常在调用方一侧）
type UserRepository interface {
    FindAll() ([]model.User, error)
    FindByID(id uint) (*model.User, error)
    // ...
}

// 实现（无需写 "implements UserRepository"）
type userRepository struct { db *gorm.DB }

func (r *userRepository) FindAll() ([]model.User, error) { ... }
func (r *userRepository) FindByID(id uint) (*model.User, error) { ... }
// 方法签名全部匹配 → 自动满足 UserRepository 接口
```

**为什么要用接口？**
- Service 依赖 `repository.UserRepository` 接口，而非 `*userRepository` 具体类型
- 测试时可以注入 mock，不需要真实数据库
- 遵循依赖倒置原则（DIP）

### 5. 实体模型

```go
// Go: model/user.go
type User struct {
    ID    uint   `json:"id"    gorm:"primaryKey"`
    Name  string `json:"name"  gorm:"size:100;not null"`
    Email string `json:"email" gorm:"size:100;uniqueIndex;not null"`
    Age   *int   `json:"age"`  // 指针类型：区分"未填写"(nil) 和"填写了0"
}
```

```java
// Java: User.java
@Entity
public class User {
    @Id @GeneratedValue
    private Long id;
    private String name;
    private String email;
    private Integer age; // 包装类型同样可以为 null
}
```

> **为什么 Age 用 `*int` 而不是 `int`？**
> Go 的 `int` 零值是 `0`，无法区分"用户填了年龄0"和"用户没填年龄"。
> 用指针 `*int` 时，未填写为 `nil`，填写了0则为 `&0`，语义更准确。

### 6. 数据访问层

```go
// Go: repository/user_repository.go
// 结构体小写（unexported），外部只能通过接口访问
type userRepository struct {
    db *gorm.DB // 通过构造函数注入，不依赖全局变量
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) FindAll() ([]model.User, error) {
    var users []model.User
    err := r.db.Find(&users).Error
    return users, err
}
```

```java
// Java: UserMapper.java (MyBatis)
@Select("SELECT * FROM users")
List<User> findAll();
```

### 7. 错误处理：多返回值 vs try-catch

**Java** 用异常：
```java
public User findById(Long id) {
    return userMapper.findById(id); // 找不到可能返回 null 或抛异常
}
```

**Go** 用多返回值，**必须显式处理每个错误**：
```go
// 错误处理是 Go 的核心理念，不能用 _ 随意忽略
user, err := s.repo.FindByID(id)
if err != nil {
    return err
}
```

**注意**：GORM 找不到记录时返回 `gorm.ErrRecordNotFound`，需要区分处理：
```go
_, err := s.repo.FindByEmail(user.Email)
if err == nil {
    return errors.New("email already exists") // 找到了 → 重复
}
if !errors.Is(err, gorm.ErrRecordNotFound) {
    return err // 真正的数据库错误
}
// err == gorm.ErrRecordNotFound → 未找到，可以创建
return s.repo.Create(user)
```

### 8. HTTP 状态码规范

| 场景 | 正确状态码 | 错误示范 |
|---|---|---|
| 创建成功 | `201 Created` | `200 OK` |
| 资源不存在 | `404 Not Found` | `400 Bad Request` |
| 服务器错误 | `500 Internal Server Error` | `400 Bad Request` |
| 参数错误 | `400 Bad Request` | `500` |

### 9. 路由配置

```go
// Go: router/router.go
func SetupRouter(userController *controller.UserController) *gin.Engine {
    r := gin.Default()

    api := r.Group("/api")
    {
        users := api.Group("/users")
        {
            users.GET("",      userController.GetAll)
            users.GET("/:id",  userController.GetByID)
            users.POST("",     userController.Create)
            users.PUT("/:id",  userController.Update)
            users.DELETE("/:id", userController.Delete)
        }
    }
    return r
}
```

```java
// Java: 使用注解方式
@RestController
@RequestMapping("/api/users")
public class UserController {
    @GetMapping          public List<User> getAll() { }
    @GetMapping("/{id}") public User getById(@PathVariable Long id) { }
    @PostMapping         public User create(@RequestBody User user) { }
    @PutMapping("/{id}") public User update(@PathVariable Long id, @RequestBody User user) { }
    @DeleteMapping("/{id}") public void delete(@PathVariable Long id) { }
}
```

---

## GORM vs MyBatis 常用操作对比

| 操作 | GORM | MyBatis |
|---|---|---|
| 查询全部 | `db.Find(&users)` | `SELECT * FROM users` |
| 按ID查询 | `db.First(&user, id)` | `SELECT * FROM users WHERE id = #{id}` |
| 条件查询 | `db.Where("name = ?", name).Find(&users)` | `SELECT * FROM users WHERE name = #{name}` |
| 插入 | `db.Create(&user)` | `INSERT INTO users ...` |
| 更新全部字段 | `db.Save(&user)` | `UPDATE users SET ...` |
| 软删除 | `db.Delete(&User{}, id)` | `UPDATE users SET deleted_at = now()` |

---

## API 接口

| 方法 | 路径 | 功能 | 成功状态码 |
|---|---|---|---|
| GET | /api/users | 获取用户列表 | 200 |
| GET | /api/users/:id | 获取单个用户 | 200 |
| POST | /api/users | 创建用户 | 201 |
| PUT | /api/users/:id | 更新用户 | 200 |
| DELETE | /api/users/:id | 删除用户 | 200 |

---

## 常用命令

```bash
# 安装依赖
go mod tidy

# 运行项目（先设置环境变量）
export $(cat .env | grep -v '#' | xargs) && go run main.go

# 编译
go build -o user-api

# 运行编译后的程序
export $(cat .env | grep -v '#' | xargs) && ./user-api

# 运行测试
go test ./...

# 查看依赖列表
go list -m all
```

---

## 下一步学习

1. **参数校验** - 使用 gin 内置的 `validator` 标签
2. **单元测试** - 用接口 mock Repository 层，`go test ./...`
3. **JWT 认证** - 使用 `github.com/golang-jwt/jwt`
4. **热加载开发** - 使用 `air` 工具（类似 Spring DevTools）
5. **配置管理进阶** - 使用 `github.com/spf13/viper` 支持多种配置源
6. **结构化日志** - 使用 `go.uber.org/zap`

---

## 参考资源

- [Gin 官方文档](https://gin-gonic.com/docs/)
- [GORM 官方文档](https://gorm.io/docs/)
- [Go 官方教程](https://go.dev/tour/)
- [Effective Go](https://go.dev/doc/effective_go)
