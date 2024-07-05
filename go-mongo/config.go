package gomongo

type Config struct {
	Name           string `yaml:"Name" json:"name"`
	Addr           string `yaml:"Addr" json:"addr"`
	User           string `yaml:"User" json:"user"`
	Password       string `yaml:"Password" json:"password"`
	EnablePassword bool   `yaml:"EnablePassword" json:"enable_password"`
	Database       string `yaml:"Database" json:"database"`
	AutoPing       bool   `yaml:"AutoPing" json:"autoPing"`
}
