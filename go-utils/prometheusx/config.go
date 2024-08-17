package prometheusx

type Config struct {
	Host string `json:"host,optional" yaml:"host,omitempty"`
	Port int    `json:"port,default=9101" yaml:"port,omitempty"`
	Path string `json:"path,default=/metrics" yaml:"path"`
}
