// gRPC 拦截器服务端演示
//
// 演示三种拦截器组合: Recovery -> Logging -> Auth -> Handler
//
// 对比Java Spring AOP:
//   Java:  @Around -> @Before -> 方法 -> @After/AfterReturning，基于注解，运行时动态织入
//   Go:    ChainUnaryInterceptor 按注册顺序形成洋葱模型，编译期确定，无反射开销
//
// 运行: go run server_interceptor/main.go
// 客户端: go run client_interceptor/main.go
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"runtime"
	"time"

	pb "grpc-demo/gen/go/demo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedDemoServiceServer
	users map[int32]*pb.User
}

func newServer() *server {
	return &server{
		users: map[int32]*pb.User{
			1: {Id: 1, Name: "Alice", Email: "alice@example.com"},
			2: {Id: 2, Name: "Bob", Email: "bob@example.com"},
			3: {Id: 3, Name: "Charlie", Email: "charlie@example.com"},
		},
	}
}

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	log.Printf("[Handler] GetUser called, id=%d", req.Id)

	if req.Id == 999 {
		// 模拟 panic 场景，由 recoveryInterceptor 捕获
		panic("模拟panic场景！")
	}

	user, ok := s.users[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "user not found: %d", req.Id)
	}

	return &pb.GetUserResponse{User: user}, nil
}

func (s *server) ListUsers(req *pb.ListUsersRequest, stream pb.DemoService_ListUsersServer) error {
	for id, user := range s.users {
		if err := stream.Send(&pb.ListUsersResponse{User: user, Page: int32(id)}); err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

func (s *server) UploadUsers(stream pb.DemoService_UploadUsersServer) error {
	count := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// io.EOF = 客户端正常关闭流，不是错误
			break
		}
		if err != nil {
			// 真实错误必须返回，不能忽略
			return err
		}
		count++
		s.users[req.User.Id] = req.User
	}
	return stream.SendAndClose(&pb.UploadUsersResponse{
		Count:   count,
		Message: fmt.Sprintf("成功上传 %d 个用户", count),
	})
}

func (s *server) Chat(stream pb.DemoService_ChatServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// 客户端正常关闭，return nil 而不是 return err
			return nil
		}
		if err != nil {
			return err
		}
		reply := &pb.ChatMessage{
			From:      "Server",
			Content:   fmt.Sprintf("Echo: %s", msg.Content),
			Timestamp: time.Now().Unix(),
		}
		if err := stream.Send(reply); err != nil {
			return err
		}
	}
}

// ========== 拦截器实现 ==========
// 拦截器签名: func(ctx, req, info, handler) (resp, error)
// 调用 handler(ctx, req) 就是调用下一层（可以是下一个拦截器或最终 handler）
// 对比Java AOP: handler 相当于 ProceedingJoinPoint.proceed()

// loggingInterceptor 记录每次请求的方法名和耗时。
//
// 对比Java: 类似 @Around 切面，在 proceed() 前后记录时间
// 注意: 调用 handler 前后分别记录时间，这是 Go 拦截器的标准写法
func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	log.Printf("[Logging] >>> 请求开始: method=%s", info.FullMethod)

	resp, err := handler(ctx, req) // 调用下一层

	log.Printf("[Logging] <<< 请求结束: method=%s, duration=%v, err=%v",
		info.FullMethod, time.Since(start), err)
	return resp, err
}

// authInterceptor 从 metadata 中读取 token 并验证。
//
// gRPC metadata 对比HTTP Header:
//   HTTP:  request.getHeader("Authorization")
//   gRPC:  metadata.FromIncomingContext(ctx) -> md.Get("authorization")
//
// 对比Java Spring Security: 类似 OncePerRequestFilter.doFilterInternal()
func authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "缺少metadata")
	}

	tokens := md.Get("authorization")
	if len(tokens) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "缺少authorization token")
	}

	if tokens[0] != "valid-token" {
		return nil, status.Errorf(codes.Unauthenticated, "无效的token: %s", tokens[0])
	}

	log.Printf("[Auth] 认证通过")
	return handler(ctx, req)
}

// recoveryInterceptor 捕获 handler 中的 panic，防止服务崩溃。
//
// 对比Java: 类似 @ControllerAdvice + @ExceptionHandler(Exception.class)
// Go惯用法: 用 defer + recover() 组合捕获 panic
//   - defer 保证函数退出时执行（无论是否 panic）
//   - recover() 只在 defer 中有效，捕获 panic 值
//   - 捕获后将 panic 转为 gRPC 错误码返回，服务不会崩溃
//
// 注意: 这里用了具名返回值 (resp interface{}, err error)
//       这样 defer 中才能修改 err 的值并返回给调用方
func recoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			// 打印堆栈信息，方便排查问题
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			log.Printf("[Recovery] panic recovered: %v\n%s", r, string(buf[:n]))
			// 将 panic 转为 Internal 错误码返回客户端
			err = status.Errorf(codes.Internal, "内部错误: %v", r)
		}
	}()
	return handler(ctx, req)
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// grpc.ChainUnaryInterceptor 将多个拦截器组合成洋葱模型
	// 执行顺序（进入）: recovery -> logging -> auth -> handler
	// 执行顺序（返回）: handler -> auth -> logging -> recovery
	//
	// 对比Java @Order 注解: Go 按注册顺序，不需要额外配置
	//
	// 为什么 recovery 放最外层?
	//   因为它需要能捕获所有层（包括 logging/auth）的 panic
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recoveryInterceptor,
			loggingInterceptor,
			authInterceptor,
		),
	)
	pb.RegisterDemoServiceServer(s, newServer())

	log.Println("gRPC server (with interceptors) listening on :50051")
	log.Println("拦截器顺序: Recovery -> Logging -> Auth -> Handler")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
