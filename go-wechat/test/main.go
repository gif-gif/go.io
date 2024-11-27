package main

import (
	gocontext "github.com/gif-gif/go.io/go-context"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram/config"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/menu"
)

type AccessTokenHandle struct {
	AccessToken string
}

func (a *AccessTokenHandle) GetAccessToken() (string, error) {
	return a.AccessToken, nil
}

func main() {
	testMiniProgram()
	<-gocontext.WithCancel().Done()
}

func testMp() {
	wc := wechat.NewWechat()
	//这里本地内存保存access_token，也可选择redis，memcache或者自定cache
	memory := cache.NewMemory()
	cfg := &offConfig.Config{
		AppID:     "wxb19e7f16eafb98c2",
		AppSecret: "b38828d87586ec284b093e0d87ca7b21",
		Token:     "123qwe",
		//EncodingAESKey: "xxxx",
		Cache: memory,
	}
	oa := wc.GetOfficialAccount(cfg)

	//ak, err := oa.GetAccessToken()
	//if err != nil {
	//	golog.WithTag("ak").Error(err)
	//	return
	//}
	//golog.WithTag("ak").Info(ak)

	oa.SetAccessTokenHandle(&AccessTokenHandle{
		AccessToken: "84_FO2Fgjm2MU73ylsnWriDGQ7XXLhNE792znTlmR_Cr7lnBeVK0ifnUcH-tzroqgUh2rW3tvSxWkfuIYH6TrN64Pl6B_V8PHswkCB-ZUrv6acFiC49p1C7avNO4dANMBaABADCR", //"84_sgRh9c8kz4umPXfg0qPveHGdj4E5kg6IUQUBTqrHK_GLnqnBtNDJxzHJLuIkgO-G6ttW2m-eMDmSmBIbgd0YvkAiQ1-UMIv9-O5dgjCPxlgLwnWdwHr60FK7jlITOKfAFABWS",
	})

	//ipList, err := officialAccount.GetBasic().GetCallbackIP()
	ipList, err := oa.GetBasic().GetAPIDomainIP()

	if err != nil {
		golog.WithTag("ipList").Error(err)
		return
	}
	golog.WithTag("ipList").Info(ipList)

	bd := oa.GetBroadcast()
	r, err := bd.SendText(nil, "hello")
	if err != nil {
		golog.WithTag("bd").Error(err)
		return
	}

	golog.WithTag("bd").Info("ok", r)
	m := oa.GetMenu()
	var buttons []*menu.Button
	buttons = append(buttons, &menu.Button{
		Type: "click",
		Name: "今日歌曲",
		Key:  "V1001_TODAY_MUSIC",
		URL:  "https://wx.acom.cc",
		SubButtons: []*menu.Button{
			&menu.Button{
				Type: "click",
				Name: "今日歌曲1",
				Key:  "V1001_TODAY_MUSIC1",
				URL:  "https://wx.acom.cc",
			}, &menu.Button{
				Type: "click",
				Name: "今日歌曲2",
				Key:  "V1001_TODAY_MUSIC2",
				URL:  "https://wx.acom.cc",
			},
		},
	})
	err = m.SetMenu(buttons)
	if err != nil {
		golog.WithTag("m").Error(err)
		return
	}

	u := oa.GetUser()
	list, err := u.ListAllUserOpenIDs()
	if err != nil {
		golog.WithTag("list").Error(err)
		return
	}

	golog.WithTag("list").Info(list)

	golog.WithTag("m").Info("ok", r)
}

func testMiniProgram() {
	wc := wechat.NewWechat()
	memory := cache.NewMemory()
	cfg := &config.Config{
		AppID:     "wxb19e7f16eafb98c2",
		AppSecret: "b38828d87586ec284b093e0d87ca7b21",
		Token:     "123qwe",
		//EncodingAESKey: "xxxx",
		Cache: memory,
	}
	mini := wc.GetMiniProgram(cfg)
	a := mini.GetAuth()
	golog.WithTag("mini").Info(a)
}
