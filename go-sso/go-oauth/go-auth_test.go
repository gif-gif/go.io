package gooauth

import (
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"testing"
)

func TestAdmobAuthUrl(t *testing.T) {
	Init(Config{
		Name: "test",
		OAuthConfig: oauth2.Config{
			ClientID:     "123",
			ClientSecret: "123",
			RedirectURL:  "URL_ADDRESS",
			Scopes:       []string{"URL_ADDRESS"},
			Endpoint:     google.Endpoint,
		},
		State: "test",
	})
	url := GetClient("test").AuthUrl()
	golog.WithTag("admob").Info(url)

	golog.WithTag("admob").Info(goutils.UrlDecode("https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fadmob.readonly+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fadmob.report"))

}
