package goetcd

type Config struct {
	Endpoints []string `json:"endpoints,optional"  yaml:"Endpoints"`
	Name      string   `json:"name,optional"  yaml:"Name"`
	Username  string   `json:"username,optional"  yaml:"Username"`
	Password  string   `json:"password,optional"  yaml:"Password"`
	TLS       *TLS     `json:"tls,optional"  yaml:"TLS"`
}

type TLS struct {
	CertFile string `json:"certFile,optional"  yaml:"CertFile"`
	KeyFile  string `json:"keyFile,optional"  yaml:"KeyFile"`
	CAFile   string `json:"caFile,optional"  yaml:"CAFile"`
}
