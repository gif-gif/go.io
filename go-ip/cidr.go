package goip

import (
	"encoding/binary"
	"fmt"
	"net"
)

type CIDR struct {
	Network   string
	FirstIP   string
	LastIP    string
	Total     uint64
	Netmask   string
	Wildcard  string
	IPVersion string
}

// PrintIPRange 打印指定 CIDR 范围内的所有 IP
func PrintIPRange(cidrStr string) ([]string, error) {
	_, ipnet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR: %v", err)
	}

	// 获取起始 IP
	ip := ipnet.IP

	// 将 IP 转换为 4 字节表示（针对 IPv4）
	ipStart := binary.BigEndian.Uint32(ip.To4())

	// 计算掩码
	mask := binary.BigEndian.Uint32(ipnet.Mask)

	// 计算最后一个 IP
	ipEnd := (ipStart & mask) | (^mask)

	ips := []string{}
	// 打印范围内的所有 IP
	for i := ipStart; i <= ipEnd; i++ {
		// 转换回 IP 地址格式
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		fmt.Println(ip)
		ips = append(ips, ip.String())
	}

	return ips, nil
}

func CalculateCIDR(cidrStr string) (*CIDR, error) {
	_, ipnet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR format: %v", err)
	}

	cidr := &CIDR{
		Network: cidrStr,
	}

	if len(ipnet.IP) == net.IPv6len {
		cidr.IPVersion = "IPv6"
	} else {
		cidr.IPVersion = "IPv4"
	}

	cidr.FirstIP = ipnet.IP.String()

	mask := net.IP(ipnet.Mask)
	lastIP := make(net.IP, len(ipnet.IP))
	for i := 0; i < len(ipnet.IP); i++ {
		lastIP[i] = ipnet.IP[i] | ^mask[i]
	}
	cidr.LastIP = lastIP.String()

	ones, bits := ipnet.Mask.Size()
	cidr.Total = 1 << uint64(bits-ones)

	cidr.Netmask = net.IP(ipnet.Mask).String()

	wildcard := make(net.IP, len(ipnet.Mask))
	for i := 0; i < len(ipnet.Mask); i++ {
		wildcard[i] = ^ipnet.Mask[i]
	}
	cidr.Wildcard = net.IP(wildcard).String()

	return cidr, nil
}

func IsIPInCIDR(ipStr, cidrStr string) (bool, error) {
	// 解析 IP 地址
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false, fmt.Errorf("invalid IP address: %s", ipStr)
	}

	// 解析 CIDR
	_, ipnet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return false, fmt.Errorf("invalid CIDR: %s, error: %v", cidrStr, err)
	}

	// 检查 IP 是否在范围内
	return ipnet.Contains(ip), nil
}
