package goes

type Config struct {
	Addrs     []string `json:"addrs" yaml:"addrs"`
	User      string   `json:"user" yaml:"user"`
	Password  string   `json:"password" yaml:"password"`
	APIKey    string   `json:"apiKey" yaml:"apiKey"`
	EnableLog bool     `json:"enableLog" yaml:"enableLog"`
}
