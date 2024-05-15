package gominio

type Config struct {
	AccessKeyId     string `json:"access_key_id" yaml:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret" yaml:"access_key_secret"`
	Endpoint        string `json:"endpoint" yaml:"endpoint"`
	Bucket          string `json:"bucket" yaml:"bucket"`
	Domain          string `json:"domain" yaml:"domain,optional"`
	UseSSL          bool   `json:"useSSL" yaml:"useSSL"`
}
