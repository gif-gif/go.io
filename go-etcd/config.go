package goetcd

type Config struct {
	Endpoints []string `json:"endpoints" yaml:"Endpoints"`
	Name      string   `json:"name" yaml:"Name"`
	Username  string   `json:"username" yaml:"Username"`
	Password  string   `json:"password" yaml:"Password"`
	TLS       *TLS     `json:"tls" yaml:"TLS"`
}

type TLS struct {
	CertFile string `json:"cert_file" yaml:"CertFile"`
	KeyFile  string `json:"key_file" yaml:"KeyFile"`
	CAFile   string `json:"ca_file" yaml:"CAFile"`
}
