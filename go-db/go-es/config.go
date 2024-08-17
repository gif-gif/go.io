package goes

type Config struct {
	Addr      string `json:"addr,optional"  yaml:"Addr"`
	User      string `json:"user,optional"  yaml:"User"`
	Password  string `json:"password,optional"  yaml:"Password"`
	EnableLog bool   `json:"enableLog,optional"  yaml:"EnableLog"`
}
