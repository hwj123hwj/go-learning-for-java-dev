// gRPC 客户端 —— 四种通信模式演示
//
// 对比Java gRPC客户端:
//   Java:  ManagedChannel channel = ManagedChannelBuilder.forAddress("localhost", 50051).usePlaintext().build()
//          DemoServiceGrpc.DemoServiceBlockingStub stub = DemoServiceGrpc.newBlockingStub(channel)
//   Go:    conn, _ := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
//          client := pb.NewDemoServiceClient(conn)
//
// 核心差异:
//   - Java 有 BlockingStub / AsyncStub / FutureStub 三种 stub，需要选择
//   - Go 只有一种 client，统一用 context 控制超时和取消
//
// 运行: 先启动 go run server/main.go，再运行 go run client/main.go
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "grpc-demo/gen/go/demo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// grpc.NewClient 创建连接（懒连接，不会立即建立TCP连接）
	// insecure.NewCredentials() 表示不使用TLS，仅用于开发测试
	//
	// 对比Java: ManagedChannelBuilder.forAddress(...).usePlaintext().build()
	//
	// Go惯用法: conn 应全局单例复用，不要每次请求都创建新连接
	//           程序退出时 defer conn.Close() 关闭连接
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close() // 对比Java: channel.shutdown().awaitTermination(5, TimeUnit.SECONDS)

	client := pb.NewDemoServiceClient(conn)

	fmt.Println("=== gRPC四种通信模式演示 ===")
	fmt.Println()

	demoUnary(client)
	fmt.Println()

	demoServerStreaming(client)
	fmt.Println()

	demoClientStreaming(client)
	fmt.Println()

	demoBidirectionalStreaming(client)
}

// demoUnary 演示一元RPC调用。
//
// Go惯用法: context.WithTimeout 控制单次请求超时
//   - context 是 Go 中取消、超时、传递请求级数据的标准机制
//   - 对比Java: 没有统一机制，各框架自己实现（CompletableFuture.orTimeout、@Timeout等）
//   - defer cancel() 确保资源释放，不管请求成功还是失败
func demoUnary(client pb.DemoServiceClient) {
	fmt.Println("--- Unary RPC ---")
	fmt.Println("对比Java: 类似普通方法调用（BlockingStub.getUser()）")
	fmt.Println()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: 1})
	if err != nil {
		log.Printf("GetUser error: %v", err)
		return
	}

	fmt.Printf("响应: user=%+v\n", resp.User)
	fmt.Println()
	fmt.Println("Java调用方式:")
	fmt.Println("  GetUserRequest request = GetUserRequest.newBuilder().setId(1).build();")
	fmt.Println("  GetUserResponse response = blockingStub.getUser(request);")
	fmt.Println()
	fmt.Println("Go调用方式:")
	fmt.Println("  resp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: 1})")
	fmt.Println()
	fmt.Println("差异: Java用Builder模式构建请求，Go用struct literal直接赋值")
}

// demoServerStreaming 演示服务端流式RPC。
//
// Go惯用法: for + stream.Recv()，直到 io.EOF
//   对比Java: Iterator<T> 或 StreamObserver 回调
func demoServerStreaming(client pb.DemoServiceClient) {
	fmt.Println("--- Server Streaming ---")
	fmt.Println("场景: 服务端分批返回数据（分页查询、日志流）")
	fmt.Println()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.ListUsers(ctx, &pb.ListUsersRequest{PageSize: 10})
	if err != nil {
		log.Printf("ListUsers error: %v", err)
		return
	}

	// Go惯用法: for 循环接收，io.EOF 表示服务端流结束
	// 对比Java BlockingStub: Iterator<ListUsersResponse> iter = blockingStub.listUsers(req)
	//                        while (iter.hasNext()) { iter.next(); }
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break // 流正常结束
		}
		if err != nil {
			log.Printf("stream recv error: %v", err)
			break
		}
		fmt.Printf("  收到: user=%+v, page=%d\n", resp.User, resp.Page)
	}

	fmt.Println()
	fmt.Println("Go惯用法: for 循环 + stream.Recv()，直到 io.EOF")
}

// demoClientStreaming 演示客户端流式RPC。
//
// Go惯用法: stream.Send() 发送多条，CloseAndRecv() 关闭并接收响应
func demoClientStreaming(client pb.DemoServiceClient) {
	fmt.Println("--- Client Streaming ---")
	fmt.Println("场景: 客户端批量上传数据")
	fmt.Println()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.UploadUsers(ctx)
	if err != nil {
		log.Printf("UploadUsers error: %v", err)
		return
	}

	users := []*pb.User{
		{Id: 10, Name: "User10", Email: "user10@example.com"},
		{Id: 11, Name: "User11", Email: "user11@example.com"},
		{Id: 12, Name: "User12", Email: "user12@example.com"},
	}

	for _, user := range users {
		if err := stream.Send(&pb.UploadUsersRequest{User: user}); err != nil {
			log.Printf("send error: %v", err)
			break
		}
		fmt.Printf("  发送: user=%+v\n", user)
		time.Sleep(100 * time.Millisecond)
	}

	// CloseAndRecv: 关闭发送端，等待服务端返回汇总响应
	// 对比Java AsyncStub: requestObserver.onCompleted() 后通过 responseObserver 回调接收
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("CloseAndRecv error: %v", err)
		return
	}

	fmt.Printf("服务端响应: count=%d, message=%s\n", resp.Count, resp.Message)
	fmt.Println()
	fmt.Println("Go惯用法: stream.Send() 逐条发送，CloseAndRecv() 关闭并接收响应")
}

// demoBidirectionalStreaming 演示双向流式RPC。
//
// Go惯用法: goroutine 分离收发
//   - 主 goroutine 负责发送
//   - 子 goroutine 负责接收
//   - 用 channel 同步两个 goroutine 的生命周期
func demoBidirectionalStreaming(client pb.DemoServiceClient) {
	fmt.Println("--- Bidirectional Streaming ---")
	fmt.Println("场景: 双向实时通信（聊天、实时协作）")
	fmt.Println()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.Chat(ctx)
	if err != nil {
		log.Printf("Chat error: %v", err)
		return
	}

	messages := []string{"Hello", "How are you?", "Goodbye"}

	// waitc channel 用于等待接收 goroutine 结束
	// 对比Java: CountDownLatch 或 CompletableFuture
	waitc := make(chan struct{})

	// 开启 goroutine 专门接收服务端消息
	// Go惯用法: goroutine + channel 是 Go 并发的核心，对比Java的线程+阻塞队列
	go func() {
		defer close(waitc) // goroutine 结束时关闭 channel，通知主 goroutine
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Printf("recv error: %v", err)
				return
			}
			fmt.Printf("  收到: from=%s, content=%s\n", resp.From, resp.Content)
		}
	}()

	// 主 goroutine 发送消息
	for _, msg := range messages {
		if err := stream.Send(&pb.ChatMessage{
			From:      "Client",
			Content:   msg,
			Timestamp: time.Now().Unix(),
		}); err != nil {
			log.Printf("send error: %v", err)
			break
		}
		fmt.Printf("  发送: content=%s\n", msg)
		time.Sleep(200 * time.Millisecond)
	}

	// 关闭发送端，通知服务端不再发送消息
	// 对比Java AsyncStub: requestObserver.onCompleted()
	stream.CloseSend()

	// 等待接收 goroutine 退出
	// <-waitc 阻塞直到 channel 被关闭
	<-waitc

	fmt.Println()
	fmt.Println("Go惯用法: goroutine 分离收发，channel 同步生命周期")
}
