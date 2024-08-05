package goes

type Config struct {
	Addr      string `json:"addr,optional" yaml:"Addr,optional"`
	User      string `json:"user,optional" yaml:"User,optional"`
	Password  string `json:"password,optional" yaml:"Password,optional"`
	EnableLog bool   `json:"enableLog,optional" yaml:"EnableLog,optional"`
}
