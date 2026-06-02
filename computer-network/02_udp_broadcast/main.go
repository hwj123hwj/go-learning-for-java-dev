package main

import (
	"fmt"
	"net"
	"time"
)

func startUDPServer() {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("解析地址失败:", err)
		return
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("服务器启动失败:", err)
		return
	}
	defer conn.Close()

	serverAddr := conn.LocalAddr().String()
	fmt.Printf("UDP服务器启动: %s\n", serverAddr)

	go func() {
		buf := make([]byte, 1024)
		for {
			n, remoteAddr, err := conn.ReadFromUDP(buf)
			if err != nil {
				return
			}
			msg := string(buf[:n])
			fmt.Printf("  [服务器] 收到来自%s: %s\n", remoteAddr, msg)

			reply := fmt.Sprintf("UDP回复: %s", msg)
			conn.WriteToUDP([]byte(reply), remoteAddr)
		}
	}()

	time.Sleep(100 * time.Millisecond)
	runUDPClient(serverAddr)
	time.Sleep(100 * time.Millisecond)
}

func runUDPClient(serverAddr string) {
	addr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Println("解析服务器地址失败:", err)
		return
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("连接服务器失败:", err)
		return
	}
	defer conn.Close()
	fmt.Printf("UDP客户端连接到: %s\n", serverAddr)

	messages := []string{"hello udp", "broadcast test", "go network"}
	for _, msg := range messages {
		_, err := conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("发送失败:", err)
			return
		}

		buf := make([]byte, 1024)
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("接收失败:", err)
			return
		}
		fmt.Printf("  [客户端] 发送: %s\n", msg)
		fmt.Printf("  [客户端] 收到: %s\n", string(buf[:n]))
	}
}

func demonstrateUDPHeader() {
	fmt.Println()
	fmt.Println("--- UDP报文首部 ---")
	fmt.Println()
	fmt.Println("  0      7 8     15 16    23 24    31")
	fmt.Println("  +--------+--------+--------+--------+")
	fmt.Println("  |   源端口(16)  |  目的端口(16) |")
	fmt.Println("  +--------+--------+--------+--------+")
	fmt.Println("  |   长度(16)    |   校验和(16)  |")
	fmt.Println("  +--------+--------+--------+--------+")
	fmt.Println("  |            数据部分            |")
	fmt.Println()
	fmt.Println("UDP首部仅8字节，TCP首部20-60字节")
	fmt.Println()
	fmt.Println("--- TCP vs UDP对比 ---")
	fmt.Println()
	fmt.Printf("%-15s %-20s %-20s\n", "特性", "TCP", "UDP")
	fmt.Printf("%-15s %-20s %-20s\n", "连接", "面向连接", "无连接")
	fmt.Printf("%-15s %-20s %-20s\n", "可靠性", "可靠传输", "不可靠")
	fmt.Printf("%-15s %-20s %-20s\n", "传输方式", "字节流", "数据报")
	fmt.Printf("%-15s %-20s %-20s\n", "首部开销", "20-60字节", "8字节")
	fmt.Printf("%-15s %-20s %-20s\n", "拥塞控制", "有", "无")
	fmt.Printf("%-15s %-20s %-20s\n", "流量控制", "有(滑动窗口)", "无")
	fmt.Printf("%-15s %-20s %-20s\n", "通信模式", "一对一", "一对一/一对多/多对多")
	fmt.Printf("%-15s %-20s %-20s\n", "应用场景", "HTTP/FTP/SMTP", "DNS/视频/直播")
}

func demonstrateARQ() {
	fmt.Println()
	fmt.Println("--- ARQ协议 ---")
	fmt.Println()
	fmt.Println("1. 停止-等待协议(SW-ARQ):")
	fmt.Println("   发送一帧 -> 等待ACK -> 发送下一帧")
	fmt.Println("   信道利用率: U = 1/(1+2a), a=传播延迟/发送时间")
	fmt.Println()
	fmt.Println("2. 后退N帧协议(GBN-ARQ):")
	fmt.Println("   发送窗口>1, 接收窗口=1")
	fmt.Println("   出错时从出错帧开始全部重传")
	fmt.Println("   序号空间: >= 发送窗口+1")
	fmt.Println()
	fmt.Println("3. 选择重传协议(SR-ARQ):")
	fmt.Println("   发送窗口>1, 接收窗口>1")
	fmt.Println("   只重传出错的帧")
	fmt.Println("   序号空间: >= 发送窗口+接收窗口")
	fmt.Println("   窗口大小: <= (MAX_SEQ+1)/2")
}

func main() {
	fmt.Println("=== UDP编程与协议原理 ===")
	fmt.Println()

	startUDPServer()
	demonstrateUDPHeader()
	demonstrateARQ()

	fmt.Println()
	fmt.Println("=== 考研要点 ===")
	fmt.Println("1. UDP: 无连接、不可靠、面向报文、首部8字节")
	fmt.Println("2. UDP校验和: 伪首部+首部+数据，可选(IPv4)")
	fmt.Println("3. 停止-等待: 最简单，效率低，U=1/(1+2a)")
	fmt.Println("4. GBN: 累积确认，出错回退N帧重传")
	fmt.Println("5. SR: 逐个确认，只重传出错帧，效率最高")
	fmt.Println("6. 滑动窗口大小限制: GBN<=2^n-1, SR<=2^(n-1)")
	fmt.Println("7. UDP适用: 实时应用(视频/语音)、DNS查询、TFTP")
	fmt.Println("8. 信道利用率: a越大(传播延迟大)，SW效率越低")
}
