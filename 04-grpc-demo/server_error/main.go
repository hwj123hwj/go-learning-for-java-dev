// gRPC 错误处理与超时控制 —— 服务端
//
// 演示 status 包的标准错误码用法。
//
// 对比Java异常体系:
//   Java:  throw new XxxException() -> @ExceptionHandler -> HTTP状态码
//   Go:    return status.Errorf(codes.Xxx, "msg") -> 客户端 status.Code(err)
//
// 核心区别: Go 没有异常，错误通过返回值传递，强迫调用方显式处理
//
// 运行: go run server_error/main.go
// 客户端: go run client_error/main.go
package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "grpc-demo/gen/go/demo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type errorServer struct {
	pb.UnimplementedDemoServiceServer
}

// GetUser 演示各种标准错误码的使用场景。
//
// gRPC 错误码 vs HTTP 状态码对照:
//   codes.OK                (0)  -> HTTP 200
//   codes.NotFound          (5)  -> HTTP 404
//   codes.PermissionDenied  (7)  -> HTTP 403
//   codes.InvalidArgument   (3)  -> HTTP 400
//   codes.Internal          (13) -> HTTP 500
//   codes.Unavailable       (14) -> HTTP 503
//   codes.DeadlineExceeded  (4)  -> HTTP 504
//   codes.Unauthenticated   (16) -> HTTP 401
//
// 对比Java: 类似 ResponseEntity<>(body, HttpStatus.NOT_FOUND)
//           但 gRPC 错误码是传输层标准，跨语言通用
func (s *errorServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	switch req.Id {
	case 1:
		// 正常返回
		return &pb.GetUserResponse{User: &pb.User{Id: 1, Name: "Alice", Email: "alice@example.com"}}, nil
	case 2:
		// codes.NotFound: 资源不存在
		// 对比Java: throw new ResourceNotFoundException("用户不存在")
		return nil, status.Error(codes.NotFound, "用户不存在")
	case 3:
		// codes.PermissionDenied: 无权访问（已认证但无权限）
		// 注意: 未认证用 codes.Unauthenticated，已认证但无权限用 codes.PermissionDenied
		return nil, status.Error(codes.PermissionDenied, "无权访问")
	case 4:
		// codes.InvalidArgument: 请求参数非法
		// 对比Java: throw new IllegalArgumentException() 或 @Valid 校验失败
		return nil, status.Error(codes.InvalidArgument, "参数错误: id无效")
	case 5:
		// 模拟慢请求，用于演示客户端超时控制
		time.Sleep(3 * time.Second)
		return &pb.GetUserResponse{User: &pb.User{Id: 5, Name: "Slow", Email: "slow@example.com"}}, nil
	default:
		// codes.Internal: 服务内部错误，不应暴露细节给客户端
		// 对比Java: throw new InternalServerException() -> HTTP 500
		return nil, status.Error(codes.Internal, "未知错误")
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterDemoServiceServer(s, &errorServer{})

	log.Println("Error handling server listening on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
