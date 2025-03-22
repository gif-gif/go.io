package gooauth

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
	"io"
	"net/http"
)

func VerifyGoogleToken(accessToken string, clientId string, ctx context.Context) (*GoogleUserInfo, error) {
	// 创建验证器
	validator, err := idtoken.NewValidator(ctx, option.WithoutAuthentication())
	if err != nil {
		return nil, fmt.Errorf("failed to create validator: %v", err)
	}

	// 验证 ID Token
	payload, err := validator.Validate(ctx, accessToken, clientId)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	// 从 payload 创建 GoogleUser 结构体
	googleUser := GoogleUserInfo{
		Email:         payload.Claims["email"].(string),
		VerifiedEmail: payload.Claims["email_verified"].(bool),
		Name:          payload.Claims["name"].(string),
		GivenName:     payload.Claims["given_name"].(string),
		FamilyName:    payload.Claims["family_name"].(string),
		Picture:       payload.Claims["picture"].(string),
		Locale:        "",
	}

	return &googleUser, nil
}

func GetUserInfoFromGoogle(accessToken string) (*GoogleUserInfo, error) {
	// 创建HTTP请求获取用户信息
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API响应错误: %s", string(data))
	}

	// 解析用户信息
	var userInfo GoogleUserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %v", err)
	}

	return &userInfo, nil
}
