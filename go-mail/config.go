package gomail

type Config struct {
	Username string `json:"username,optional" yaml:"Username,optional"`
	Password string `json:"password,optional" yaml:"Password,optional"`
	Host     string `json:"host,optional" yaml:"Host,optional"`
	Port     int    `json:"port,optional" yaml:"Port,optional"`
	TLS      bool   `json:"tls,optional" yaml:"TLS,optional"`
}
