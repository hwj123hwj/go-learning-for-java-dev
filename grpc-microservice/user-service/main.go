// 用户微服务
//
// 独立运行在 :50061，提供用户的 CRUD 操作。
//
// 对比Java Spring Boot:
//   Java:  @SpringBootApplication + @GrpcService，配置文件指定端口
//   Go:    main() 手动创建 listener + server，没有框架自动扫描
//
// 对比Java Spring Cloud 服务注册:
//   Java:  application.yml 配置 eureka/nacos，服务自动注册
//   Go:    本示例直接硬编码地址；生产中通常用 consul/etcd + grpc resolver
//
// 运行: go run user-service/main.go
package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "grpc-microservice/proto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// userServiceServer 是用户服务的实现。
//
// sync.RWMutex 保护并发读写:
//   - RLock/RUnlock 用于读操作（允许多个 goroutine 并发读）
//   - Lock/Unlock   用于写操作（互斥）
//
// 对比Java: 类似 ReadWriteLock，但 Go 的更简洁
// 对比Java synchronized: synchronized 是互斥锁，不区分读写
type userServiceServer struct {
	pb.UnimplementedUserServiceServer
	mu     sync.RWMutex    // 保护 users map 的并发访问
	users  map[int32]*pb.User
	nextID int32
}

func newUserServiceServer() *userServiceServer {
	return &userServiceServer{
		users: map[int32]*pb.User{
			1: {Id: 1, Name: "Alice", Email: "alice@example.com"},
			2: {Id: 2, Name: "Bob", Email: "bob@example.com"},
		},
		nextID: 3,
	}
}

func (s *userServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// RLock: 读锁，允许并发读，阻塞写
	// 对比Java: readWriteLock.readLock().lock()
	s.mu.RLock()
	defer s.mu.RUnlock() // Go惯用法: defer 确保锁一定会释放，不会忘记

	user, ok := s.users[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "user not found: %d", req.Id)
	}

	log.Printf("[UserService] GetUser: id=%d, name=%s", req.Id, user.Name)
	return &pb.GetUserResponse{User: user}, nil
}

func (s *userServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// Lock: 写锁，互斥
	// 对比Java: readWriteLock.writeLock().lock()
	s.mu.Lock()
	defer s.mu.Unlock()

	user := &pb.User{
		Id:    s.nextID,
		Name:  req.Name,
		Email: req.Email,
	}
	s.users[s.nextID] = user
	s.nextID++

	log.Printf("[UserService] CreateUser: id=%d, name=%s", user.Id, user.Name)
	return &pb.CreateUserResponse{User: user}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50061")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, newUserServiceServer())

	log.Println("UserService listening on :50061")
	log.Println("对比Java: 类似 Spring Boot 应用，独立端口运行")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
