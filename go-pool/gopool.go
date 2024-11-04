package gopool

import (
	"github.com/alitto/pond/v2"
)

type GoPool struct {
	pool pond.Pool
}

type PoolStat struct {
	RunningWorkers  int64
	IdleWorkers     int
	SubmittedTasks  uint64
	WaitingTasks    uint64
	SuccessfulTasks uint64
	FailedTasks     uint64
	CompletedTasks  uint64
}

//
//func New(maxWorkers int, options ...pond.Option) *GoPool {
//	p := pond.NewPool(maxWorkers, options...)
//	return &GoPool{pool: p}
//}
//
//func (g *GoPool) StopAndWait() {
//	g.pool.StopAndWait()
//}
//
//// 如果池已停止并且不再接受任务，则Stopped返回true，否则返回false。
//func (g *GoPool) Stopped() bool {
//	return g.pool.Stopped()
//}
//
//// 使用固定数量的Worker创建一个无缓冲（阻塞）池，提交任务等待
//func NewFixedSizePool(maxWorkers int) *GoPool {
//	return New(maxWorkers)
//}
//
//// Create a context that will be cancelled
//// Tasks being processed will continue until they finish, but queued tasks are cancelled.
//func NewContextPool(maxWorkers int, ctx context.Context) *GoPool {
//	return New(maxWorkers)
//}
//
////	 Create a task group
////
////			group := pool.Group()
////
////			// Submit a group of related tasks
////			for i := 0; i < 20; i++ {
////				n := i
////				group.Submit(func() {
////					fmt.Printf("Running group task #%d\n", n)
////				})
////			}
////
////			// Wait for all tasks in the group to complete
////			group.Wait()
////		}
//func (g *GoPool) NewTaskGroup() pond.TaskGroup {
//	return g.pool.NewGroup()
//}
//
//func (g *GoPool) Submit(fn func()) {
//	g.pool.Submit(fn)
//}
//
//func (g *GoPool) PrometheusHandler() {
//	prometheus.MustRegister(prometheus.NewGaugeFunc(
//		prometheus.GaugeOpts{
//			Name: "pool_workers_running",
//			Help: "Number of running worker goroutines",
//		},
//		func() float64 {
//			return float64(g.pool.RunningWorkers())
//		}))
//	prometheus.MustRegister(prometheus.NewGaugeFunc(
//		prometheus.GaugeOpts{
//			Name: "pool_workers_idle",
//			Help: "Number of idle worker goroutines",
//		},
//		func() float64 {
//			return 0
//		}))
//
//	// Task metrics
//	prometheus.MustRegister(prometheus.NewCounterFunc(
//		prometheus.CounterOpts{
//			Name: "pool_tasks_submitted_total",
//			Help: "Number of tasks submitted",
//		},
//		func() float64 {
//			return float64(g.pool.SubmittedTasks())
//		}))
//	prometheus.MustRegister(prometheus.NewGaugeFunc(
//		prometheus.GaugeOpts{
//			Name: "pool_tasks_waiting_total",
//			Help: "Number of tasks waiting in the queue",
//		},
//		func() float64 {
//			return float64(g.pool.WaitingTasks())
//		}))
//	prometheus.MustRegister(prometheus.NewCounterFunc(
//		prometheus.CounterOpts{
//			Name: "pool_tasks_successful_total",
//			Help: "Number of tasks that completed successfully",
//		},
//		func() float64 {
//			return float64(g.pool.SuccessfulTasks())
//		}))
//	prometheus.MustRegister(prometheus.NewCounterFunc(
//		prometheus.CounterOpts{
//			Name: "pool_tasks_failed_total",
//			Help: "Number of tasks that completed with panic",
//		},
//		func() float64 {
//			return float64(g.pool.FailedTasks())
//		}))
//	prometheus.MustRegister(prometheus.NewCounterFunc(
//		prometheus.CounterOpts{
//			Name: "pool_tasks_completed_total",
//			Help: "Number of tasks that completed either successfully or with panic",
//		},
//		func() float64 {
//			return float64(g.pool.CompletedTasks())
//		}))
//	// Expose the registered metrics via HTTP
//	//http.Handle("/metrics", promhttp.Handler())
//}
//
//func (g *GoPool) PoolStats() PoolStat {
//	return PoolStat{
//		RunningWorkers:  g.pool.RunningWorkers(),
//		SubmittedTasks:  g.pool.SubmittedTasks(),
//		WaitingTasks:    g.pool.WaitingTasks(),
//		SuccessfulTasks: g.pool.SuccessfulTasks(),
//		FailedTasks:     g.pool.FailedTasks(),
//		CompletedTasks:  g.pool.CompletedTasks(),
//	}
//}
//
//func (g *GoPool) PrintPoolStats() {
//	ps := g.PoolStats()
//	golog.WithTag("GoPool").InfoF("RunningWorkers: %d", ps.RunningWorkers)
//	golog.WithTag("GoPool").InfoF("IdleWorkers: %d", ps.IdleWorkers)
//	golog.WithTag("GoPool").InfoF("SubmittedTasks: %d", ps.SubmittedTasks)
//	golog.WithTag("GoPool").InfoF("WaitingTasks: %d", ps.WaitingTasks)
//	golog.WithTag("GoPool").InfoF("SuccessfulTasks: %d", ps.SuccessfulTasks)
//	golog.WithTag("GoPool").InfoF("FailedTasks: %d", ps.FailedTasks)
//	golog.WithTag("GoPool").InfoF("CompletedTasks: %d", ps.CompletedTasks)
//	golog.WithTag("GoPool").InfoF("----------------------------------------------------------------")
//}
