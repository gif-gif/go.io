package gocache

type Config struct {
	Name              string
	DefaultExpiration int64 //秒
	CleanupInterval   int64 //秒
}
