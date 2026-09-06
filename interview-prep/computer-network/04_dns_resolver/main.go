package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type DNSRecord struct {
	Domain string
	Type   string
	Value  string
	TTL    int
}

type DNSServer struct {
	Records map[string][]DNSRecord
}

func NewDNSServer() *DNSServer {
	return &DNSServer{
		Records: make(map[string][]DNSRecord),
	}
}

func (ds *DNSServer) AddRecord(domain, recordType, value string, ttl int) {
	ds.Records[domain] = append(ds.Records[domain], DNSRecord{
		Domain: domain,
		Type:   recordType,
		Value:  value,
		TTL:    ttl,
	})
}

func (ds *DNSServer) Query(domain, recordType string) []DNSRecord {
	records, ok := ds.Records[domain]
	if !ok {
		return nil
	}
	var result []DNSRecord
	for _, r := range records {
		if r.Type == recordType {
			result = append(result, r)
		}
	}
	return result
}

func (ds *DNSServer) RecursiveQuery(domain, recordType string) []DNSRecord {
	if records := ds.Query(domain, recordType); len(records) > 0 {
		return records
	}
	parts := strings.Split(domain, ".")
	for i := 1; i < len(parts); i++ {
		parentDomain := strings.Join(parts[i:], ".")
		if records := ds.Query(parentDomain, "NS"); len(records) > 0 {
			fmt.Printf("  递归查询: 委派给 %s -> %s\n", parentDomain, records[0].Value)
			if records := ds.Query(domain, recordType); len(records) > 0 {
				return records
			}
		}
	}
	return nil
}

func parseIP(ipStr string) uint32 {
	parts := strings.Split(ipStr, ".")
	var result uint32
	for _, p := range parts {
		v, _ := strconv.Atoi(p)
		result = (result << 8) | uint32(v)
	}
	return result
}

func ipToString(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		(ip>>24)&0xFF, (ip>>16)&0xFF, (ip>>8)&0xFF, ip&0xFF)
}

func subnetInfo(ipStr string, prefixLen int) {
	ip := parseIP(ipStr)
	var mask uint32 = 0xFFFFFFFF << (32 - prefixLen)
	network := ip & mask
	broadcast := network | ^mask
	firstHost := network + 1
	lastHost := broadcast - 1
	hostCount := uint32(1) << (32 - prefixLen)
	if prefixLen == 32 {
		hostCount = 1
	} else if prefixLen == 31 {
		hostCount = 2
	} else {
		hostCount -= 2
	}

	fmt.Printf("IP地址:     %s/%d\n", ipStr, prefixLen)
	fmt.Printf("子网掩码:   %s\n", ipToString(mask))
	fmt.Printf("网络地址:   %s\n", ipToString(network))
	fmt.Printf("广播地址:   %s\n", ipToString(broadcast))
	fmt.Printf("第一个主机: %s\n", ipToString(firstHost))
	fmt.Printf("最后一个主机: %s\n", ipToString(lastHost))
	fmt.Printf("主机数量:   %d\n", hostCount)
}

func demonstrateDNSResolution() {
	fmt.Println("--- DNS解析过程 ---")
	fmt.Println()
	fmt.Println("浏览器访问 www.example.com 的DNS解析过程:")
	fmt.Println()
	fmt.Println("  1. 浏览器检查本地缓存")
	fmt.Println("  2. 查询本地DNS缓存(操作系统)")
	fmt.Println("  3. 向本地域名服务器(LDNS)发送递归查询")
	fmt.Println("  4. LDNS向根域名服务器发送迭代查询")
	fmt.Println("     LDNS -> 根DNS: \"www.example.com的IP?\"")
	fmt.Println("     根DNS -> LDNS: \"去问.com的TLD服务器\"")
	fmt.Println("  5. LDNS向.com TLD服务器查询")
	fmt.Println("     LDNS -> .com TLD: \"www.example.com的IP?\"")
	fmt.Println("     .com TLD -> LDNS: \"去问example.com的权威DNS\"")
	fmt.Println("  6. LDNS向example.com权威DNS查询")
	fmt.Println("     LDNS -> 权威DNS: \"www.example.com的IP?\"")
	fmt.Println("     权威DNS -> LDNS: \"93.184.216.34\"")
	fmt.Println("  7. LDNS缓存结果并返回给浏览器")
	fmt.Println()
	fmt.Println("递归查询 vs 迭代查询:")
	fmt.Println("  递归: 客户端只问一次，服务器负责到底(客户端->LDNS)")
	fmt.Println("  迭代: 客户端多次询问不同服务器(LDNS->根->TLD->权威)")
}

func demonstrateCIDR() {
	fmt.Println()
	fmt.Println("--- CIDR与子网划分 ---")
	fmt.Println()

	fmt.Println("示例1: 192.168.1.100/24")
	subnetInfo("192.168.1.100", 24)
	fmt.Println()

	fmt.Println("示例2: 10.0.0.50/16")
	subnetInfo("10.0.0.50", 16)
	fmt.Println()

	fmt.Println("示例3: 172.16.5.30/28")
	subnetInfo("172.16.5.30", 28)
	fmt.Println()

	fmt.Println("子网划分例题:")
	fmt.Println("  某公司有4个部门，分别需要50、25、12、5台主机")
	fmt.Println("  给定192.168.1.0/24，如何划分?")
	fmt.Println()
	fmt.Println("  部门1(50台): 192.168.1.0/26   (62台主机, 掩码255.255.255.192)")
	fmt.Println("  部门2(25台): 192.168.1.64/27  (30台主机, 掩码255.255.255.224)")
	fmt.Println("  部门3(12台): 192.168.1.96/28  (14台主机, 掩码255.255.255.240)")
	fmt.Println("  部门4(5台):  192.168.1.112/29 (6台主机,  掩码255.255.255.248)")
}

func main() {
	fmt.Println("=== DNS与网络层 ===")
	fmt.Println()

	fmt.Println("--- 模拟DNS服务器 ---")
	dns := NewDNSServer()
	dns.AddRecord("example.com", "A", "93.184.216.34", 3600)
	dns.AddRecord("example.com", "MX", "mail.example.com", 3600)
	dns.AddRecord("www.example.com", "A", "93.184.216.34", 3600)
	dns.AddRecord("mail.example.com", "A", "93.184.216.35", 3600)
	dns.AddRecord("com", "NS", "a.gtld-servers.net", 86400)
	dns.AddRecord(".", "NS", "a.root-servers.net", 86400)

	fmt.Println("查询 www.example.com A记录:")
	if records := dns.Query("www.example.com", "A"); len(records) > 0 {
		for _, r := range records {
			fmt.Printf("  %s %s %s (TTL=%d)\n", r.Domain, r.Type, r.Value, r.TTL)
		}
	}
	fmt.Println()

	fmt.Println("查询 example.com MX记录:")
	if records := dns.Query("example.com", "MX"); len(records) > 0 {
		for _, r := range records {
			fmt.Printf("  %s %s %s (TTL=%d)\n", r.Domain, r.Type, r.Value, r.TTL)
		}
	}
	fmt.Println()

	demonstrateDNSResolution()
	demonstrateCIDR()

	fmt.Println()
	fmt.Println("=== 考研要点 ===")
	fmt.Println("1. DNS: 域名系统，将域名映射为IP地址")
	fmt.Println("2. DNS记录类型: A(IPv4), AAAA(IPv6), CNAME(别名), MX(邮件), NS(域名服务器)")
	fmt.Println("3. 递归查询: 客户端->本地DNS，服务器负责到底")
	fmt.Println("4. 迭代查询: 本地DNS->根DNS->TLD->权威DNS，逐级查询")
	fmt.Println("5. CIDR: 无类域间路由，a.b.c.d/n，n为网络前缀长度")
	fmt.Println("6. 子网掩码: 1对应网络号，0对应主机号")
	fmt.Println("7. 网络地址: IP & 掩码，主机号全0")
	fmt.Println("8. 广播地址: 网络地址 | ~掩码，主机号全1")
	fmt.Println("9. ARP: IP地址 -> MAC地址，同一子网内解析")
	fmt.Println("10. DHCP: 自动分配IP地址，四步: Discover/Offer/Request/Ack")
}

func init() {
	_ = net.Dial
}
