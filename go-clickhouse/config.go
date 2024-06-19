package goclickhouse

type Config struct {
	Driver       string `yaml:"driver" json:"driver"`
	Addr         string `yaml:"addr" json:"addr"`
	User         string `yaml:"user" json:"user"`
	Password     string `yaml:"password" json:"password"`
	Database     string `yaml:"database" json:"database"`
	ReadTimeout  int32  `yaml:"read_timeout" json:"read_timeout"`
	WriteTimeout int32  `yaml:"write_timeout" json:"write_timeout"`
	AltHosts     string `yaml:"alt_hosts" json:"alt_hosts"`
	AutoPing     bool   `yaml:"auto_ping" json:"auto_ping"`
	Debug        bool   `yaml:"debug" json:"debug"`
}
