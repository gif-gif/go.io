package gomongo

type Config struct {
	Name           string `yaml:"Name" json:"name,optional"`
	Addr           string `yaml:"Addr" json:"addr,optional"`
	User           string `yaml:"User" json:"user,optional"`
	Password       string `yaml:"Password" json:"password,optional"`
	EnablePassword bool   `yaml:"EnablePassword" json:"enablePassword,optional"`
	Database       string `yaml:"Database" json:"database,optional"`
	AutoPing       bool   `yaml:"AutoPing" json:"autoPing,optional"`
}
