package gojwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	Name               string `yaml:"Name" json:"name,optional"`
	AccessSecret       string `yaml:"AccessSecret" json:"accessSecret,optional"`
	AccessExpire       int64  `yaml:"AccessExpire" json:"accessExpire,optional"`
	RefreshTokenExpire int64  `yaml:"RefreshTokenExpire" json:"refreshTokenExpire,optional"`
	JwtEncryptSecret   string `yaml:"JwtEncryptSecret" json:"jwtEncryptSecret,optional"` //Token加密秘钥
}

type GoJwt struct {
	Config Config
}

func New(config Config) *GoJwt {
	return &GoJwt{
		Config: config,
	}
}

func GetJwtToken(secretKey string, iat, seconds int64, params map[string]any) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	for key, val := range params {
		claims[key] = val
	}
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

func (c *GoJwt) GeneratedTokens(params map[string]any) (accessToken, refreshToken string, expires int64, err error) {
	now := time.Now()
	expires = now.Add(time.Duration(c.Config.AccessExpire) * time.Second).Unix()
	nowUnix := now.Unix()
	accessToken, err = GetJwtToken(c.Config.AccessSecret, nowUnix, c.Config.AccessExpire, params)
	if err != nil {
		return
	}
	params["isRefreshToken"] = 1
	refreshToken, err = GetJwtToken(c.Config.AccessSecret, nowUnix, c.Config.RefreshTokenExpire, params)
	if err != nil {
		return
	}
	return
}

// 解析token
func (c *GoJwt) ParseToken(tokenString string) (map[string]interface{}, error) {
	if strings.Contains(tokenString, "Bearer ") {
		tokenString = tokenString[7:]
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(c.Config.AccessSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func (c *GoJwt) ParseTokenEx(tokenString string) (map[string]interface{}, bool, error) {
	// 1. 去除空格
	tokenString = strings.TrimSpace(tokenString)

	// 2. 检查空字符串
	if tokenString == "" {
		return nil, false, errors.New("token is empty")
	}

	// 3. 处理 Bearer 前缀
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimSpace(tokenString[7:])
		if tokenString == "" {
			return nil, false, errors.New("token is empty after Bearer prefix")
		}
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(c.Config.AccessSecret), nil
	})

	if err != nil {
		return nil, false, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims, true, nil
	} else {
		return nil, false, err
	}
}

func (c *GoJwt) IsValidToken(tokenString string) bool {
	// 1. 去除空格
	tokenString = strings.TrimSpace(tokenString)

	// 2. 检查空字符串
	if tokenString == "" {
		return false
	}

	// 3. 处理 Bearer 前缀
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimSpace(tokenString[7:])
		if tokenString == "" {
			return false
		}
	}

	// 4. 解析并验证 token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(c.Config.AccessSecret), nil
	})

	// 5. 检查错误和有效性
	if err != nil {
		return false
	}

	return token.Valid
}

// 验证Refresh Token并生成新的Access Token
func (c *GoJwt) RefreshAccessToken(refreshTokenStr string) (accessToken, refreshToken string, expires int64, err error) {
	// 1. 去除空格
	refreshTokenStr = strings.TrimSpace(refreshTokenStr)

	// 2. 检查空字符串
	if refreshTokenStr == "" {
		return "", "", 0, errors.New("refresh token is empty")
	}

	// 3. 处理 Bearer 前缀
	if strings.HasPrefix(refreshTokenStr, "Bearer ") {
		refreshTokenStr = strings.TrimSpace(refreshTokenStr[7:])
		if refreshTokenStr == "" {
			return "", "", 0, errors.New("refresh token is empty")
		}
	}

	params, ok, err := c.ParseTokenEx(refreshTokenStr)
	if err != nil {
		return "", "", 0, err
	}
	if !ok {
		return "", "", 0, fmt.Errorf("invalid refresh token")
	}
	return c.GeneratedTokens(params)
}
