package goetcd

type Config struct {
	Endpoints []string `json:"endpoints,optional" yaml:"Endpoints,optional"`
	Name      string   `json:"name,optional" yaml:"Name,optional"`
	Username  string   `json:"username,optional" yaml:"Username,optional"`
	Password  string   `json:"password,optional" yaml:"Password,optional"`
	TLS       *TLS     `json:"tls,optional" yaml:"TLS,optional"`
}

type TLS struct {
	CertFile string `json:"certFile,optional" yaml:"CertFile,optional"`
	KeyFile  string `json:"keyFile,optional" yaml:"KeyFile,optional"`
	CAFile   string `json:"caFile,optional" yaml:"CAFile,optional"`
}
