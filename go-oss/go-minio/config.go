package gominio

type Config struct {
	Name            string `json:"name,optional" yaml:"Name,optional"`
	AccessKeyId     string `json:"accessKeyId,optional" yaml:"AccessKeyId,optional"`
	AccessKeySecret string `json:"accessKeySecret,optional" yaml:"AccessKeySecret,optional"`
	Endpoint        string `json:"endpoint,optional" yaml:"Endpoint,optional"`
	Bucket          string `json:"bucket,optional" yaml:"Bucket,optional"`
	Dir             string `json:"dir,optional" yaml:"Dir,optional"`
	Domain          string `json:"domain,optional" yaml:"Domain,optional,optional"`
	UseSSL          bool   `json:"useSSL,optional" yaml:"UseSSL,optional"`
}
