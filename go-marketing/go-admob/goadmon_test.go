package goadmob

import (
	"context"
	golog "github.com/gif-gif/go.io/go-log"
	"golang.org/x/oauth2"
	"google.golang.org/api/admob/v1"
	"google.golang.org/api/option"
	"testing"
)

func TestAdmob(t *testing.T) {
	ctx := context.Background()
	Init(Config{
		ClientId:     ".apps.googleusercontent.com",
		ClientSecret: "",
		RedirectUrl:  "URL_ADDRESS",
		State:        "123456",
	})
	select {}
}

// ya29.a0AeDClZA_Fa7P107WdgLtsygrei4DadnbugnQ5HAQ7qxCC5o6d3_VV-8VqOvZNH2MeqATbpGf-vRgnRGpA8RQgAEjFBVtJvbLdlKtQ-kXOviFIAM1QjwgZrjV-ipxUk0tBZ5HB6zt1Qa5o8AITh7VP29TGfRK2lT7uxkgvobcaCgYKAYwSARISFQHGX2MiYWz4pU5sbWqu7s8B9PEICg0175
// 1//0ecJCv_hmW44mCgYIARAAGA4SNwF-L9Irk_ZNfBfC28DxaCWrUwsMVo0_wZhV9xdHQnU4Z1wNXSUwRUeiCo7wM-FaQ0HRtogYBwE

func TestEtcdConfigListener(t *testing.T) {
	ctx := context.Background()
	config := &oauth2.Config{
		ClientID:     "273488495628-a56cdd6vrnkm5i5ors5vl2bmrj3rh622.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "https://bangbox.jidianle.cc",
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://accounts.google.com/o/oauth2/auth",
			TokenURL:  "https://oauth2.googleapis.com/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

	token, err := config.TokenSource(ctx, &oauth2.Token{
		RefreshToken: "1//0ecJCv_hmW44mCgYIARAAGA4SNwF-L9Irk_ZNfBfC28DxaCWrUwsMVo0_wZhV9xdHQnU4Z1wNXSUwRUeiCo7wM-FaQ0HRtogYBwE",
	}).Token()

	//token, err := config.Exchange(ctx, "4/0AVG7fiRTFumd6aP08kwz_FPtq7x6-1iEQpqUBK5UNlbJTQqVicCIgS4ritCFFQD3vDmGfw")
	if err != nil {
		golog.WithTag("goadmob").Error("token:" + err.Error())
		return
	}

	golog.WithTag("goadmob").Info(token.AccessToken, token.RefreshToken, token.Expiry)

	admobService, err := admob.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))
	if err != nil {
		golog.WithTag("goadmob").Error(err)
		return
	}
	res, err := admobService.Accounts.Apps.List("accounts/pub-9876543210987654").Do()
	if err != nil {
		golog.WithTag("goadmob").Error(err)
		return
	}

	golog.WithTag("goadmob").Info(res)

	select {}
}
