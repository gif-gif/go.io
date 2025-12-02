package govault

import (
	"context"
	"fmt"

	vault "github.com/hashicorp/vault/api"
)

type Config struct {
	Name    string `yaml:"Name" json:"name,optional"`
	Address string `yaml:"Address" json:"address,required"`
	Token   string `yaml:"Token" json:"token,required"`
}

type GoVault struct {
	Client *vault.Client
	Config Config
}

func New(c Config) (*GoVault, error) {
	config := vault.DefaultConfig()
	config.Address = c.Address

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, err
	}

	client.SetToken(c.Token)

	return &GoVault{Client: client, Config: c}, nil
}

// 健康检查
func (vm *GoVault) HealthCheck() error {
	health, err := vm.Client.Sys().Health()
	if err != nil {
		return err
	}

	if !health.Initialized {
		return fmt.Errorf("GoVault 未初始化")
	}

	if health.Sealed {
		return fmt.Errorf("GoVault 已封存")
	}

	return nil
}

//	data := map[string]interface{}{
//		"username": "admin",
//		"password": "secret123",
//	}
func (vm *GoVault) WriteSecretV1(path string, data map[string]interface{}) error {
	_, err := vm.Client.Logical().Write(path, data)
	return err
}

// KV v1 使用 Logical API
//
//	secret, err := Client.Logical().Read("secret/myapp/config")
//	if err != nil {
//	    return err
//	}
//
//	if secret == nil {
//	    return fmt.Errorf("密钥不存在")
//	}
//
//	username := secret.Data["username"].(string)
//	password := secret.Data["password"].(string)
//
//	fmt.Printf("Username: %s, Password: %s\n", username, password)
//	return nil
func (vm *GoVault) ReadSecretV1(path string) (map[string]interface{}, error) {
	// KV v1 使用 Logical API
	secret, err := vm.Client.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("读取 secret 失败: %w", err)
	}

	if secret == nil {
		return nil, fmt.Errorf("密钥不存在")
	}
	return secret.Data, nil
}

func (vm *GoVault) DeleteSecretV1(path string) error {
	// 删除密钥
	_, err := vm.Client.Logical().Delete(path)
	if err != nil {
		return fmt.Errorf("删除密钥失败: %w", err)
	}

	return nil
}

func (vm *GoVault) WriteSecretV2(path string, data map[string]interface{}) error {
	// 写入密钥到 KV v2
	_, err := vm.Client.KVv2("secret").Put(context.Background(), path, data)
	if err != nil {
		return fmt.Errorf("写入密钥失败: %w", err)
	}
	return nil
}

func (vm *GoVault) ReadSecretV2(path string) (map[string]interface{}, error) {
	// 读取密钥
	secret, err := vm.Client.KVv2("secret").Get(context.Background(), path)
	if err != nil {
		return nil, fmt.Errorf("读取密钥失败: %w", err)
	}
	// 提取数据
	data := secret.Data
	return data, nil
}

func (vm *GoVault) DeleteSecretV2(path string) error {
	// 删除密钥
	err := vm.Client.KVv2("secret").Delete(context.Background(), path)
	if err != nil {
		return fmt.Errorf("删除密钥失败: %w", err)
	}

	return nil
}
