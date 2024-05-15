package gominio

type Config struct {
	AccessKeyId     string `json:"access_key_id" yaml:"AccessKeyId"`
	AccessKeySecret string `json:"access_key_secret" yaml:"AccessKeySecret"`
	Endpoint        string `json:"endpoint" yaml:"Endpoint"`
	Bucket          string `json:"bucket" yaml:"Bucket"`
	Domain          string `json:"domain" yaml:"Domain,optional"`
	UseSSL          bool   `json:"useSSL" yaml:"UseSSL"`
}
