package gooauth

import (
	golog "github.com/gif-gif/go.io/go-log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"testing"
)

func TestAdmobAuthUrl(t *testing.T) {
	Init(Config{
		Name:         "test",
		ClientId:     "123",
		ClientSecret: "123",
		RedirectURL:  "URL_ADDRESS",
		Scopes:       []string{"URL_ADDRESS"},
		Endpoint: &Endpoint{
			AuthURL:  google.Endpoint.AuthURL,
			TokenURL: google.Endpoint.TokenURL,
			//AuthStyle: google.Endpoint.AuthStyle,
			AuthStyle: oauth2.AuthStyleInHeader,
		},
		State: "test",
	})

	param1 := oauth2.SetAuthURLParam("param1", "value1")
	param2 := oauth2.SetAuthURLParam("param2", "value2")
	url := GetClient("test").AuthUrl(param1, param2)
	golog.WithTag("admob").Info(url)

	//golog.WithTag("admob").Info(goutils.UrlDecode("https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fadmob.readonly+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fadmob.report"))

}
