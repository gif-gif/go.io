package goo_es

type Config struct {
	Addr      string `json:"addr" yaml:"addr"`
	User      string `json:"user" yaml:"user"`
	Password  string `json:"password" yaml:"password"`
	EnableLog bool   `json:"enable_log" yaml:"enable_log"`
}
