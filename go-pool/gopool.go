package gopool

import (
	"context"
	"github.com/alitto/pond"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/prometheus/client_golang/prometheus"
)

type GoPool struct {
	pool *pond.WorkerPool
}

type PoolStat struct {
	RunningWorkers  int
	IdleWorkers     int
	SubmittedTasks  uint64
	WaitingTasks    uint64
	SuccessfulTasks uint64
	FailedTasks     uint64
	CompletedTasks  uint64
}

func New(maxWorkers, maxCapacity int, options ...pond.Option) *GoPool {
	p := pond.New(maxWorkers, maxCapacity, options...)
	return &GoPool{pool: p}
}

func (g *GoPool) StopAndWait() {
	g.pool.StopAndWait()
}

// Create a buffered (non-blocking) pool that can scale up to maxWorkers workers
//
// and has a buffer capacity of maxCapacity tasks
//
// 创建一个缓冲（非阻塞）池，最多可扩展到maxWorkers个Worker，缓冲容量为maxCapacity个任务(大于这个会阻塞等待提交)
func NewDynamicSizePool(maxWorkers, maxCapacity int) *GoPool {
	return New(maxWorkers, maxCapacity)
}

// 使用固定数量的Worker创建一个无缓冲（阻塞）池，提交任务等待
func NewFixedSizePool(maxWorkers, minWorkers int) *GoPool {
	return New(maxWorkers, 0, pond.MinWorkers(minWorkers))
}

// Create a context that will be cancelled
// Tasks being processed will continue until they finish, but queued tasks are cancelled.
func NewContextPool(maxWorkers, maxCapacity int, ctx context.Context) *GoPool {
	//ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	//defer stop()
	// Create a pool and pass the context to it.
	return New(maxWorkers, maxCapacity, pond.Context(ctx))
}

//	 Create a task group
//
//			group := pool.Group()
//
//			// Submit a group of related tasks
//			for i := 0; i < 20; i++ {
//				n := i
//				group.Submit(func() {
//					fmt.Printf("Running group task #%d\n", n)
//				})
//			}
//
//			// Wait for all tasks in the group to complete
//			group.Wait()
//		}
func (g *GoPool) NewTaskGroup() *pond.TaskGroup {
	return g.pool.Group()
}

// group, ctx := pool.GroupContext(context.Background())
//
//	group.Submit(func() error {
//		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
//			resp, err := http.DefaultClient.Do(req)
//			if err == nil {
//				resp.Body.Close()
//			}
//		return err
//	})
//
// Wait for all fn to complete.
//
//	err := group.Wait()
//	if err != nil {
//		fmt.Printf("Failed to Error: %v", err)
//	} else {
//		fmt.Println("Successfully all")
//	}
//
// Create a task group associated to a context
func (g *GoPool) NewGroupContext() (*pond.TaskGroupWithContext, context.Context) {
	return g.pool.GroupContext(context.Background())
}

func (g *GoPool) Submit(fn func()) {
	g.pool.Submit(fn)
}

func (g *GoPool) PrometheusHandler() {
	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "pool_workers_running",
			Help: "Number of running worker goroutines",
		},
		func() float64 {
			return float64(g.pool.RunningWorkers())
		}))
	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "pool_workers_idle",
			Help: "Number of idle worker goroutines",
		},
		func() float64 {
			return float64(g.pool.IdleWorkers())
		}))

	// Task metrics
	prometheus.MustRegister(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "pool_tasks_submitted_total",
			Help: "Number of tasks submitted",
		},
		func() float64 {
			return float64(g.pool.SubmittedTasks())
		}))
	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "pool_tasks_waiting_total",
			Help: "Number of tasks waiting in the queue",
		},
		func() float64 {
			return float64(g.pool.WaitingTasks())
		}))
	prometheus.MustRegister(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "pool_tasks_successful_total",
			Help: "Number of tasks that completed successfully",
		},
		func() float64 {
			return float64(g.pool.SuccessfulTasks())
		}))
	prometheus.MustRegister(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "pool_tasks_failed_total",
			Help: "Number of tasks that completed with panic",
		},
		func() float64 {
			return float64(g.pool.FailedTasks())
		}))
	prometheus.MustRegister(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "pool_tasks_completed_total",
			Help: "Number of tasks that completed either successfully or with panic",
		},
		func() float64 {
			return float64(g.pool.CompletedTasks())
		}))
	// Expose the registered metrics via HTTP
	//http.Handle("/metrics", promhttp.Handler())
}

func (g *GoPool) PoolStats() PoolStat {
	return PoolStat{
		RunningWorkers:  g.pool.RunningWorkers(),
		IdleWorkers:     g.pool.IdleWorkers(),
		SubmittedTasks:  g.pool.SubmittedTasks(),
		WaitingTasks:    g.pool.WaitingTasks(),
		SuccessfulTasks: g.pool.SuccessfulTasks(),
		FailedTasks:     g.pool.FailedTasks(),
		CompletedTasks:  g.pool.CompletedTasks(),
	}
}

func (g *GoPool) PrintPoolStats() {
	ps := g.PoolStats()
	golog.WithTag("GoPool").InfoF("RunningWorkers: %d", ps.RunningWorkers)
	golog.WithTag("GoPool").InfoF("IdleWorkers: %d", ps.IdleWorkers)
	golog.WithTag("GoPool").InfoF("SubmittedTasks: %d", ps.SubmittedTasks)
	golog.WithTag("GoPool").InfoF("WaitingTasks: %d", ps.WaitingTasks)
	golog.WithTag("GoPool").InfoF("SuccessfulTasks: %d", ps.SuccessfulTasks)
	golog.WithTag("GoPool").InfoF("FailedTasks: %d", ps.FailedTasks)
	golog.WithTag("GoPool").InfoF("CompletedTasks: %d", ps.CompletedTasks)
	golog.WithTag("GoPool").InfoF("----------------------------------------------------------------")
}
