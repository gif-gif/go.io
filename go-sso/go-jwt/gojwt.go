package gojwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

type Config struct {
	Name               string `yaml:"Name" json:"name,optional"`
	AccessSecret       string `yaml:"AccessSecret" json:"accessSecret,optional"`
	AccessExpire       int64  `yaml:"AccessExpire" json:"accessExpire,optional"`
	RefreshTokenExpire int64  `yaml:"RefreshTokenExpire" json:"refreshTokenExpire,optional"`
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

func (c *GoJwt) IsValidToken(tokenString string) bool {
	if strings.Contains(tokenString, "Bearer ") {
		tokenString = strings.TrimSpace(tokenString[7:])
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(c.Config.AccessSecret), nil
	})

	if err != nil {
		return false
	}
	return token.Valid
}

// 验证Refresh Token并生成新的Access Token
func (c *GoJwt) RefreshAccessToken(refreshTokenStr string) (accessToken, refreshToken string, expires int64, err error) {
	if strings.Contains(refreshTokenStr, "Bearer ") {
		refreshTokenStr = refreshTokenStr[7:]
	}
	if !c.IsValidToken(refreshTokenStr) {
		return "", "", 0, fmt.Errorf("invalid refresh token")
	}

	params, err := c.ParseToken(refreshTokenStr)
	if err != nil {
		return "", "", 0, fmt.Errorf("invalid refresh token")
	}
	return c.GeneratedTokens(params)
}
