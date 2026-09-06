package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

func startTCPServer() {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("服务器启动失败:", err)
		return
	}
	addr := listener.Addr().String()
	fmt.Printf("TCP Echo服务器启动: %s\n", addr)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go handleConnection(conn)
		}
	}()

	time.Sleep(100 * time.Millisecond)
	runTCPClient(addr)
	time.Sleep(100 * time.Millisecond)
	listener.Close()
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	remoteAddr := conn.RemoteAddr().String()
	fmt.Printf("  [服务器] 客户端连接: %s\n", remoteAddr)

	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Printf("  [服务器] 读取错误: %v\n", err)
			}
			break
		}
		fmt.Printf("  [服务器] 收到: %s", line)
		upper := toUpper(line)
		conn.Write([]byte("ECHO: " + upper))
	}
	fmt.Printf("  [服务器] 客户端断开: %s\n", remoteAddr)
}

func toUpper(s string) string {
	result := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			c -= 32
		}
		result = append(result, c)
	}
	return string(result)
}

func runTCPClient(addr string) {
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		fmt.Println("连接服务器失败:", err)
		return
	}
	defer conn.Close()
	fmt.Printf("TCP客户端连接到: %s\n", addr)

	messages := []string{"hello tcp\n", "computer network\n", "go language\n"}
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
		fmt.Printf("  [客户端] 发送: %s", msg)
		fmt.Printf("  [客户端] 收到: %s", string(buf[:n]))
	}
}

func demonstrateTCPStates() {
	fmt.Println()
	fmt.Println("--- TCP三次握手与四次挥手 ---")
	fmt.Println()
	fmt.Println("三次握手(建立连接):")
	fmt.Println("  1. Client -> Server: SYN=1, seq=x")
	fmt.Println("  2. Server -> Client: SYN=1, ACK=1, seq=y, ack=x+1")
	fmt.Println("  3. Client -> Server: ACK=1, seq=x+1, ack=y+1")
	fmt.Println()
	fmt.Println("四次挥手(断开连接):")
	fmt.Println("  1. Client -> Server: FIN=1, seq=u")
	fmt.Println("  2. Server -> Client: ACK=1, seq=v, ack=u+1")
	fmt.Println("     (服务器进入CLOSE_WAIT，客户端进入FIN_WAIT_2)")
	fmt.Println("  3. Server -> Client: FIN=1, ACK=1, seq=w, ack=u+1")
	fmt.Println("  4. Client -> Server: ACK=1, seq=u+1, ack=w+1")
	fmt.Println("     (客户端进入TIME_WAIT，等待2MSL后关闭)")
	fmt.Println()
	fmt.Println("为什么是三次握手?")
	fmt.Println("  防止已失效的连接请求报文突然传到服务器产生错误")
	fmt.Println()
	fmt.Println("为什么是四次挥手?")
	fmt.Println("  TCP全双工，每个方向需单独关闭")
	fmt.Println("  服务器收到FIN后可能还有数据要发，不能立即关闭")
	fmt.Println()
	fmt.Println("TIME_WAIT(2MSL)的意义:")
	fmt.Println("  1. 确保最后一个ACK能到达服务器")
	fmt.Println("  2. 等待本连接所有报文消失，避免新连接收到旧报文")
}

func demonstrateSlidingWindow() {
	fmt.Println()
	fmt.Println("--- TCP滑动窗口 ---")
	fmt.Println()
	fmt.Println("发送窗口:")
	fmt.Println("  已确认     |  已发送未确认  |  可发送未发送  |  不可发送")
	fmt.Println()
	fmt.Println("接收窗口:")
	fmt.Println("  已确认     |  可接收        |  不可接收")
	fmt.Println()
	fmt.Println("流量控制: 接收方通过窗口大小通知发送方")
	fmt.Println("  - 窗口为0: 接收方缓冲区满，发送方停止发送")
	fmt.Println("  - 糊涂窗口综合征: 避免发送小报文(Nagle算法)")
	fmt.Println()
	fmt.Println("拥塞控制:")
	fmt.Println("  1. 慢开始: cwnd从1开始，每RTT翻倍(指数增长)")
	fmt.Println("  2. 拥塞避免: cwnd每RTT加1(线性增长)")
	fmt.Println("  3. 快重传: 收到3个重复ACK，立即重传(不等超时)")
	fmt.Println("  4. 快恢复: ssthresh=cwnd/2, cwnd=ssthresh(不回1)")
}

func main() {
	fmt.Println("=== TCP编程与协议原理 ===")
	fmt.Println()

	startTCPServer()

	demonstrateTCPStates()
	demonstrateSlidingWindow()

	fmt.Println()
	fmt.Println("=== 考研要点 ===")
	fmt.Println("1. TCP: 面向连接、可靠传输、全双工、字节流")
	fmt.Println("2. 可靠传输: 序号、确认、超时重传、校验和")
	fmt.Println("3. 流量控制: 滑动窗口，接收方控制发送速率")
	fmt.Println("4. 拥塞控制: 慢开始+拥塞避免+快重传+快恢复")
	fmt.Println("5. 三次握手: SYN -> SYN+ACK -> ACK")
	fmt.Println("6. 四次挥手: FIN -> ACK -> FIN+ACK -> ACK")
	fmt.Println("7. TIME_WAIT: 等待2MSL，确保对方收到最终ACK")
	fmt.Println("8. TCP vs UDP: TCP可靠但慢，UDP不可靠但快")
}
