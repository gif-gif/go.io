package goip

import (
	"errors"
	"fmt"
	"time"

	"github.com/miekg/dns"
)

// 只返回A记录IPS
func LookupIPList(domain string, dnsServer string) ([]string, error) {
	c := new(dns.Client)
	c.Timeout = 5 * time.Second

	m := new(dns.Msg)
	// 设置查询域名，TypeA 表示查询 IPv4
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	m.RecursionDesired = true

	// 执行查询
	r, _, err := c.Exchange(m, dnsServer)
	if err != nil {
		return []string{}, err
	}

	if r.Rcode != dns.RcodeSuccess {
		fmt.Printf("DNS 返回错误码: %v\n", r.Rcode)
		return []string{}, errors.New(fmt.Sprintf("DNS 返回错误码: %v\n", r.Rcode))
	}

	ips := []string{}
	fmt.Printf("来自 %s 的结果:\n", dnsServer)
	for _, ans := range r.Answer {
		// 类型断言为 A 记录
		if a, ok := ans.(*dns.A); ok {
			ips = append(ips, a.A.String())
		}
	}
	return ips, nil
}
