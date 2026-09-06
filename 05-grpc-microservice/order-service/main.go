// 订单微服务
//
// 独立运行在 :50062，演示服务间 gRPC 调用。
// 创建订单时，通过 gRPC 调用 UserService 获取用户信息。
//
// 架构: Client -> OrderService(:50062) -> UserService(:50061)
//
// 连接管理最佳实践（对比Java）:
//   Java:  @FeignClient 由 Spring 管理连接池和生命周期，开发者无需关心
//   Go:    grpc.NewClient 在服务启动时创建一次，全局复用
//          程序退出时 defer conn.Close() 关闭连接
//          不要在每次请求时创建新连接（开销大，且不会被自动释放）
//
// 运行: 先启动 go run user-service/main.go，再运行本文件
package main

import (
	"context"
	"log"
	"net"
	"sync"

	orderpb "grpc-microservice/proto/order"
	userpb "grpc-microservice/proto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// orderServiceServer 持有对 UserService 的 gRPC 连接。
//
// 关键设计: userConn 和 userClient 在服务启动时初始化，全程复用。
// 对比Java @Autowired: Spring 自动注入 FeignClient，Go 需要手动管理
type orderServiceServer struct {
	orderpb.UnimplementedOrderServiceServer
	mu         sync.RWMutex
	orders     map[int32]*orderpb.Order
	nextID     int32
	userConn   *grpc.ClientConn          // 保存连接引用，用于程序退出时关闭
	userClient userpb.UserServiceClient  // gRPC 客户端，线程安全，可并发使用
}

func newOrderServiceServer() *orderServiceServer {
	// 在服务启动时建立到 UserService 的连接
	// grpc.NewClient 是懒连接，不会立即建立 TCP 连接，首次 RPC 调用时才真正连接
	// 对比Java: new FeignClient() 或 WebClient.create("http://user-service")
	conn, err := grpc.NewClient("localhost:50061", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to user service: %v", err)
	}

	log.Println("Connected to UserService at localhost:50061")

	return &orderServiceServer{
		orders:     make(map[int32]*orderpb.Order),
		nextID:     1,
		userConn:   conn,
		userClient: userpb.NewUserServiceClient(conn),
	}
}

// CreateOrder 创建订单，并通过 gRPC 调用 UserService 验证用户是否存在。
//
// 注意 context 传递: 将入参 ctx 透传给 userClient.GetUser()
//   - 这样客户端的超时/取消会自动传播到下游服务调用
//   - 对比Java Feign: RequestInterceptor 传递 Header，但超时需要单独配置
//   - Go 的 context 机制天然支持分布式超时传播
func (s *orderServiceServer) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	// 调用 UserService，传递 context 以便超时/取消传播
	userResp, err := s.userClient.GetUser(ctx, &userpb.GetUserRequest{Id: req.UserId})
	if err != nil {
		log.Printf("[OrderService] Failed to get user: %v", err)
		// 将上游错误包装后返回，不要直接暴露内部服务错误
		return nil, status.Errorf(codes.NotFound, "user not found: %d", req.UserId)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 存储核心数据，不冗余存 user_name
	order := &orderpb.Order{
		Id:       s.nextID,
		UserId:   req.UserId,
		Product:  req.Product,
		Quantity: req.Quantity,
	}
	s.orders[s.nextID] = order
	s.nextID++

	// 响应时构建 OrderView（包含 user_name），与存储模型分离
	view := &orderpb.Order{
		Id:       order.Id,
		UserId:   order.UserId,
		UserName: userResp.User.Name, // 从 UserService 响应填充
		Product:  order.Product,
		Quantity: order.Quantity,
	}

	log.Printf("[OrderService] CreateOrder: id=%d, user=%s, product=%s", order.Id, view.UserName, order.Product)
	return &orderpb.CreateOrderResponse{Order: view}, nil
}

func (s *orderServiceServer) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.GetOrderResponse, error) {
	s.mu.RLock()
	order, ok := s.orders[req.Id]
	s.mu.RUnlock()

	if !ok {
		return nil, status.Errorf(codes.NotFound, "order not found: %d", req.Id)
	}

	// 查询订单时也实时调用 UserService 获取最新用户名
	userResp, err := s.userClient.GetUser(ctx, &userpb.GetUserRequest{Id: order.UserId})
	userName := "unknown"
	if err != nil {
		// 非关键错误，降级处理而不是直接返回错误
		log.Printf("[OrderService] Failed to get user for order %d: %v", req.Id, err)
	} else {
		userName = userResp.User.Name
	}

	view := &orderpb.Order{
		Id:       order.Id,
		UserId:   order.UserId,
		UserName: userName,
		Product:  order.Product,
		Quantity: order.Quantity,
	}

	log.Printf("[OrderService] GetOrder: id=%d", req.Id)
	return &orderpb.GetOrderResponse{Order: view}, nil
}

func main() {
	server := newOrderServiceServer()
	// defer 确保程序退出时关闭到 UserService 的连接
	// 对比Java: Spring 容器关闭时自动调用 @PreDestroy
	defer server.userConn.Close()

	lis, err := net.Listen("tcp", ":50062")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	orderpb.RegisterOrderServiceServer(s, server)

	log.Println("OrderService listening on :50062")
	log.Println("依赖: UserService at localhost:50061")
	log.Println()
	log.Println("对比Java Spring Cloud:")
	log.Println("  Java: @FeignClient 自动生成客户端，Eureka/Nacos 服务发现")
	log.Println("  Go:   grpc.NewClient 手动创建，连接全局单例复用")
	log.Println("        生产环境推荐: consul/etcd + grpc resolver 实现服务发现")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
