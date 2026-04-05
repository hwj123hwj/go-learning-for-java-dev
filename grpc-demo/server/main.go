// gRPC 服务端 —— 四种通信模式演示
//
// 对比Java: 类似 @GrpcService 注解的服务实现类
// 核心区别:
//   - Java 用继承 XXXGrpc.XXXImplBase 并 @Override 方法
//   - Go  用嵌入 pb.UnimplementedXXXServer（前向兼容），然后实现对应方法
//
// 运行: go run server/main.go
// 客户端: go run client/main.go
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	pb "grpc-demo/gen/go/demo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// server 实现 DemoServiceServer 接口。
//
// 对比Java: 类似 @Service + implements DemoServiceGrpc.DemoServiceImplBase
// 嵌入 UnimplementedDemoServiceServer 的作用:
//   - 提供所有 RPC 方法的默认实现（返回 codes.Unimplemented）
//   - 新增 proto 方法时不会导致编译报错，符合开闭原则
//   - Java 里对应的是 ImplBase 抽象类，未实现的方法默认返回 UNIMPLEMENTED
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

// GetUser 演示 Unary RPC —— 最简单的请求-响应模式。
//
// 对比Java gRPC:
//   Java:  public void getUser(GetUserRequest req, StreamObserver<GetUserResponse> observer)
//          observer.onNext(response); observer.onCompleted();
//   Go:    func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error)
//
// Go 的函数签名更直观: 直接返回 (响应, error)，不需要手动回调 observer。
func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	log.Printf("[Unary] GetUser called, id=%d", req.Id)

	user, ok := s.users[req.Id]
	if !ok {
		// Go惯用法: 用 status.Errorf 返回带错误码的错误，不要用 fmt.Errorf
		// 原因: 客户端通过 status.Code(err) 识别错误类型
		//       fmt.Errorf 返回普通 error，gRPC 会把它包装成 codes.Unknown
		// 对比Java: throw new StatusRuntimeException(Status.NOT_FOUND.withDescription("..."))
		return nil, status.Errorf(codes.NotFound, "user not found: %d", req.Id)
	}

	return &pb.GetUserResponse{User: user}, nil
}

// ListUsers 演示 Server Streaming —— 服务端推送多条消息。
//
// 对比Java gRPC:
//   Java:  调用 responseObserver.onNext(item) 发送每条，最后 onCompleted()
//   Go:    调用 stream.Send(item) 发送，return nil 表示流正常结束
//
// 场景: 分页查询、日志推送、实时事件流
func (s *server) ListUsers(req *pb.ListUsersRequest, stream pb.DemoService_ListUsersServer) error {
	log.Printf("[ServerStream] ListUsers called, page_size=%d", req.PageSize)

	page := int32(1)
	for id, user := range s.users {
		if err := stream.Send(&pb.ListUsersResponse{
			User: user,
			Page: page,
		}); err != nil {
			// Send 出错通常意味着客户端断开，直接返回
			return err
		}
		log.Printf("[ServerStream] sent user id=%d, page=%d", id, page)
		time.Sleep(200 * time.Millisecond)
		page++
	}

	// return nil 表示流正常结束，gRPC 框架自动向客户端发送 EOF 信号
	return nil
}

// UploadUsers 演示 Client Streaming —— 客户端批量上传。
//
// 对比Java gRPC:
//   Java:  返回 StreamObserver，框架回调 onNext/onCompleted/onError，逻辑分散
//   Go:    主动循环调用 stream.Recv() 拉取数据，逻辑集中，更线性
//
// Bug陷阱: io.EOF 和真实错误必须分开处理，不能统一 break 或 return err
func (s *server) UploadUsers(stream pb.DemoService_UploadUsersServer) error {
	log.Println("[ClientStream] UploadUsers called")

	count := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// io.EOF = 客户端正常关闭流，不是错误
			// 对比Java: 等同于 StreamObserver.onCompleted() 被调用
			break
		}
		if err != nil {
			// 真实错误（网络中断等）必须返回，不能忽略
			// Go没有 try-catch，错误必须显式向上传递
			return err
		}

		count++
		s.users[req.User.Id] = req.User
		log.Printf("[ClientStream] received user id=%d, name=%s", req.User.Id, req.User.Name)
	}

	// SendAndClose: 发送唯一响应并关闭流（Client Streaming 专用）
	return stream.SendAndClose(&pb.UploadUsersResponse{
		Count:   count,
		Message: fmt.Sprintf("成功上传 %d 个用户", count),
	})
}

// Chat 演示 Bidirectional Streaming —— 全双工通信。
//
// 对比Java gRPC:
//   Java:  回调模式，逻辑分散在 onNext/onCompleted/onError 里
//   Go:    for 循环同步读取，逻辑线性，更容易阅读和调试
//
// Bug陷阱: 收到 io.EOF 必须 return nil，不能 return err
//   原因: io.EOF 表示客户端正常关闭，不是错误
//         return err 会让客户端收到 codes.Unknown 错误码
func (s *server) Chat(stream pb.DemoService_ChatServer) error {
	log.Println("[BidirectionalStream] Chat started")

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// 客户端调用了 CloseSend()，表示不再发送
			// 必须 return nil，不能 return err
			log.Println("[BidirectionalStream] Chat ended normally")
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("[BidirectionalStream] received: from=%s, content=%s", msg.From, msg.Content)

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

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// grpc.NewServer() 对比Java: 类似 ServerBuilder.forPort(50051).build()
	// 需要拦截器时: grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor1, interceptor2))
	s := grpc.NewServer()

	// 注册服务实现，对比Java: 类似 serverBuilder.addService(new DemoServiceImpl())
	pb.RegisterDemoServiceServer(s, newServer())

	log.Println("gRPC server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
