package main

import (
	"fmt"
	"math"
)

type TCPSegment struct {
	SrcPort    int
	DstPort    int
	SeqNum     int
	AckNum     int
	SYN        bool
	ACK        bool
	FIN        bool
	WindowSize int
	Data       string
}

type IPPacket struct {
	Version    int
	IHL        int
	TTL        int
	Protocol   int
	SrcIP      string
	DstIP      string
	Data       interface{}
}

type EthernetFrame struct {
	SrcMAC  string
	DstMAC  string
	Type    int
	Data    interface{}
	FCS     string
}

func encapsulateHTTP() {
	fmt.Println("--- 数据封装过程 ---")
	fmt.Println()
	fmt.Println("应用层: HTTP请求")
	httpData := "GET /index.html HTTP/1.1\r\nHost: example.com\r\n\r\n"
	fmt.Printf("  数据: %q\n", httpData)

	fmt.Println()
	fmt.Println("传输层: 添加TCP首部(20字节)")
	tcpSeg := TCPSegment{
		SrcPort: 12345, DstPort: 80,
		SeqNum: 1, AckNum: 0,
		SYN: true, ACK: false, FIN: false,
		WindowSize: 65535,
		Data:       httpData,
	}
	fmt.Printf("  源端口=%d, 目的端口=%d, SYN=%v\n", tcpSeg.SrcPort, tcpSeg.DstPort, tcpSeg.SYN)

	fmt.Println()
	fmt.Println("网络层: 添加IP首部(20字节)")
	ipPkt := IPPacket{
		Version: 4, IHL: 5, TTL: 64,
		Protocol: 6,
		SrcIP:    "192.168.1.100",
		DstIP:    "93.184.216.34",
		Data:     tcpSeg,
	}
	fmt.Printf("  源IP=%s, 目的IP=%s, TTL=%d, 协议=%d(TCP)\n",
		ipPkt.SrcIP, ipPkt.DstIP, ipPkt.TTL, ipPkt.Protocol)

	fmt.Println()
	fmt.Println("数据链路层: 添加以太网帧头尾(14+4字节)")
	ethFrame := EthernetFrame{
		SrcMAC: "AA:BB:CC:DD:EE:FF",
		DstMAC: "11:22:33:44:55:66",
		Type:   0x0800,
		Data:   ipPkt,
		FCS:    "CRC32校验",
	}
	fmt.Printf("  源MAC=%s, 目的MAC=%s, 类型=0x%04X(IPv4)\n",
		ethFrame.SrcMAC, ethFrame.DstMAC, ethFrame.Type)

	fmt.Println()
	fmt.Println("物理层: 转换为比特流在物理介质上传输")
	fmt.Println()
	fmt.Println("封装总结:")
	fmt.Println("  [以太网头][IP头][TCP头][HTTP数据][FCS]")
	fmt.Println("   14字节   20字节 20字节  变长    4字节")
}

func demonstrateOSIModel() {
	fmt.Println()
	fmt.Println("--- OSI七层模型 vs TCP/IP四层模型 ---")
	fmt.Println()
	fmt.Printf("%-4s %-16s %-16s %-20s %-15s\n", "层", "OSI", "TCP/IP", "功能", "协议/设备")
	fmt.Printf("%-4s %-16s %-16s %-20s %-15s\n",
		"7", "应用层", "应用层", "为用户提供服务", "HTTP/FTP/DNS")
	fmt.Printf("%-4s %-16s %-16s %-20s %-15s\n",
		"6", "表示层", "应用层", "数据格式转换/加密", "SSL/TLS/JPEG")
	fmt.Printf("%-4s %-16s %-16s %-20s %-15s\n",
		"5", "会话层", "应用层", "建立/管理会话", "RPC/SQL")
	fmt.Printf("%-4s %-16s %-16s %-20s %-15s\n",
		"4", "传输层", "传输层", "端到端可靠传输", "TCP/UDP")
	fmt.Printf("%-4s %-16s %-16s %-20s %-15s\n",
		"3", "网络层", "网际层", "路由选择/寻址", "IP/ICMP/路由器")
	fmt.Printf("%-4s %-16s %-16s %-20s %-15s\n",
		"2", "数据链路层", "网络接口层", "帧传输/差错控制", "以太网/交换机")
	fmt.Printf("%-4s %-16s %-16s %-20s %-15s\n",
		"1", "物理层", "网络接口层", "比特流传输", "网线/集线器")
}

func demonstrateCSMACCD() {
	fmt.Println()
	fmt.Println("--- CSMA/CD协议 ---")
	fmt.Println()
	fmt.Println("以太网介质访问控制: CSMA/CD(载波监听多路访问/冲突检测)")
	fmt.Println()
	fmt.Println("工作流程:")
	fmt.Println("  1. 先听后发: 监听信道是否空闲")
	fmt.Println("  2. 边发边听: 发送数据同时检测冲突")
	fmt.Println("  3. 冲突停止: 检测到冲突立即停止发送")
	fmt.Println("  4. 随机重发: 等待随机时间后重新尝试(二进制指数退避)")
	fmt.Println()
	fmt.Println("最小帧长 = 2 * 传播延迟 * 数据率")
	fmt.Println("  以太网最小帧长: 64字节 (10Mbps, 2km, 51.2μs)")
	fmt.Println()
	fmt.Println("争用期(2τ):")
	fmt.Println("  单程传播延迟τ = 5μs (10Mbps以太网)")
	fmt.Println("  争用期 = 2τ = 51.2μs")
	fmt.Println("  最短帧长 = 10Mbps × 51.2μs = 512bit = 64字节")
	fmt.Println()
	fmt.Println("二进制指数退避:")
	fmt.Println("  第k次冲突: 从{0,1,...,2^k-1}中随机选r")
	fmt.Println("  等待 r × 2τ 后重发")
	fmt.Println("  k最大取10，超过16次放弃")
}

func calculateChecksum() {
	fmt.Println()
	fmt.Println("--- IP首部校验和计算 ---")
	fmt.Println()

	header := []int{
		0x4500, 0x001c, 0x0000, 0x0000,
		0x4001, 0x0000, 0xC0A8, 0x0164,
		0x5DB8, 0xD822,
	}

	fmt.Println("IP首部(校验和字段置0):")
	for i, h := range header {
		fmt.Printf("  0x%04X", h)
		if (i+1)%5 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()

	sum := 0
	for _, h := range header {
		sum += h
	}
	for sum > 0xFFFF {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}
	checksum := ^sum & 0xFFFF

	fmt.Printf("累加和: 0x%04X\n", sum)
	fmt.Printf("取反得到校验和: 0x%04X\n", checksum)
	fmt.Println()
	fmt.Println("校验和算法: 反码求和")
	fmt.Println("  1. 将校验和字段置0")
	fmt.Println("  2. 将所有16位字相加")
	fmt.Println("  3. 将进位加到低位(回卷)")
	fmt.Println("  4. 取反得到校验和")
	fmt.Println("  5. 接收方: 所有字(含校验和)相加，结果全1则正确")
}

func demonstrateIPFragmentation() {
	fmt.Println()
	fmt.Println("--- IP分片 ---")
	fmt.Println()
	fmt.Println("例题: 4000字节IP数据报，MTU=1500字节，如何分片?")
	fmt.Println()

	totalLen := 4000
	ipHeader := 20
	dataLen := totalLen - ipHeader
	mtu := 1500
	maxDataPerFragment := mtu - ipHeader

	fmt.Printf("原始数据报: 总长=%d, 首部=%d, 数据=%d\n", totalLen, ipHeader, dataLen)
	fmt.Printf("MTU=%d, 每片最大数据=%d\n", mtu, maxDataPerFragment)
	fmt.Println()

	offset := 0
	fragNum := 1
	for dataLen > 0 {
		chunk := maxDataPerFragment
		if chunk > dataLen {
			chunk = dataLen
		}
		mf := 1
		if dataLen-chunk <= 0 {
			mf = 0
		}
		fmt.Printf("片%d: 数据=%d字节, 总长=%d, MF=%d, 偏移=%d\n",
			fragNum, chunk, chunk+ipHeader, mf, offset/8)
		offset += chunk
		dataLen -= chunk
		fragNum++
	}

	fmt.Println()
	fmt.Println("分片要点:")
	fmt.Println("  1. 标识(Identification): 同一数据报的分片标识相同")
	fmt.Println("  2. MF(More Fragments): MF=1表示后面还有分片")
	fmt.Println("  3. 片偏移: 以8字节为单位")
	fmt.Println("  4. 只有最后一片MF=0")
	fmt.Println("  5. 每片都要加IP首部(20字节)")
}

func main() {
	fmt.Println("=== 协议栈与网络原理 ===")
	fmt.Println()

	encapsulateHTTP()
	demonstrateOSIModel()
	demonstrateCSMACCD()
	calculateChecksum()
	demonstrateIPFragmentation()

	fmt.Println()
	fmt.Println("=== 考研要点 ===")
	fmt.Println("1. OSI七层: 应表会传网数物")
	fmt.Println("2. TCP/IP四层: 应用层、传输层、网际层、网络接口层")
	fmt.Println("3. 封装: 应用数据->TCP段->IP包->以太网帧->比特流")
	fmt.Println("4. 解封装: 反向过程，逐层去掉首部")
	fmt.Println("5. CSMA/CD: 先听后发、边发边听、冲突停止、随机重发")
	fmt.Println("6. 最小帧长64字节: 保证在发送完前能检测到冲突")
	fmt.Println("7. IP分片: 超过MTU时分片，片偏移以8字节为单位")
	fmt.Println("8. 校验和: 反码求和，IP首部校验和只校验首部")
	fmt.Println("9. MTU: 最大传输单元，以太网1500字节")
	fmt.Printf("10. IP地址分类: A(1-126) B(128-191) C(192-223) D(224-239组播) E(240-255保留)\n")
}

func init() {
	_ = math.MaxFloat64
}
