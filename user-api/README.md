# 用户管理系统 - Go Web项目实战

> 专为Java开发者设计的Go Web开发入门项目，通过对比学习快速掌握Gin + GORM技术栈。

## 快速开始

```bash
cd user-api
go run main.go
```

访问 http://localhost:8080 即可看到前端界面。

---

## 项目结构

```
user-api/
├── main.go                    # 入口文件（对比Application.java）
├── go.mod                     # 模块管理（对比pom.xml）
├── go.sum                     # 依赖校验（对比pom.xml中的版本锁定）
├── config/
│   ├── config.go              # 配置管理
│   └── database.go            # 数据库连接
├── model/
│   └── user.go                # 数据模型（对比Entity/POJO）
├── repository/
│   └── user_repository.go     # 数据访问层（对比Mapper/DAO）
├── service/
│   └── user_service.go        # 业务逻辑层（对比Service）
├── controller/
│   └── user_controller.go     # 控制器层（对比Controller）
├── router/
│   └── router.go              # 路由配置（对比@RequestMapping）
└── static/
    └── index.html             # 前端页面
```

---

## Java vs Go 对比速查表

### 1. 项目初始化

| Java Spring Boot | Go |
|------------------|-----|
| `spring init` 创建项目 | `go mod init user-api` |
| `pom.xml` 管理依赖 | `go.mod` 管理依赖 |
| `mvn install` 安装依赖 | `go mod tidy` 安装依赖 |

### 2. 入口文件

```go
// Go: main.go
package main

func main() {
    // 初始化数据库
    // 组装依赖
    // 启动服务器
    r.Run(":8080")
}
```

```java
// Java: Application.java
@SpringBootApplication
public class Application {
    public static void main(String[] args) {
        SpringApplication.run(Application.class, args);
    }
}
```

### 3. 实体类

```go
// Go: model/user.go
type User struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Name      string         `json:"name" gorm:"size:100;not null"`
    Email     string         `json:"email" gorm:"uniqueIndex"`
    Age       int            `json:"age"`
    CreatedAt time.Time      `json:"created_at"`
}
```

```java
// Java: User.java
@Entity
@Table(name = "users")
public class User {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    @Column(nullable = false, length = 100)
    private String name;
    
    @Column(unique = true)
    private String email;
    
    private Integer age;
    
    @CreatedDate
    private LocalDateTime createdAt;
}
```

### 4. 数据访问层

```go
// Go: repository/user_repository.go
func (r *UserRepository) FindAll() ([]model.User, error) {
    var users []model.User
    err := config.DB.Find(&users).Error
    return users, err
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
    var user model.User
    err := config.DB.First(&user, id).Error
    return &user, err
}
```

```java
// Java: UserMapper.java 或 UserMapper.xml
@Select("SELECT * FROM users")
List<User> findAll();

@Select("SELECT * FROM users WHERE id = #{id}")
User findById(Long id);
```

### 5. 业务逻辑层

```go
// Go: service/user_service.go
type UserService struct {
    repo *repository.UserRepository
}

func (s *UserService) CreateUser(user *model.User) error {
    existing, _ := s.repo.FindByEmail(user.Email)
    if existing.ID != 0 {
        return errors.New("email already exists")
    }
    return s.repo.Create(user)
}
```

```java
// Java: UserService.java
@Service
public class UserService {
    @Autowired
    private UserMapper userMapper;
    
    public void createUser(User user) {
        User existing = userMapper.findByEmail(user.getEmail());
        if (existing != null) {
            throw new RuntimeException("email already exists");
        }
        userMapper.insert(user);
    }
}
```

### 6. 控制器层

```go
// Go: controller/user_controller.go
func (c *UserController) GetAll(ctx *gin.Context) {
    users, err := c.service.GetAllUsers()
    if err != nil {
        ctx.JSON(500, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(200, users)
}
```

```java
// Java: UserController.java
@RestController
@RequestMapping("/api/users")
public class UserController {
    @Autowired
    private UserService userService;
    
    @GetMapping
    public List<User> getAll() {
        return userService.getAllUsers();
    }
}
```

### 7. 路由配置

```go
// Go: router/router.go
func SetupRouter(userController *controller.UserController) *gin.Engine {
    r := gin.Default()
    
    api := r.Group("/api")
    {
        users := api.Group("/users")
        {
            users.GET("", userController.GetAll)
            users.GET("/:id", userController.GetByID)
            users.POST("", userController.Create)
            users.PUT("/:id", userController.Update)
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
    @GetMapping
    public List<User> getAll() { }
    
    @GetMapping("/{id}")
    public User getById(@PathVariable Long id) { }
    
    @PostMapping
    public User create(@RequestBody User user) { }
    
    @PutMapping("/{id}")
    public User update(@PathVariable Long id, @RequestBody User user) { }
    
    @DeleteMapping("/{id}")
    public void delete(@PathVariable Long id) { }
}
```

---

## GORM vs MyBatis 常用操作对比

| 操作 | GORM | MyBatis |
|------|------|---------|
| 查询全部 | `DB.Find(&users)` | `SELECT * FROM users` |
| 查询单条 | `DB.First(&user, id)` | `SELECT * FROM users WHERE id = #{id}` |
| 条件查询 | `DB.Where("name = ?", name).Find(&users)` | `SELECT * FROM users WHERE name = #{name}` |
| 插入 | `DB.Create(&user)` | `INSERT INTO users ...` |
| 更新 | `DB.Save(&user)` | `UPDATE users SET ...` |
| 删除 | `DB.Delete(&User{}, id)` | `DELETE FROM users WHERE id = #{id}` |

---

## API接口

| 方法 | 路径 | 功能 |
|------|------|------|
| GET | /api/users | 获取用户列表 |
| GET | /api/users/:id | 获取单个用户 |
| POST | /api/users | 创建用户 |
| PUT | /api/users/:id | 更新用户 |
| DELETE | /api/users/:id | 删除用户 |

---

## 核心概念对比

### 依赖注入

| Java Spring | Go |
|-------------|-----|
| `@Autowired` 自动注入 | 手动组装依赖 |
| Spring容器管理 | 显式传递依赖 |

```go
// Go: 手动组装
userRepo := repository.NewUserRepository()
userService := service.NewUserService(userRepo)
userController := controller.NewUserController(userService)
```

### 错误处理

| Java | Go |
|------|-----|
| try-catch | 多返回值 + error |
| 异常向上抛 | 显式处理每个错误 |

```go
// Go: 显式错误处理
user, err := service.GetUserByID(id)
if err != nil {
    ctx.JSON(404, gin.H{"error": "user not found"})
    return
}
```

### 空值处理

| Java | Go |
|------|-----|
| `null` | `nil` |
| NullPointerException | 不会panic，但需要检查 |

---

## 常用命令

```bash
# 运行项目
go run main.go

# 编译项目
go build -o user-api

# 运行编译后的程序
./user-api

# 安装依赖
go mod tidy

# 查看依赖
go list -m all
```

---

## 下一步学习

1. **参数校验** - 使用 `validator` 库
2. **JWT认证** - 使用 `jwt-go`
3. **日志中间件** - 自定义Gin中间件
4. **配置管理** - 使用 `viper` 读取配置文件
5. **单元测试** - 使用 `go test`
6. **热加载** - 使用 `air` 工具

---

## 参考资源

- [Gin官方文档](https://gin-gonic.com/docs/)
- [GORM官方文档](https://gorm.io/docs/)
- [Go官方教程](https://go.dev/tour/)
