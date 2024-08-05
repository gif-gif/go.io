package goalioss

type Config struct {
	AccessKeyId     string `json:"accessKeyId,optional" yaml:"AccessKeyId,optional"`
	AccessKeySecret string `json:"accessKeySecret,optional" yaml:"AccessKeySecret,optional"`
	Endpoint        string `json:"endpoint,optional" yaml:"Endpoint,optional"`
	Bucket          string `json:"bucket,optional" yaml:"Bucket,optional"`
	Domain          string `json:"domain,optional" yaml:"Domain,optional"`
}
