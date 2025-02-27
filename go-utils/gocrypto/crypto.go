package gocrypto

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"math/big"
	"net/url"
	"os"
	"strings"
)

//md5
//sha1
//sha256
//sha512

const (
	HashingAlgorithmMd5    = "md5"
	HashingAlgorithmSha1   = "sha1"
	HashingAlgorithmSha256 = "sha256"
	HashingAlgorithmSha512 = "sha512"
)

//hex（十六进制）
//base64

const (
	EncodingHex    = "hex"
	EncodingBase64 = "base64"
)

// CreateHash 计算文件的哈希值
// data: 输入数据
// hashingAlgorithm: 哈希算法名称 ("md5", "sha1", "sha256", "sha512" 等)
// encoding: 编码方式 ("hex" 或 "base64")
func CreateHash(data []byte, hashingAlgorithm string, encoding string) (string, error) {
	// 选择哈希算法
	var h hash.Hash
	switch hashingAlgorithm {
	case HashingAlgorithmMd5:
		h = md5.New()
	case HashingAlgorithmSha1:
		h = sha1.New()
	case HashingAlgorithmSha256:
		h = sha256.New()
	case HashingAlgorithmSha512:
		h = sha512.New()
	default:
		return "", fmt.Errorf("unsupported hashing algorithm: %s", hashingAlgorithm)
	}

	// 写入数据
	h.Write(data)

	// 获取哈希值
	hashBytes := h.Sum(nil)

	// 根据指定编码格式返回结果
	switch encoding {
	case EncodingHex:
		return hex.EncodeToString(hashBytes), nil
	case EncodingBase64:
		return base64.StdEncoding.EncodeToString(hashBytes), nil
	default:
		return "", fmt.Errorf("unsupported encoding: %s", encoding)
	}
}

// 生成随机字符串，可用于生成随机密钥， len 为长度
func GenerateKey(len int64) (string, error) {
	// 生成32字节（256位）的密钥
	key := make([]byte, len)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}

// 生成 AES 密钥
func GenerateAESKey() (string, error) {
	// 生成32字节（256位）的密钥
	key, err := GenerateKey(32)
	if err != nil {
		return "", err
	}
	return key, nil
}

// 生成 AES 密钥和 IV
func GenerateAESKeyAndIV() (string, string, error) {
	// 生成 16 字节（128 位）的 Key
	key, err := GenerateKey(16)
	if err != nil {
		return "", "", err
	}

	// 生成 16 字节（128 位）的 IV
	iv, err := GenerateKey(32)
	if err != nil {
		return "", "", err
	}

	return key, iv, nil
}

// 计算文件md5(支持超大文件)
func CalculateFileMD5(filePath string) (string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 创建MD5哈希对象
	hash := md5.New()

	// 创建一个缓冲区，逐块读取文件内容
	buffer := make([]byte, 1024*1024) // 1MB 缓冲区
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return "", err
		}
		if n == 0 {
			break
		}
		// 更新哈希值
		if _, err := hash.Write(buffer[:n]); err != nil {
			return "", err
		}
	}

	// 计算最终的哈希值
	hashInBytes := hash.Sum(nil)
	hashInString := fmt.Sprintf("%x", hashInBytes)

	return hashInString, nil
}

// MD5 大写
func MD5(buf []byte) string {
	h := md5.New()
	h.Write(buf)
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

// Md5小写
func Md5(buf []byte) string {
	h := md5.New()
	h.Write(buf)
	return strings.ToLower(hex.EncodeToString(h.Sum(nil)))
}

// 小写
func SHA1(buf []byte) string {
	h := sha1.New()
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil))
}

func SHA256(buf, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil))
}

func HMacMd5(buf, key []byte) string {
	h := hmac.New(md5.New, key)
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil))
}

func HMacSha1(buf, key []byte) string {
	h := hmac.New(sha1.New, key)
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil))
}

func HMacSha256(buf, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil))
}

func Base64Encode(buf []byte) string {
	return base64.StdEncoding.EncodeToString(buf)
}

func Base64Decode(str string) []byte {
	var count = (4 - len(str)%4) % 4
	str += strings.Repeat("=", count)
	buf, _ := base64.StdEncoding.DecodeString(str)
	return buf
}

func SHAWithRSA(key, data []byte) (string, error) {
	pkey, err := x509.ParsePKCS8PrivateKey(key)
	if err != nil {
		return "", err
	}

	h := crypto.Hash.New(crypto.SHA1)
	h.Write(data)
	hashed := h.Sum(nil)

	buf, err := rsa.SignPKCS1v15(rand.Reader, pkey.(*rsa.PrivateKey), crypto.SHA1, hashed)
	if err != nil {
		return "", err
	}
	return Base64Encode(buf), nil
}

func AESECBEncrypt(data, key []byte) ([]byte, error) {
	cb, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := cb.BlockSize()
	paddingSize := blockSize - len(data)%blockSize
	if paddingSize != 0 {
		data = append(data, bytes.Repeat([]byte{byte(0)}, paddingSize)...)
	}
	encrypted := make([]byte, len(data))
	for bs, be := 0, blockSize; bs < len(data); bs, be = bs+blockSize, be+blockSize {
		cb.Encrypt(encrypted[bs:be], data[bs:be])
	}
	return encrypted, nil
}

func AESECBDecrypt(buf, key []byte) ([]byte, error) {
	cb, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := cb.BlockSize()
	decrypted := make([]byte, len(buf))
	for bs, be := 0, blockSize; bs < len(buf); bs, be = bs+blockSize, be+blockSize {
		cb.Decrypt(decrypted[bs:be], buf[bs:be])
	}
	paddingSize := int(decrypted[len(decrypted)-1])
	return decrypted[0 : len(decrypted)-paddingSize], nil
}

func AESCBCEncrypt(rawData, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// block 大小 16
	blockSize := block.BlockSize()

	// 填充原文
	rawData = pkcs7padding(rawData, blockSize)

	// 定义密码数据
	var cipherData []byte

	// 如果iv为空，生成随机iv，并附加到加密数据前面，否则单独生成加密数据
	if iv == nil {
		// 初始化加密数据
		cipherData = make([]byte, blockSize+len(rawData))
		// 定义向量
		iv = cipherData[:blockSize]
		// 填充向量IV， ReadFull从rand.Reader精确地读取len(b)字节数据填充进iv，rand.Reader是一个全局、共享的密码用强随机数生成器
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return nil, err
		}
		// 加密
		mode := cipher.NewCBCEncrypter(block, iv)
		mode.CryptBlocks(cipherData[blockSize:], rawData)
	} else {
		// 初始化加密数据
		cipherData = make([]byte, len(rawData))
		// 定义向量
		iv = iv[:blockSize]
		// 加密
		mode := cipher.NewCBCEncrypter(block, iv)
		mode.CryptBlocks(cipherData, rawData)
	}

	return cipherData, nil
}

func AESCBCDecrypt(cipherData, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// block 大小 16
	blockSize := block.BlockSize()

	// 加密串长度
	l := len(cipherData)

	// 校验长度
	if l < blockSize {
		return nil, errors.New("encrypt data too short")
	}

	// 定义原始数据
	var origData []byte

	// 如果iv为空，需要获取前16位作为随机iv
	if iv == nil {
		// 定义向量
		iv = cipherData[:blockSize]
		// 定义真实加密串
		cipherData = cipherData[blockSize:]
		// 初始化原始数据
		origData = make([]byte, l-blockSize)
	} else {
		// 定义向量
		iv = iv[:blockSize]
		// 初始化原始数据
		origData = make([]byte, l)
	}

	// 解密
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(origData, cipherData)
	origData = pkcs7unpadding(origData)

	return origData, nil
}

func pkcs7padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func pkcs7unpadding(origData []byte) []byte {
	l := len(origData)
	unPadding := int(origData[l-1])
	if l < unPadding {
		return nil
	}
	return origData[:(l - unPadding)]
}

func SessionId() string {
	buf := make([]byte, 32)
	_, err := rand.Read(buf)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(buf)
}

const (
	base59key = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ."
)

// 如果遇到特殊字符，需要用 url.PathEscape(str) 解决
func Base59Encoding(strByte []byte, key ...string) string {
	strByte = []byte(url.PathEscape(string(strByte)))
	if l := len(key); l == 0 || key[0] == "" {
		key = []string{base59key}
	}
	base := int64(59)
	strTen := big.NewInt(0).SetBytes(strByte)
	keyByte := []byte(key[0])
	var modSlice []byte
	for strTen.Cmp(big.NewInt(0)) > 0 {
		mod := big.NewInt(0)
		strTen5 := big.NewInt(base)
		strTen.DivMod(strTen, strTen5, mod)
		modSlice = append(modSlice, keyByte[mod.Int64()])
	}
	for _, elem := range strByte {
		if elem != 0 {
			break
		}
		if elem == 0 {
			modSlice = append(modSlice, byte('1'))
		}
	}
	ReverseModSlice := reverseByteArr(modSlice)
	return string(ReverseModSlice)
}

func reverseByteArr(bytes []byte) []byte {
	for i := 0; i < len(bytes)/2; i++ {
		bytes[i], bytes[len(bytes)-1-i] = bytes[len(bytes)-1-i], bytes[i]
	}
	return bytes
}

func Base59Decoding(strByte []byte, key ...string) []byte {
	if l := len(key); l == 0 || key[0] == "" {
		key = []string{base59key}
	}
	base := int64(59)
	ret := big.NewInt(0)
	for _, byteElem := range strByte {
		index := bytes.IndexByte([]byte(key[0]), byteElem)
		ret.Mul(ret, big.NewInt(base))
		ret.Add(ret, big.NewInt(int64(index)))
	}
	str, _ := url.PathUnescape(string(ret.Bytes()))
	return []byte(str)
}

func UrlEncode(str string) string {
	return url.QueryEscape(str)
}

func UrlDecode(str string) string {
	str, _ = url.QueryUnescape(str)
	return str
}
