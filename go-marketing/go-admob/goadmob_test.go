package goadmob

import (
	"context"
	gocontext "github.com/gif-gif/go.io/go-context"
	"github.com/gif-gif/go.io/go-db/gogorm"
	golog "github.com/gif-gif/go.io/go-log"
	gooauth "github.com/gif-gif/go.io/go-sso/go-oauth"
	"google.golang.org/api/admob/v1"
	"testing"
	"time"
)

func init() {
	//gogorm.Init(gogorm.Config{
	//	DataSource: "root:111111@tcp(127.0.0.1)/admob?charset=utf8mb4&parseTime=True&loc=Local",
	//})
}

func TestAdmobApps(t *testing.T) {
	err := Default().Refresh()
	if err != nil {
		golog.WithTag("goadmob").Error(err.Error())
		return
	}

	apps, err := Default().GetApps(100, "")
	if err != nil {
		golog.WithTag("goadmob").Error(err.Error())
		return
	}

	type AdAccount struct {
		Id         int64
		AccountId  string
		Channel    string
		Status     int64
		CreateTime int64
		UpdateTime int64
	}

	adAccount := AdAccount{}
	db := gogorm.Default().DB
	err = db.Table("ad_account").Select("id,account_id,channel,status,create_time,update_time").First(&adAccount).Error
	if err != nil {
		golog.WithTag("godb").Error(err.Error())
		return
	}

	for _, app := range apps.Apps {
		err = db.Table("ad_app").Exec("insert into ad_app(app_code,platform,title,ad_account_id,app_pub_id,app_store_id,status,create_time,update_time)values(?,?,?,?,?,?,?,?,?)", app.LinkedAppInfo.DisplayName, app.Platform, app.LinkedAppInfo.DisplayName, adAccount.Id, app.AppId, app.LinkedAppInfo.AppStoreId, 1, time.Now().Unix(), time.Now().Unix()).Error
		if err != nil {
			golog.WithTag("godb").Error(err.Error())
			return
		}
		golog.WithTag("admob").WithField("appId", app.AppId).Info("OK")
	}

	<-gocontext.WithCancel().Done()
}

type AdApp struct {
	Id       int64
	AppPubId string
}

func TestAdmobAdUnits(t *testing.T) {
	var adAppMap = make(map[string]int64)
	db := gogorm.Default().DB
	adApps := []AdApp{}
	err := db.Table("ad_app").Select("id,app_pub_id").Scan(&adApps).Error
	if err != nil {
		golog.WithTag("goadmob").Error(err.Error())
		return
	}

	for _, app := range adApps {
		adAppMap[app.AppPubId] = app.Id
	}

	err = Default().Refresh()
	if err != nil {
		golog.WithTag("goadmob").Error(err.Error())
		return
	}

	apps, err := Default().GetAdUnits(1000, "")
	if err != nil {
		golog.WithTag("goadmob").Error(err.Error())
		return
	}

	for _, app := range apps.AdUnits {
		appId := adAppMap[app.AppId]
		err = db.Table("ad_info").Exec("insert into ad_info(ad_app_id,title,ad_type,ad_unit,status,create_time,update_time)values(?,?,?,?,?,?,?)", appId, app.DisplayName, app.AdFormat, app.AdUnitId, 1, time.Now().Unix(), time.Now().Unix()).Error
		if err != nil {
			golog.WithTag("goadmob").Error(err.Error())
			return
		}
		golog.WithTag("admob").WithField("appId", app.AdUnitId).Info("OK")
	}

	<-gocontext.WithCancel().Done()
}

func TestAdmobReport(t *testing.T) {
	var adAppMap = make(map[string]int64)
	db := gogorm.Default().DB
	adApps := []AdApp{}
	err := db.Table("ad_app").Select("id,app_pub_id").Scan(&adApps).Error
	if err != nil {
		golog.WithTag("goadmob").Error(err.Error())
		return
	}

	for _, app := range adApps {
		adAppMap[app.AppPubId] = app.Id
	}

	err = Default().Refresh()
	if err != nil {
		golog.WithTag("goadmob").Error(err.Error())
		return
	}

	req := &ReportReq{
		MaxReportRows: 100,
		Dimensions:    []string{"DATE", "APP", "COUNTRY"},
		AdUnits:       []string{"ca-app-pub-4328354313035484/2826310677"},
		Metrics:       []string{"AD_REQUESTS", "CLICKS", "ESTIMATED_EARNINGS", "IMPRESSIONS", "IMPRESSION_CTR", "IMPRESSION_RPM", "MATCHED_REQUESTS", "MATCH_RATE", "SHOW_RATE"},
		EndDate: admob.Date{
			Day:   21,
			Month: 8,
			Year:  2024,
		},
		StartDate: admob.Date{
			Day:   20,
			Month: 8,
			Year:  2024,
		},
	}

	res, err := Default().GetReport(req)
	if err != nil {
		golog.WithTag("goadmob").Error(err.Error())
		return
	}

	golog.WithTag("admob").Info(res)

	<-gocontext.WithCancel().Done()
}

func TestAdmobAuthUrl(t *testing.T) {
	Init(context.Background(), Config{
		Name:      "admob",
		AccountId: "123",
		AuthConfig: gooauth.Config{
			ClientSecret: "secret",
			RedirectURL:  "https://test.com",
		},
	})
	url := Default().AuthUrl()
	golog.WithTag("admob").Info(url)

}
func TestAdmobAuth(t *testing.T) {

}
