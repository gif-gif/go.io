# 通用进程生命周期管理
```go
NewRunner 创建一个新的 Runner
示例:
	// userService
	userService := process.NewRunner(process.Options{
	    Name:     "userService",
	    ExecPath: "./userService",
	    Args:     []string{"run", "-c", "./cfg.json"},
	})

	// apiService
	apiService := process.NewRunner(process.Options{
	    Name:     "apiService",
	    ExecPath: "./apiService",
	    Args:     []string{"run", "-c", "./config.json"},
	})
```