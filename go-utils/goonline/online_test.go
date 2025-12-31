package goonline

import (
	"log"
	"testing"
	"time"

	goetcd "github.com/gif-gif/go.io/go-etcd"
)

func TestOnlineManager(t *testing.T) {
	etcdEndpoints := []string{"localhost:2379"}
	usersEtcdConfigs := goetcd.Config{
		Endpoints: etcdEndpoints,
	}

	goetcd.Init(usersEtcdConfigs)
	client := goetcd.Default().Client
	userConfigs := &Config{
		EntityType: "user",
	}

	// 1. 创建用户在线管理器
	userManager, err := New(client, userConfigs)
	if err != nil {
		log.Fatalf("创建用户管理器失败: %v", err)
	}
	defer userManager.Close()
	// 2. 创建服务器在线管理器
	serverConfigs := &Config{
		EntityType: "server",
	}
	serverManager, err := New(client, serverConfigs)
	if err != nil {
		log.Fatalf("创建服务器管理器失败: %v", err)
	}
	defer serverManager.Close()
	// 3. 创建设备在线管理器
	deviceConfigs := &Config{
		EntityType: "device",
	}
	deviceManager, err := New(client, deviceConfigs)
	if err != nil {
		log.Fatalf("创建设备管理器失败: %v", err)
	}
	defer deviceManager.Close()

	// 4. 创建多个监听器
	userWatcher, err := NewOnlineWatcher(client, userConfigs)
	if err != nil {
		log.Fatalf("创建用户监听器失败: %v", err)
	}
	defer userWatcher.Close()
	go userWatcher.Watch()

	serverWatcher, _ := NewOnlineWatcher(client, serverConfigs)
	defer serverWatcher.Close()
	go serverWatcher.Watch()

	deviceWatcher, _ := NewOnlineWatcher(client, deviceConfigs)
	defer deviceWatcher.Close()
	go deviceWatcher.Watch()

	time.Sleep(1 * time.Second)

	// 5. 用户上线
	log.Println("\n=== 用户上线 ===")
	userManager.SetOnline("user001", 0, map[string]interface{}{
		"name": "张三",
		"ip":   "192.168.1.100",
		"role": "admin",
	})
	userManager.SetOnline("user002", 0, map[string]interface{}{
		"name": "李四",
		"ip":   "192.168.1.101",
		"role": "user",
	})

	time.Sleep(1 * time.Second)

	// 6. 服务器上线
	log.Println("\n=== 服务器上线 ===")
	serverManager.SetOnline("server001", 0, map[string]interface{}{
		"hostname": "web-01",
		"ip":       "10.0.1.10",
		"region":   "us-east-1",
		"cpu":      "16核",
		"memory":   "32GB",
	})
	serverManager.SetOnline("server002", 0, map[string]interface{}{
		"hostname": "db-01",
		"ip":       "10.0.1.20",
		"region":   "us-east-1",
		"cpu":      "32核",
		"memory":   "64GB",
	})

	time.Sleep(1 * time.Second)

	// 7. 设备上线
	log.Println("\n=== 设备上线 ===")
	deviceManager.SetOnline("device001", 0, map[string]interface{}{
		"type":     "camera",
		"location": "门口",
		"status":   "正常",
	})
	deviceManager.SetOnline("device002", 0, map[string]interface{}{
		"type":     "sensor",
		"location": "仓库",
		"temp":     "25℃",
	})

	time.Sleep(2 * time.Second)

	// 8. 查询各类型在线数量
	log.Println("\n=== 在线统计 ===")
	userCount, _ := userManager.GetOnlineCount()
	serverCount, _ := serverManager.GetOnlineCount()
	deviceCount, _ := deviceManager.GetOnlineCount()
	log.Printf("用户在线: %d", userCount)
	log.Printf("服务器在线: %d", serverCount)
	log.Printf("设备在线: %d", deviceCount)

	// 9. 模拟续租
	log.Println("\n=== 模拟续租 ===")
	time.Sleep(3 * time.Second)

	// 用户续租
	userManager.SetOnline("user001", 0, map[string]interface{}{
		"name": "张三",
		"ip":   "192.168.1.100",
		"role": "super_admin", // 权限升级
	})

	// 服务器续租并更新状态
	serverManager.SetOnline("server001", 0, map[string]interface{}{
		"hostname": "web-01",
		"ip":       "10.0.1.10",
		"region":   "us-east-1",
		"cpu":      "16核",
		"memory":   "32GB",
		"load":     "0.8", // 新增负载信息
	})

	time.Sleep(2 * time.Second)

	// 10. 查看详细信息
	log.Println("\n=== 查看详细信息 ===")
	user, _ := userManager.Get("user001")
	if user != nil {
		log.Printf("用户 user001: 上线时间=%s, 更新时间=%s, 数据=%v",
			user.OnlineAt, user.UpdateAt, user.Data)
	}

	server, _ := serverManager.Get("server001")
	if server != nil {
		log.Printf("服务器 server001: %v", server.Data)
	}

	// 11. 下线操作
	log.Println("\n=== 下线操作 ===")
	time.Sleep(2 * time.Second)
	userManager.SetOffline("user002")
	deviceManager.SetOffline("device001")

	// 12. 最终统计
	time.Sleep(2 * time.Second)
	log.Println("\n=== 最终在线统计 ===")

	users, _ := userManager.GetOnlineList()
	log.Printf("\n在线用户 (%d):", len(users))
	for _, u := range users {
		log.Printf("  - %s: %v", u.ID, u.Data)
	}

	servers, _ := serverManager.GetOnlineList()
	log.Printf("\n在线服务器 (%d):", len(servers))
	for _, s := range servers {
		log.Printf("  - %s: %v", s.ID, s.Data)
	}

	devices, _ := deviceManager.GetOnlineList()
	log.Printf("\n在线设备 (%d):", len(devices))
	for _, d := range devices {
		log.Printf("  - %s: %v", d.ID, d.Data)
	}

	log.Println("\n程序运行10秒后退出...")
	time.Sleep(10 * time.Second)
}
