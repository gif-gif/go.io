package goo_mail

var (
	__mail iMail
)

func Init(conf Config) {
	__mail = New(conf)
}

func Send(msg Message) error {
	return __mail.Send(msg)
}
