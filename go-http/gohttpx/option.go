package gohttpx

type Option struct {
	Name  string
	Value interface{}
}

func TlsOption(caCrtFile, clientCrtFile, clientKeyFile string) Option {
	return Option{Name: "tls", Value: map[string]string{
		"caCrtFile":     caCrtFile,
		"clientCrtFile": clientCrtFile,
		"clientKeyFile": clientKeyFile,
	}}
}

func ContentTypeXmlOption() Option {
	return Option{Name: "content-type-xml", Value: CONTENT_TYPE_XML}
}

func ContentTypeJsonOption() Option {
	return Option{Name: "content-type-xml", Value: CONTENT_TYPE_JSON}
}

func ContentTypeFormOption() Option {
	return Option{Name: "content-type-xml", Value: CONTENT_TYPE_FORM}
}

func HeaderOption(field, value string) Option {
	return Option{Name: "header", Value: map[string]string{
		field: value,
	}}
}

func DebugOption() Option {
	return Option{Name: "debug", Value: true}
}
