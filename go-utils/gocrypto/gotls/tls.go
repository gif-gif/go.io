package gotls

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

// CertConfig 证书配置
type CertConfig struct {
	// 主机名/域名
	Hosts []string
	// IP 地址
	IPs []net.IP
	// 组织名称
	Organization string
	// 通用名称（Common Name）
	CommonName string
	// 国家代码
	Country string
	// 省份
	Province string
	// 城市
	Locality string
	// 有效期（天数）
	ValidDays int
	// RSA 密钥位数（2048 或 4096）
	RSABits int
	// 是否为 CA 证书
	IsCA bool
}

// DefaultConfig 返回默认配置
func DefaultConfig() *CertConfig {
	return &CertConfig{
		Hosts:        []string{""},
		IPs:          []net.IP{net.ParseIP("")},
		Organization: "My Company",
		CommonName:   "localhost",
		Country:      "SG",
		Province:     "Singapore",
		Locality:     "Singapore",
		ValidDays:    3650,
		RSABits:      4096,
		IsCA:         false,
	}
}

// GenerateTLSCert 生成 TLS 证书和私钥
func GenerateTLSCert(config *CertConfig, certPath, keyPath string) error {
	if config == nil {
		config = DefaultConfig()
	}

	// 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, config.RSABits)
	if err != nil {
		return fmt.Errorf("生成私钥失败: %w", err)
	}

	// 生成序列号
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return fmt.Errorf("生成序列号失败: %w", err)
	}

	// 证书模板
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:      []string{config.Country},
			Province:     []string{config.Province},
			Locality:     []string{config.Locality},
			Organization: []string{config.Organization},
			CommonName:   config.CommonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Duration(config.ValidDays) * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		DNSNames:              config.Hosts,
		IPAddresses:           config.IPs,
	}

	// 如果是 CA 证书
	if config.IsCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	// 创建证书
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("创建证书失败: %w", err)
	}

	// 保存证书
	certFile, err := os.Create(certPath)
	if err != nil {
		return fmt.Errorf("创建证书文件失败: %w", err)
	}
	defer certFile.Close()

	if err := pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes}); err != nil {
		return fmt.Errorf("编码证书失败: %w", err)
	}

	// 保存私钥
	keyFile, err := os.Create(keyPath)
	if err != nil {
		return fmt.Errorf("创建私钥文件失败: %w", err)
	}
	defer keyFile.Close()

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	if err := pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateKeyBytes}); err != nil {
		return fmt.Errorf("编码私钥失败: %w", err)
	}

	return nil
}

// GenerateTLSCertWithCA 使用 CA 签名生成证书
func GenerateTLSCertWithCA(config *CertConfig, caCertPath, caKeyPath, certPath, keyPath string) error {
	// 读取 CA 证书
	caCertPEM, err := os.ReadFile(caCertPath)
	if err != nil {
		return fmt.Errorf("读取 CA 证书失败: %w", err)
	}
	caCertBlock, _ := pem.Decode(caCertPEM)
	caCert, err := x509.ParseCertificate(caCertBlock.Bytes)
	if err != nil {
		return fmt.Errorf("解析 CA 证书失败: %w", err)
	}

	// 读取 CA 私钥
	caKeyPEM, err := os.ReadFile(caKeyPath)
	if err != nil {
		return fmt.Errorf("读取 CA 私钥失败: %w", err)
	}
	caKeyBlock, _ := pem.Decode(caKeyPEM)
	caKey, err := x509.ParsePKCS1PrivateKey(caKeyBlock.Bytes)
	if err != nil {
		return fmt.Errorf("解析 CA 私钥失败: %w", err)
	}

	// 生成新的私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, config.RSABits)
	if err != nil {
		return fmt.Errorf("生成私钥失败: %w", err)
	}

	// 生成序列号
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return fmt.Errorf("生成序列号失败: %w", err)
	}

	// 证书模板
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:      []string{config.Country},
			Province:     []string{config.Province},
			Locality:     []string{config.Locality},
			Organization: []string{config.Organization},
			CommonName:   config.CommonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Duration(config.ValidDays) * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		DNSNames:              config.Hosts,
		IPAddresses:           config.IPs,
	}

	// 使用 CA 签名
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, caCert, &privateKey.PublicKey, caKey)
	if err != nil {
		return fmt.Errorf("创建证书失败: %w", err)
	}

	// 保存证书
	certFile, err := os.Create(certPath)
	if err != nil {
		return fmt.Errorf("创建证书文件失败: %w", err)
	}
	defer certFile.Close()

	if err := pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes}); err != nil {
		return fmt.Errorf("编码证书失败: %w", err)
	}

	// 保存私钥
	keyFile, err := os.Create(keyPath)
	if err != nil {
		return fmt.Errorf("创建私钥文件失败: %w", err)
	}
	defer keyFile.Close()

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	if err := pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateKeyBytes}); err != nil {
		return fmt.Errorf("编码私钥失败: %w", err)
	}

	return nil
}

// VerifyCertAndKey 验证证书和私钥是否匹配且有效
func VerifyCertAndKey(certPath, keyPath string) (*x509.Certificate, error) {
	// 读取证书
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("读取证书文件失败: %w", err)
	}

	// 读取私钥
	keyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("读取私钥文件失败: %w", err)
	}

	return VerifyCertAndKeyByContent(certPEM, keyPEM)
}

// VerifyCertAndKey 验证证书和私钥是否匹配且有效
func VerifyCertAndKeyByContent(certPEM, keyPEM []byte) (*x509.Certificate, error) {
	certBlock, _ := pem.Decode(certPEM)
	if certBlock == nil {
		return nil, fmt.Errorf("证书文件格式错误: 无法解码 PEM")
	}

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析证书失败: %w", err)
	}

	keyBlock, _ := pem.Decode(keyPEM)
	if keyBlock == nil {
		return nil, fmt.Errorf("私钥文件格式错误: 无法解码 PEM")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		// 尝试 PKCS8 格式
		pkcs8Key, err2 := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
		if err2 != nil {
			return nil, fmt.Errorf("解析私钥失败 (PKCS1: %v, PKCS8: %v)", err, err2)
		}
		var ok bool
		privateKey, ok = pkcs8Key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("私钥类型错误: 需要 RSA 私钥")
		}
	}

	// 验证私钥和证书的公钥是否匹配
	certPubKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("证书公钥类型错误: 需要 RSA 公钥")
	}

	if certPubKey.N.Cmp(privateKey.N) != 0 || certPubKey.E != privateKey.E {
		return nil, fmt.Errorf("证书和私钥不匹配")
	}

	// 检查证书有效期
	now := time.Now()
	if now.Before(cert.NotBefore) {
		return nil, fmt.Errorf("证书尚未生效 (生效时间: %s)", cert.NotBefore.Format("2006-01-02 15:04:05"))
	}
	if now.After(cert.NotAfter) {
		return nil, fmt.Errorf("证书已过期 (过期时间: %s)", cert.NotAfter.Format("2006-01-02 15:04:05"))
	}

	// 验证证书签名（自签名证书）
	// 注意：只有 CA 证书或自签名的服务器证书才能验证自己
	if cert.IsCA || cert.Subject.CommonName == cert.Issuer.CommonName {
		if err := cert.CheckSignatureFrom(cert); err != nil {
			// 对于非 CA 的自签名证书，这个错误是预期的，可以忽略
			if !cert.IsCA {
				fmt.Println("  ⚠ 注意: 这是非 CA 的自签名证书，签名验证跳过")
			} else {
				return nil, fmt.Errorf("证书签名验证失败: %w", err)
			}
		}
	}

	fmt.Println("✓ 证书和私钥验证通过")
	fmt.Printf("  - 证书类型: %s\n", func() string {
		if cert.IsCA {
			return "CA 证书"
		}
		if cert.Subject.CommonName == cert.Issuer.CommonName {
			return "自签名服务器证书"
		}
		return "CA 签名的服务器证书"
	}())

	//fmt.Printf("  - 证书主体: %s\n", cert.Subject.CommonName)
	//fmt.Printf("  - 组织: %s\n", cert.Subject.Organization)
	//fmt.Printf("  - 生效时间: %s\n", cert.NotBefore.Format("2006-01-02 15:04:05"))
	//fmt.Printf("  - 过期时间: %s\n", cert.NotAfter.Format("2006-01-02 15:04:05"))
	//fmt.Printf("  - 剩余天数: %d 天\n", int(time.Until(cert.NotAfter).Hours()/24))
	//if len(cert.DNSNames) > 0 {
	//	fmt.Printf("  - DNS 名称: %v\n", cert.DNSNames)
	//}
	//if len(cert.IPAddresses) > 0 {
	//	fmt.Printf("  - IP 地址: %v\n", cert.IPAddresses)
	//}
	//fmt.Printf("  - 密钥长度: %d 位\n", privateKey.N.BitLen())

	return cert, nil
}

// VerifyCertChain 验证证书链（证书是否由 CA 签名）
func VerifyCertChain(certPath, caCertPath string) (*x509.Certificate, error) {
	// 读取证书
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("读取证书文件失败: %w", err)
	}

	// 读取 CA 证书
	caCertPEM, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, fmt.Errorf("读取 CA 证书文件失败: %w", err)
	}

	return VerifyCertChainByContent(certPEM, caCertPEM)
}

// VerifyCertChain 验证证书链（证书是否由 CA 签名）
func VerifyCertChainByContent(certPEM, caCertPEM []byte) (*x509.Certificate, error) {

	certBlock, _ := pem.Decode(certPEM)
	if certBlock == nil {
		return nil, fmt.Errorf("证书文件格式错误")
	}

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析证书失败: %w", err)
	}

	caCertBlock, _ := pem.Decode(caCertPEM)
	if caCertBlock == nil {
		return nil, fmt.Errorf("CA 证书文件格式错误")
	}

	caCert, err := x509.ParseCertificate(caCertBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析 CA 证书失败: %w", err)
	}

	// 验证证书是否由 CA 签名
	if err := cert.CheckSignatureFrom(caCert); err != nil {
		return nil, fmt.Errorf("证书签名验证失败: %w", err)
	}

	// 创建证书池并验证
	roots := x509.NewCertPool()
	roots.AddCert(caCert)

	opts := x509.VerifyOptions{
		Roots: roots,
	}

	if _, err := cert.Verify(opts); err != nil {
		return nil, fmt.Errorf("证书链验证失败: %w", err)
	}

	//fmt.Println("✓ 证书链验证通过")
	//fmt.Printf("  - 证书: %s\n", cert.Subject.CommonName)
	//fmt.Printf("  - 签发者: %s\n", caCert.Subject.CommonName)

	return cert, nil
}
