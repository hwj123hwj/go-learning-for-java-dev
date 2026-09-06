// 微服务演示客户端
//
// 演示完整调用链: Client -> OrderService(:50062) -> UserService(:50061)
//
// 运行顺序:
//   1. go run user-service/main.go
//   2. go run order-service/main.go
//   3. go run client/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	orderpb "grpc-microservice/proto/order"
	userpb "grpc-microservice/proto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 同时与两个服务建立连接
	// Go惯用法: 每个连接全局单例，defer 关闭
	userConn, err := grpc.NewClient("localhost:50061", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to user service: %v", err)
	}
	defer userConn.Close()

	orderConn, err := grpc.NewClient("localhost:50062", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to order service: %v", err)
	}
	defer orderConn.Close()

	userClient := userpb.NewUserServiceClient(userConn)
	orderClient := orderpb.NewOrderServiceClient(orderConn)

	fmt.Println("=== gRPC微服务调用链演示 ===")
	fmt.Println("架构: Client -> OrderService(:50062) -> UserService(:50061)")
	fmt.Println()

	// 所有请求共用一个带超时的 context
	// 注意: 这里的超时会传播到 OrderService 对 UserService 的调用
	// 即: 如果整体超时5秒，那么服务间调用也必须在5秒内完成
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("--- 步骤1: 直接调用UserService创建用户 ---")
	createUserResp, err := userClient.CreateUser(ctx, &userpb.CreateUserRequest{
		Name:  "Charlie",
		Email: "charlie@example.com",
	})
	if err != nil {
		log.Printf("CreateUser error: %v", err)
	} else {
		fmt.Printf("创建用户: id=%d, name=%s\n", createUserResp.User.Id, createUserResp.User.Name)
	}
	fmt.Println()

	fmt.Println("--- 步骤2: 调用OrderService创建订单（内部调用UserService）---")
	fmt.Println("调用链: Client -> OrderService -> UserService（验证用户存在）")
	createOrderResp, err := orderClient.CreateOrder(ctx, &orderpb.CreateOrderRequest{
		UserId:   1,
		Product:  "iPhone 15",
		Quantity: 2,
	})
	if err != nil {
		log.Printf("CreateOrder error: %v", err)
	} else {
		fmt.Printf("创建订单: id=%d, user=%s, product=%s, quantity=%d\n",
			createOrderResp.Order.Id,
			createOrderResp.Order.UserName,
			createOrderResp.Order.Product,
			createOrderResp.Order.Quantity)
	}
	fmt.Println()

	fmt.Println("--- 步骤3: 查询订单 ---")
	getOrderResp, err := orderClient.GetOrder(ctx, &orderpb.GetOrderRequest{Id: 1})
	if err != nil {
		log.Printf("GetOrder error: %v", err)
	} else {
		fmt.Printf("查询订单: id=%d, user=%s, product=%s\n",
			getOrderResp.Order.Id,
			getOrderResp.Order.UserName,
			getOrderResp.Order.Product)
	}
	fmt.Println()

	fmt.Println("--- 步骤4: 测试不存在的用户（错误传播）---")
	_, err = orderClient.CreateOrder(ctx, &orderpb.CreateOrderRequest{
		UserId:   999,
		Product:  "MacBook",
		Quantity: 1,
	})
	if err != nil {
		fmt.Printf("预期错误（用户不存在）: %v\n", err)
	}
	fmt.Println()

	fmt.Println("=== 对比Java Spring Cloud ===")
	fmt.Println("Java:  @FeignClient(name=\"user-service\") 声明式调用")
	fmt.Println("       Ribbon/LoadBalancer 自动负载均衡")
	fmt.Println("       Eureka/Nacos 服务注册发现")
	fmt.Println()
	fmt.Println("Go:    grpc.NewClient() 手动创建连接，显式管理生命周期")
	fmt.Println("       需要负载均衡时: grpc.WithDefaultServiceConfig + round_robin")
	fmt.Println("       需要服务发现时: 集成 consul/etcd + grpc resolver")
	fmt.Println("       优点: 行为透明，无魔法，便于调试")
}
