package gominio

type Config struct {
	Name            string `json:"name,optional"  yaml:"Name"`
	AccessKeyId     string `json:"accessKeyId,optional"  yaml:"AccessKeyId"`
	AccessKeySecret string `json:"accessKeySecret,optional"  yaml:"AccessKeySecret"`
	Endpoint        string `json:"endpoint,optional"  yaml:"Endpoint"`
	//Bucket          string `json:"bucket,optional"  yaml:"Bucket"`
	//Dir    string `json:"dir,optional"  yaml:"Dir"`
	Domain string `json:"domain,optional"  yaml:"Domain,optional"`
	UseSSL bool   `json:"useSSL,optional"  yaml:"UseSSL"`
	Open   bool   `json:"open,optional"  yaml:"Open"`
}
