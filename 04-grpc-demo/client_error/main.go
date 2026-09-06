// gRPC 错误处理与超时控制 —— 客户端
//
// 演示如何解析 gRPC 错误码，以及 context 超时和取消机制。
//
// Go错误处理核心理念（对比Java）:
//   Java: try { ... } catch (XxxException e) { ... }  // 异常驱动
//   Go:   resp, err := call(); if err != nil { ... }   // 返回值驱动
//
//   优点: 错误处理强制可见，不会被意外吞掉
//   难点: 每个调用都要处理 err，代码量增加（但逻辑更清晰）
//
// 运行: 先启动 go run server_error/main.go，再运行本文件
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "grpc-demo/gen/go/demo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	conn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewDemoServiceClient(conn)

	fmt.Println("=== gRPC错误处理与超时控制 ===")
	fmt.Println()

	demoErrorCodes(client)
	fmt.Println()

	demoTimeout(client)
	fmt.Println()

	demoCancel(client)
}

// demoErrorCodes 演示如何识别和处理不同的 gRPC 错误码。
//
// status.Convert(err) 是解析 gRPC 错误的标准方式
// 对比Java: e.getStatus().getCode() / Status.fromThrowable(e).getCode()
func demoErrorCodes(client pb.DemoServiceClient) {
	fmt.Println("--- 错误码演示 ---")
	fmt.Println()

	testCases := []struct {
		id   int32
		desc string
	}{
		{1, "正常返回"},
		{2, "NotFound错误"},
		{3, "PermissionDenied错误"},
		{4, "InvalidArgument错误"},
		{99, "Internal错误"},
	}

	for _, tc := range testCases {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		resp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: tc.id})
		cancel() // 立即 cancel 而不是 defer，避免循环中 context 积累

		if err != nil {
			// status.Convert 将 error 转换为 *status.Status，可以读取 Code 和 Message
			// 对比Java: StatusRuntimeException e; e.getStatus().getCode()
			st := status.Convert(err)
			fmt.Printf("ID=%d (%s): code=%s, message=%s\n", tc.id, tc.desc, st.Code(), st.Message())
		} else {
			fmt.Printf("ID=%d (%s): user=%+v\n", tc.id, tc.desc, resp.User)
		}
	}

	fmt.Println()
	fmt.Println("常用错误码速查:")
	fmt.Println("  codes.OK                -> HTTP 200 成功")
	fmt.Println("  codes.NotFound          -> HTTP 404 资源不存在")
	fmt.Println("  codes.PermissionDenied  -> HTTP 403 已认证但无权限")
	fmt.Println("  codes.Unauthenticated   -> HTTP 401 未认证")
	fmt.Println("  codes.InvalidArgument   -> HTTP 400 参数错误")
	fmt.Println("  codes.Internal          -> HTTP 500 内部错误")
	fmt.Println("  codes.Unavailable       -> HTTP 503 服务不可用")
	fmt.Println("  codes.DeadlineExceeded  -> HTTP 504 超时")
}

// demoTimeout 演示 context.WithTimeout 超时控制。
//
// Go惯用法: context.WithTimeout 是控制超时的统一方式
//   - 超时会自动传播到所有子调用（RPC、数据库、HTTP请求等）
//   - 对比Java: 没有统一机制，每个框架自己实现
//     gRPC Java: stub.withDeadlineAfter(1, TimeUnit.SECONDS)
//     Spring:    @Timeout(1) 注解
func demoTimeout(client pb.DemoServiceClient) {
	fmt.Println("--- 超时控制演示 ---")
	fmt.Println()

	fmt.Println("请求ID=5会延迟3秒，设置1秒超时:")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	start := time.Now()
	_, err := client.GetUser(ctx, &pb.GetUserRequest{Id: 5})
	duration := time.Since(start)

	if err != nil {
		st := status.Convert(err)
		fmt.Printf("错误: code=%s, message=%s\n", st.Code(), st.Message())
		fmt.Printf("实际耗时: %v（约等于设置的超时时间）\n", duration.Round(time.Millisecond))
	}

	fmt.Println()
	fmt.Println("超时后错误码: codes.DeadlineExceeded")
	fmt.Println("Go惯用法: context 作为第一个参数传递，超时自动传播到所有下游调用")
}

// demoCancel 演示 context.WithCancel 主动取消请求。
//
// 场景: 用户主动点击"取消"按钮，或者某个并发请求先返回后取消其他请求
// 对比Java: Future.cancel(true) 或 CompletableFuture + Thread.interrupt()
func demoCancel(client pb.DemoServiceClient) {
	fmt.Println("--- 取消请求演示 ---")
	fmt.Println()

	// WithCancel 返回一个可手动取消的 context
	ctx, cancel := context.WithCancel(context.Background())

	// 模拟500ms后用户主动取消（例如点击UI上的取消按钮）
	go func() {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("用户主动取消请求...")
		cancel() // 调用 cancel() 立即取消所有使用该 ctx 的操作
	}()

	start := time.Now()
	_, err := client.GetUser(ctx, &pb.GetUserRequest{Id: 5}) // ID=5 会慢3秒
	duration := time.Since(start)

	if err != nil {
		st := status.Convert(err)
		fmt.Printf("错误: code=%s\n", st.Code()) // 期望: codes.Canceled
		fmt.Printf("实际耗时: %v（约等于取消前的等待时间）\n", duration.Round(time.Millisecond))
	}

	fmt.Println()
	fmt.Println("取消后错误码: codes.Canceled")
	fmt.Println("对比Java: Future.cancel() 只能取消未开始的任务，")
	fmt.Println("         Go context.cancel() 可以取消正在执行的 RPC 调用")
}
