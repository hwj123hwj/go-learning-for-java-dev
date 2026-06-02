package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

func startHTTPServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("  [服务器] %s %s %s\n", r.Method, r.URL.Path, r.Proto)
		fmt.Printf("  [服务器] Host: %s\n", r.Host)
		fmt.Printf("  [服务器] User-Agent: %s\n", r.Header.Get("User-Agent"))
		fmt.Printf("  [服务器] Accept: %s\n", r.Header.Get("Accept"))
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Custom-Header", "go-learning")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from Go HTTP Server!"))
	})

	mux.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"hello","status":"ok"}`))
	})

	mux.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusMovedPermanently)
	})

	mux.HandleFunc("/notfound", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 Not Found"))
	})

	server := &http.Server{
		Addr:    "127.0.0.1:0",
		Handler: mux,
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("HTTP服务器启动失败:", err)
		return
	}
	addr := listener.Addr().String()
	fmt.Printf("HTTP服务器启动: http://%s\n", addr)

	go server.Serve(listener)
	time.Sleep(100 * time.Millisecond)

	runHTTPClient(addr)

	server.Close()
}

func runHTTPClient(addr string) {
	fmt.Println()
	fmt.Println("--- HTTP请求/响应 ---")

	urls := []string{
		fmt.Sprintf("http://%s/", addr),
		fmt.Sprintf("http://%s/json", addr),
		fmt.Sprintf("http://%s/redirect", addr),
		fmt.Sprintf("http://%s/notfound", addr),
	}

	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("请求失败: %s -> %v\n", url, err)
			continue
		}

		fmt.Printf("\n请求: GET %s\n", url)
		fmt.Printf("响应状态: %d %s\n", resp.StatusCode, resp.Status)
		fmt.Printf("响应头:")
		for k, v := range resp.Header {
			fmt.Printf("  %s: %s\n", k, strings.Join(v, ", "))
		}

		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			fmt.Printf("重定向到: %s\n", resp.Header.Get("Location"))
		}

		if resp.StatusCode == http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("响应体: %s\n", string(body))
		}
		resp.Body.Close()
	}
}

func demonstrateHTTPProtocol() {
	fmt.Println()
	fmt.Println("--- HTTP协议详解 ---")
	fmt.Println()
	fmt.Println("HTTP请求报文:")
	fmt.Println("  GET /index.html HTTP/1.1        <- 请求行(方法 URL 版本)")
	fmt.Println("  Host: www.example.com          <- 首部行")
	fmt.Println("  User-Agent: Mozilla/5.0")
	fmt.Println("  Accept: text/html")
	fmt.Println("  Connection: keep-alive")
	fmt.Println("                                  <- 空行")
	fmt.Println("  (请求体，GET通常为空)")
	fmt.Println()
	fmt.Println("HTTP响应报文:")
	fmt.Println("  HTTP/1.1 200 OK                 <- 状态行(版本 状态码 短语)")
	fmt.Println("  Content-Type: text/html          <- 首部行")
	fmt.Println("  Content-Length: 1234")
	fmt.Println("  Connection: keep-alive")
	fmt.Println("                                  <- 空行")
	fmt.Println("  <html>...</html>                <- 实体体")
	fmt.Println()
	fmt.Println("--- HTTP状态码 ---")
	fmt.Println("  1xx: 信息性   100 Continue")
	fmt.Println("  2xx: 成功     200 OK, 201 Created, 204 No Content")
	fmt.Println("  3xx: 重定向   301 永久, 302 临时, 304 未修改")
	fmt.Println("  4xx: 客户端错 400 Bad Request, 401 未授权, 403 禁止, 404 未找到")
	fmt.Println("  5xx: 服务器错 500 内部错误, 502 网关错误, 503 服务不可用")
	fmt.Println()
	fmt.Println("--- HTTP方法 ---")
	fmt.Println("  GET:    获取资源(幂等，可缓存)")
	fmt.Println("  POST:   提交数据(非幂等)")
	fmt.Println("  PUT:    替换资源(幂等)")
	fmt.Println("  DELETE: 删除资源(幂等)")
	fmt.Println("  HEAD:   获取首部(无响应体)")
	fmt.Println("  OPTIONS: 获取支持的方法")
	fmt.Println()
	fmt.Println("--- HTTP/1.0 vs HTTP/1.1 vs HTTP/2 ---")
	fmt.Println("  HTTP/1.0: 每次请求新建连接")
	fmt.Println("  HTTP/1.1: 持久连接(keep-alive)，管道化")
	fmt.Println("  HTTP/2:   多路复用，头部压缩，服务器推送")
	fmt.Println("  HTTPS:    HTTP + TLS，端口443")
}

func demonstrateRawHTTP() {
	fmt.Println()
	fmt.Println("--- 原始HTTP请求(手动构造) ---")

	conn, err := net.DialTimeout("tcp", "example.com:80", 5*time.Second)
	if err != nil {
		fmt.Println("连接失败:", err)
		return
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	request := "GET / HTTP/1.1\r\nHost: example.com\r\nConnection: close\r\n\r\n"
	fmt.Printf("发送原始请求:\n%s\n", request)

	_, err = conn.Write([]byte(request))
	if err != nil {
		fmt.Println("发送失败:", err)
		return
	}

	reader := bufio.NewReader(conn)
	lineCount := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if lineCount < 15 {
			fmt.Printf("  %s", line)
		}
		lineCount++
	}
	if lineCount > 15 {
		fmt.Printf("  ... (共%d行)\n", lineCount)
	}
}

func main() {
	fmt.Println("=== HTTP协议编程与原理 ===")
	fmt.Println()

	startHTTPServer()
	demonstrateHTTPProtocol()
	demonstrateRawHTTP()

	fmt.Println()
	fmt.Println("=== 考研要点 ===")
	fmt.Println("1. HTTP: 无状态、应用层协议、基于TCP")
	fmt.Println("2. 请求报文: 请求行 + 首部行 + 空行 + 实体体")
	fmt.Println("3. 响应报文: 状态行 + 首部行 + 空行 + 实体体")
	fmt.Println("4. Cookie: 服务器Set-Cookie，客户端Cookie，解决无状态问题")
	fmt.Println("5. 持久连接: HTTP/1.1默认keep-alive，减少连接开销")
	fmt.Println("6. HTTPS: HTTP+TLS，加密+认证+完整性")
	fmt.Println("7. HTTP/2: 二进制分帧、多路复用、头部压缩")
	fmt.Println("8. 幂等性: GET/PUT/DELETE幂等，POST非幂等")
}
