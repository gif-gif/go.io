package govault

import (
	"context"
	"fmt"

	vault "github.com/hashicorp/vault/api"
)

type Config struct {
	Name      string `yaml:"Name" json:"name,optional"`
	Address   string `yaml:"Address" json:"address,required"`
	Token     string `yaml:"Token" json:"token,required"`
	MountPath string `yaml:"MountPath" json:"mountPath,optional"`
}

type GoVault struct {
	Client *vault.Client
	Config Config
}

func New(c Config) (*GoVault, error) {
	if c.MountPath == "" {
		c.MountPath = "secret"
	}
	config := vault.DefaultConfig()
	config.Address = c.Address

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, err
	}

	client.SetToken(c.Token)

	return &GoVault{Client: client, Config: c}, nil
}

func NewVault(config *vault.Config) (*GoVault, error) {
	client, err := vault.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &GoVault{Client: client, Config: Config{Address: config.Address}}, nil
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

// 历史元数据数据列表
func (vm *GoVault) GetVersionsAsListV2(path string) ([]vault.KVVersionMetadata, error) {
	secret, err := vm.Client.KVv2(vm.Config.MountPath).GetVersionsAsList(context.Background(), path)
	if err != nil {
		return nil, fmt.Errorf("GetVersionsAsListV2 读取密钥失败: %w", err)
	}
	// 提取数据
	return secret, nil
}

// V2 写入时，path 为 ${MountPath}/data/${path}
func (vm *GoVault) WriteV2(path string, data map[string]interface{}) error {
	// 写入密钥到 KV v2
	_, err := vm.Client.KVv2(vm.Config.MountPath).Put(context.Background(), path, data)
	if err != nil {
		return fmt.Errorf("写入密钥失败: %w", err)
	}
	return nil
}

func (vm *GoVault) ReadV2(path string) (map[string]interface{}, error) {
	// 读取密钥
	secret, err := vm.Client.KVv2(vm.Config.MountPath).Get(context.Background(), path)
	if err != nil {
		return nil, fmt.Errorf("读取密钥失败: %w", err)
	}
	// 提取数据
	data := secret.Data
	return data, nil
}

func (vm *GoVault) ReadByVersionV2(path string, version int) (map[string]interface{}, error) {
	// 读取密钥
	secret, err := vm.Client.KVv2(vm.Config.MountPath).GetVersion(context.Background(), path, version)
	if err != nil {
		return nil, fmt.Errorf("读取密钥失败: %w", err)
	}
	// 提取数据
	data := secret.Data
	return data, nil
}

// 软删除
func (vm *GoVault) DeleteV2(path string) error {
	// 删除密钥
	err := vm.Client.KVv2(vm.Config.MountPath).Delete(context.Background(), path)
	if err != nil {
		return fmt.Errorf("删除密钥失败: %w", err)
	}

	return nil
}

func (vm *GoVault) DestroyV2(path string, versions []int) error {
	// 删除密钥
	err := vm.Client.KVv2(vm.Config.MountPath).Destroy(context.Background(), path, versions)
	if err != nil {
		return fmt.Errorf("销毁密钥失败: %w", err)
	}

	return nil
}

func (vm *GoVault) UndeleteV2(path string, versions []int) error {
	// 删除密钥
	err := vm.Client.KVv2(vm.Config.MountPath).Undelete(context.Background(), path, versions)
	if err != nil {
		return fmt.Errorf("撤销删除密钥失败: %w", err)
	}

	return nil
}

// 常用用户名和密码 方法----------------------------- kafka redis mysql starrocks minio elk etcd n8n postgres

func (vm *GoVault) CommonUsernamePasswordWrite(path string, data UserNameAndPassword) error {
	return vm.WriteV2(path, map[string]interface{}{
		"username": data.Username,
		"password": data.Password,
	})
}

func (vm *GoVault) CommonUsernamePasswordRead(path string) (*UserNameAndPassword, error) {
	data, err := vm.ReadV2(path)
	if err != nil {
		return nil, fmt.Errorf("读取常用用户名和密码失败: %w", err)
	}
	// 提取数据
	username := data["username"].(string)
	password := data["password"].(string)
	return &UserNameAndPassword{Username: username, Password: password}, nil
}

// 常用Mysql 用户和密码 方法-----------------------------

func (vm *GoVault) CommonWriteMysql(path string, mysqlDataSource MysqlDataSource) error {
	return vm.WriteV2(path, ParseMap(mysqlDataSource))
}

func (vm *GoVault) CommonReadMysql(path string) (*MysqlDataSource, error) {
	data, err := vm.ReadV2(path)
	if err != nil {
		return nil, fmt.Errorf("读取常用Mysql用户名和密码失败: %w", err)
	}
	// 提取数据
	mysqlDataSource := ParseData(data)
	return &mysqlDataSource, nil
}

func (vm *GoVault) CommonSecretWrite(path string, secret string) error {
	return vm.WriteV2(path, map[string]interface{}{
		"secretKey": secret,
	})
}

func (vm *GoVault) CommonSecretRead(path string) (string, error) {
	data, err := vm.ReadV2(path)
	if err != nil {
		return "", fmt.Errorf("读取常用Secret失败: %w", err)
	}
	// 提取数据
	secret := data["secretKey"].(string)
	return secret, nil
}
