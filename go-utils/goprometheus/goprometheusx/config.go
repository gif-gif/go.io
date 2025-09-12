package goprometheusx

type Config struct {
	Host string `json:"host,optional" yaml:"host"`
	Port int    `json:"port,default=9101" yaml:"port"`
	Path string `json:"path,default=/metrics" yaml:"path"`
}
