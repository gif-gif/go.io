package goprometheus

type Config struct {
	Name    string   `json:"name,optional" yaml:"Name"`
	Address string   `json:"address,optional" yaml:"Address" default:"0.0.0.0:9090"`
	Filters []string `json:"filters,optional" yaml:"Filters"`
}

func (c *Config) GetAddress() string {
	if c.Address == "" {
		c.Address = "0.0.0.0:9090"
	}
	return c.Address
}
