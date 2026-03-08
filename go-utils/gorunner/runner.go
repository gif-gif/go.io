package gorunner

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// Options 启动选项
type Options struct {
	Name              string        // 进程名称，用于日志标识
	ExecPath          string        // 可执行文件路径
	Args              []string      // 启动参数
	PidFile           string        // PID 文件路径，为空则不使用，例如 ./xray.pid
	RestartDelay      time.Duration // 异常退出后重启延迟，默认 2s
	StartRetryDelay   time.Duration // 启动失败重试延迟，默认 5s
	KillCheckInterval time.Duration // kill 检查间隔，默认 300ms
}

func (o *Options) setDefaults() {
	if o.RestartDelay == 0 {
		o.RestartDelay = 2 * time.Second
	}
	if o.StartRetryDelay == 0 {
		o.StartRetryDelay = 5 * time.Second
	}
	if o.KillCheckInterval == 0 {
		o.KillCheckInterval = 300 * time.Millisecond
	}
	if o.Name == "" {
		o.Name = o.ExecPath
	}
}

// ---- PID 文件工具 ----

func writePidFile(path string, pid int) error {
	if path == "" {
		return nil
	}
	return os.WriteFile(path, []byte(strconv.Itoa(pid)), 0644)
}

func readPidFile(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, fmt.Errorf("invalid pid file content: %w", err)
	}
	return pid, nil
}

func removePidFile(path string) {
	if path != "" {
		os.Remove(path)
	}
}

// findLivingProcess 通过 PID 查找存活进程，返回 nil 表示进程不存在
func findLivingProcess(pid int) *os.Process {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return nil
	}
	// Linux/macOS 下 FindProcess 总是成功，需要发送 signal 0 验证进程是否存活
	if err := proc.Signal(syscall.Signal(0)); err != nil {
		return nil
	}
	return proc
}

// ---- instance ----

type instance struct {
	cmd        *exec.Cmd
	ctx        context.Context
	cancel     context.CancelFunc
	isKilled   bool
	isExited   bool
	mutex      sync.Mutex
	opts       Options
	instanceId uint32
}

func newInstance(opts Options, id uint32) *instance {
	ctx, cancel := context.WithCancel(context.Background())
	return &instance{
		opts:       opts,
		instanceId: id,
		ctx:        ctx,
		cancel:     cancel,
	}
}

// run 统一的启动循环：
//   - attachProc != nil 表示先接管已有进程，守护它，死了再走正常启动流程
//   - attachProc == nil 直接走正常启动流程
func (ins *instance) run(attachProc *os.Process) {
	log.Infof("[%s] run instance: %v current PID:%d", ins.opts.Name, ins.instanceId, os.Getpid())

	go func() {
		defer func() {
			ins.mutex.Lock()
			defer ins.mutex.Unlock()
			ins.isExited = true
			removePidFile(ins.opts.PidFile)
		}()

		// 如果有接管的进程，先守护它
		if attachProc != nil {
			log.Infof("[%s] attached to previous process PID:%d instance: %v",
				ins.opts.Name, attachProc.Pid, ins.instanceId)

			// 设置伪 cmd，供 kill() 使用
			ins.cmd = &exec.Cmd{Process: attachProc}

		watchLoop:
			for {
				select {
				case <-ins.ctx.Done():
					log.Infof("[%s] attached watcher exiting instance: %v", ins.opts.Name, ins.instanceId)
					return
				case <-time.After(ins.opts.KillCheckInterval):
				}

				ins.mutex.Lock()
				killed := ins.isKilled
				ins.mutex.Unlock()
				if killed {
					return
				}

				if err := attachProc.Signal(syscall.Signal(0)); err != nil {
					log.Infof("[%s] attached process PID:%d exited: %v, starting new...", ins.opts.Name, attachProc.Pid, err)
					removePidFile(ins.opts.PidFile)
					ins.mutex.Lock()
					ins.cmd = nil
					ins.mutex.Unlock()
					break watchLoop // 跳出守护，进入正常启动流程
				}
			}
		}

		// 正常启动循环
		for {
			select {
			case <-ins.ctx.Done():
				log.Infof("[%s] instance exiting: %v", ins.opts.Name, ins.instanceId)
				return
			default:
			}

			log.Infof("[%s] running instance: %v cmd: %s %v", ins.opts.Name, ins.instanceId, ins.opts.ExecPath, ins.opts.Args)

			ins.cmd = exec.CommandContext(ins.ctx, ins.opts.ExecPath, ins.opts.Args...)
			ins.cmd.Stdout = os.Stdout
			ins.cmd.Stderr = os.Stderr

			if err := ins.cmd.Start(); err != nil {
				log.Warnf("[%s] failed to start instance: %v error: %v",
					ins.opts.Name, ins.instanceId, err)
				time.Sleep(ins.opts.StartRetryDelay)
				continue
			}

			pid := ins.cmd.Process.Pid
			log.Infof("[%s] started PID:%d instance: %v", ins.opts.Name, pid, ins.instanceId)

			if err := writePidFile(ins.opts.PidFile, pid); err != nil {
				log.Warnf("[%s] write pid file failed: %v", ins.opts.Name, err)
			}

			waitErr := ins.cmd.Wait()

			ins.mutex.Lock()
			killed := ins.isKilled
			ins.mutex.Unlock()

			if killed {
				log.Infof("[%s] killed instance: %v", ins.opts.Name, ins.instanceId)
				return
			}

			if ins.ctx.Err() != nil {
				log.Infof("[%s] shutting down instance: %v", ins.opts.Name, ins.instanceId)
				return
			}

			if waitErr != nil {
				log.Warnf("[%s] exited with error: %v instance: %v, restarting...", ins.opts.Name, waitErr, ins.instanceId)
			} else {
				log.Infof("[%s] exited normally instance: %v, restarting...", ins.opts.Name, ins.instanceId)
			}

			time.Sleep(ins.opts.RestartDelay)
		}
	}()
}

func (ins *instance) kill() {
	log.Infof("[%s] killing instance: %v", ins.opts.Name, ins.instanceId)

	ins.mutex.Lock()
	ins.isKilled = true
	ins.cancel()
	proc := ins.cmd.Process
	ins.mutex.Unlock()

	if proc == nil {
		return
	}

	go func() {
		for {
			if err := proc.Kill(); err != nil {
				log.Warnf("[%s] kill failed: %v instance: %v",
					ins.opts.Name, err, ins.instanceId)
				return
			}
			time.Sleep(ins.opts.KillCheckInterval)

			ins.mutex.Lock()
			exited := ins.isExited
			ins.mutex.Unlock()
			if exited {
				return
			}
		}
	}()
}

// ---- Runner ----

type Runner struct {
	mutex       sync.Mutex
	opts        Options
	instance    *instance
	instanceIdx uint32
}

func NewRunner(opts Options) *Runner {
	opts.setDefaults()
	return &Runner{
		opts:        opts,
		instanceIdx: rand.Uint32(),
	}
}

// UpdateOptions 动态更新配置（下次 Run 时生效）
func (r *Runner) UpdateOptions(fn func(opts *Options)) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	fn(&r.opts)
	r.opts.setDefaults()
}

// Run 统一入口：自动检测 PID 文件，有存活进程则接管，否则直接启动
// 若已有实例在运行，先 kill 再执行
func (r *Runner) Run() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.instance != nil {
		r.instance.kill()
	}

	// 尝试从 PID 文件找到遗留进程
	var attachProc *os.Process
	if r.opts.PidFile != "" {
		if pid, err := readPidFile(r.opts.PidFile); err == nil {
			if proc := findLivingProcess(pid); proc != nil {
				attachProc = proc
			} else {
				log.Infof("[%s] previous process PID:%d not found, starting new", r.opts.Name, pid)
				removePidFile(r.opts.PidFile)
			}
		}
	}

	ins := newInstance(r.opts, r.instanceIdx)
	r.instanceIdx++
	r.instance = ins
	ins.run(attachProc)
}

// Kill 终止当前进程
func (r *Runner) Kill() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.instance == nil {
		return
	}
	r.instance.kill()
	r.instance = nil
}

// IsRunning 返回当前是否有运行中的实例
func (r *Runner) IsRunning() bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.instance == nil {
		return false
	}
	r.instance.mutex.Lock()
	defer r.instance.mutex.Unlock()
	return !r.instance.isExited
}
