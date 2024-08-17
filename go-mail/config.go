package gomail

type Config struct {
	Username string `json:"username,optional"  yaml:"Username"`
	Password string `json:"password,optional"  yaml:"Password"`
	Host     string `json:"host,optional"  yaml:"Host"`
	Port     int    `json:"port,optional"  yaml:"Port"`
	TLS      bool   `json:"tls,optional"  yaml:"TLS"`
}
