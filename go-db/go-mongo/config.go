package gomongo

type Config struct {
	Name           string `yaml:"Name,optional" json:"name,optional"`
	Addr           string `yaml:"Addr,optional" json:"addr,optional"`
	User           string `yaml:"User,optional" json:"user,optional"`
	Password       string `yaml:"Password,optional" json:"password,optional"`
	EnablePassword bool   `yaml:"EnablePassword,optional" json:"enablePassword,optional"`
	Database       string `yaml:"Database,optional" json:"database,optional"`
	AutoPing       bool   `yaml:"AutoPing,optional" json:"autoPing,optional"`
}
