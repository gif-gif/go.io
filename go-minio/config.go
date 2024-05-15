package gominio

type Config struct {
	AccessKeyId     string `json:",optional" yaml:"AccessKeyId"`
	AccessKeySecret string `json:",optional" yaml:"AccessKeySecret"`
	Endpoint        string `json:",optional" yaml:"Endpoint"`
	Bucket          string `json:",optional" yaml:"Bucket"`
	Domain          string `json:",optional" yaml:"Domain,optional"`
	UseSSL          bool   `json:",optional" yaml:"UseSSL"`
}
