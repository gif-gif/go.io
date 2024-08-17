package goalioss

type Config struct {
	Name            string `json:"name,optional"  yaml:"Name"`
	AccessKeyId     string `json:"accessKeyId,optional"  yaml:"AccessKeyId"`
	AccessKeySecret string `json:"accessKeySecret,optional"  yaml:"AccessKeySecret"`
	Endpoint        string `json:"endpoint,optional"  yaml:"Endpoint"`
	Bucket          string `json:"bucket,optional"  yaml:"Bucket"`
	Domain          string `json:"domain,optional"  yaml:"Domain"`
	Open            bool   `json:"open,optional"  yaml:"Open"`
}
