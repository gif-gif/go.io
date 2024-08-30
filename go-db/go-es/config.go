package goes

type Config struct {
	Name      string `json:"name,optional"  yaml:"Name"`
	Addr      string `json:"addr,optional"  yaml:"Addr"`
	User      string `json:"user,optional"  yaml:"User"`
	Password  string `json:"password,optional"  yaml:"Password"`
	EnableLog bool   `json:"enableLog,optional"  yaml:"EnableLog"`
}
