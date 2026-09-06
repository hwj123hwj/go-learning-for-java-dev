// gRPC 拦截器客户端演示
//
// 演示如何在客户端发送 metadata（类似HTTP Header）
//
// 对比Java:
//   Java:  stub.withCallCredentials(credentials) 或 ClientInterceptor 注入 Header
//   Go:    metadata.AppendToOutgoingContext(ctx, "key", "value")
//
// 运行: 先启动 go run server_interceptor/main.go，再运行本文件
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "grpc-demo/gen/go/demo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewDemoServiceClient(conn)

	fmt.Println("=== gRPC拦截器演示 ===")
	fmt.Println()

	demoWithAuth(client)
	fmt.Println()

	demoWithoutAuth(client)
	fmt.Println()

	demoPanicRecovery(client)
}

// demoWithAuth 演示携带有效 token 的请求。
//
// metadata.AppendToOutgoingContext 向 context 追加出站 metadata
// 对比Java: Metadata headers = new Metadata(); headers.put(AUTH_HEADER_KEY, "valid-token")
//           stub.withInterceptors(MetadataUtils.newAttachHeadersInterceptor(headers))
func demoWithAuth(client pb.DemoServiceClient) {
	fmt.Println("--- 带有效token的请求 ---")

	// metadata.AppendToOutgoingContext 是 Go 向 gRPC 请求附加 Header 的标准方式
	// key 会被自动转成小写（gRPC 规范），对应服务端 md.Get("authorization")
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "valid-token")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: 1})
	if err != nil {
		log.Printf("GetUser error: %v", err)
		return
	}

	fmt.Printf("响应: user=%+v\n", resp.User)
	fmt.Println()
	fmt.Println("服务端拦截器执行顺序:")
	fmt.Println("  1. RecoveryInterceptor（最外层，捕获所有panic）")
	fmt.Println("  2. LoggingInterceptor（记录开始时间）")
	fmt.Println("  3. AuthInterceptor（验证token）")
	fmt.Println("  4. Handler（执行业务逻辑）")
	fmt.Println("  3. AuthInterceptor 返回")
	fmt.Println("  2. LoggingInterceptor 记录结束时间")
	fmt.Println("  1. RecoveryInterceptor 返回")
}

func demoWithoutAuth(client pb.DemoServiceClient) {
	fmt.Println("--- 不带token的请求（预期被拦截）---")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.GetUser(ctx, &pb.GetUserRequest{Id: 1})
	if err != nil {
		fmt.Printf("预期错误: %v\n", err)
	}

	fmt.Println()
	fmt.Println("对比Java Spring Security:")
	fmt.Println("  Java: SecurityFilterChain 拦截请求，返回 401 Unauthorized")
	fmt.Println("  Go:   authInterceptor 返回 codes.Unauthenticated")
}

func demoPanicRecovery(client pb.DemoServiceClient) {
	fmt.Println("--- Panic恢复演示 (id=999 触发服务端panic) ---")

	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "valid-token")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := client.GetUser(ctx, &pb.GetUserRequest{Id: 999})
	if err != nil {
		fmt.Printf("捕获的错误: %v\n", err)
	}

	fmt.Println()
	fmt.Println("对比Java @ControllerAdvice:")
	fmt.Println("  Java: @ExceptionHandler(Exception.class) 捕获全局异常，返回统一错误体")
	fmt.Println("  Go:   recoveryInterceptor 用 defer+recover() 捕获panic，返回 codes.Internal")
	fmt.Println("  效果相同: 服务不崩溃，客户端收到友好的错误码")
}
