# 用户管理系统 - 进阶版

> 基于user-api扩展的进阶版本，包含参数校验、JWT认证、日志中间件等企业级功能。

## 与基础版对比

| 功能 | user-api（基础版） | user-api-advanced（进阶版） |
|------|-------------------|---------------------------|
| CRUD接口 | ✅ | ✅ |
| 参数校验 | ❌ | ✅ validator |
| JWT认证 | ❌ | ✅ jwt-go |
| 日志中间件 | ❌ | ✅ 自定义中间件 |
| 配置管理 | 硬编码 | ✅ viper + yaml |
| 单元测试 | ❌ | ✅ go test |

## 快速开始

```bash
cd user-api-advanced
go run main.go
```

访问 http://localhost:8080

---

## 项目结构

```
user-api-advanced/
├── main.go                    # 入口文件
├── config.yaml                # 配置文件
├── config/
│   ├── config.go              # 配置管理（viper）
│   └── database.go            # 数据库连接
├── middleware/                # 中间件
│   ├── auth.go                # JWT认证中间件
│   └── logger.go              # 日志中间件
├── model/
│   └── user.go                # 数据模型
├── repository/
│   └── user_repository.go     # 数据访问层
├── service/
│   ├── user_service.go        # 业务逻辑层
│   └── user_service_test.go   # 单元测试
├── controller/
│   ├── user_controller.go     # 用户控制器
│   └── auth_controller.go     # 认证控制器
├── dto/
│   └── user_dto.go            # 数据传输对象（含校验规则）
├── router/
│   └── router.go              # 路由配置
├── utils/                     # 工具函数
│   ├── response.go            # 统一响应
│   ├── jwt.go                 # JWT工具
│   └── password.go            # 密码加密
└── static/
    └── index.html             # 前端页面
```

---

## Java vs Go 对比详解

### 1. 参数校验

**Java Spring Validation:**
```java
public class CreateUserRequest {
    @NotBlank(message = "姓名不能为空")
    @Size(min = 2, max = 50, message = "姓名长度2-50")
    private String name;
    
    @NotBlank(message = "邮箱不能为空")
    @Email(message = "邮箱格式不正确")
    private String email;
}
```

**Go validator:**
```go
type CreateUserRequest struct {
    Name  string `json:"name" binding:"required,min=2,max=50"`
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"required,gte=0,lte=150"`
}
```

### 2. JWT认证

**Java Spring Security:**
```java
@Component
public class JwtTokenProvider {
    public String generateToken(UserDetails userDetails) {
        return Jwts.builder()
            .setSubject(userDetails.getUsername())
            .setExpiration(new Date(System.currentTimeMillis() + 86400000))
            .signWith(SignatureAlgorithm.HS256, secretKey)
            .compact();
    }
}
```

**Go jwt-go:**
```go
func GenerateToken(userID uint) (string, error) {
    claims := Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}
```

### 3. 中间件

**Java Spring AOP:**
```java
@Aspect
@Component
public class LoggingAspect {
    @Around("execution(* com.example.controller.*.*(..))")
    public Object logAround(ProceedingJoinPoint joinPoint) throws Throwable {
        long start = System.currentTimeMillis();
        Object result = joinPoint.proceed();
        long elapsed = System.currentTimeMillis() - start;
        log.info("{} took {}ms", joinPoint.getSignature(), elapsed);
        return result;
    }
}
```

**Go Gin中间件:**
```go
func LoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        startTime := time.Now()
        c.Next()
        latency := time.Since(startTime)
        log.Printf("[%s] %s %d %v", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), latency)
    }
}
```

### 4. 配置管理

**Java application.yml:**
```yaml
server:
  port: 8080
spring:
  datasource:
    url: jdbc:postgresql://localhost:5433/go_user_api
    username: root
    password: 15671040800q
```

**Go config.yaml:**
```yaml
server:
  port: "8080"
database:
  host: "localhost"
  port: "5433"
  user: "root"
  password: "15671040800q"
  name: "go_user_api"
jwt:
  secret: "your-secret-key"
  expire: "24h"
```

### 5. 单元测试

**Java JUnit + Mockito:**
```java
@ExtendWith(MockitoExtension.class)
class UserServiceTest {
    @Mock
    private UserRepository userRepository;
    
    @InjectMocks
    private UserService userService;
    
    @Test
    void shouldCreateUser() {
        when(userRepository.findByEmail(any())).thenReturn(Optional.empty());
        userService.createUser(user);
        verify(userRepository).save(user);
    }
}
```

**Go testing + testify:**
```go
func TestUserService_CreateUser_Success(t *testing.T) {
    mockRepo := new(MockUserRepository)
    userService := NewUserService(mockRepo)
    
    user := &model.User{Name: "Alice", Email: "alice@test.com"}
    mockRepo.On("FindByEmail", "alice@test.com").Return(&model.User{}, nil)
    mockRepo.On("Create", user).Return(nil)
    
    err := userService.CreateUser(user)
    assert.NoError(t, err)
}
```

---

## API接口

| 方法 | 路径 | 功能 | 认证 |
|------|------|------|------|
| POST | /api/auth/register | 注册 | ❌ |
| POST | /api/auth/login | 登录 | ❌ |
| GET | /api/users | 用户列表 | ✅ |
| GET | /api/users/:id | 用户详情 | ✅ |
| POST | /api/users | 创建用户 | ✅ |
| PUT | /api/users/:id | 更新用户 | ✅ |
| DELETE | /api/users/:id | 删除用户 | ✅ |

---

## 运行测试

```bash
cd user-api-advanced
go test ./... -v
```

---

## 核心概念对比总结

| 概念 | Java | Go |
|------|------|-----|
| 依赖注入 | Spring自动注入 | 手动组装 |
| 错误处理 | try-catch | 多返回值 + error |
| 空值 | null | nil |
| 接口 | 显式implements | 隐式实现（duck typing） |
| 并发 | 线程 + ExecutorService | goroutine + channel |
| 测试 | JUnit + Mockito | testing + testify |
