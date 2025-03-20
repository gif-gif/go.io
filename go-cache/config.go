package gocache

type Config struct {
	Name              string `yaml:"Name" json:"name,optional"`
	DefaultExpiration int64  `yaml:"DefaultExpiration" json:"defaultExpiration,optional"` //秒
	CleanupInterval   int64  `yaml:"CleanupInterval" json:"cleanupInterval,optional"`     //秒
}
