package gometa

import (
	"context"
	"fmt"
	gohttp "github.com/gif-gif/go.io/go-http"
	"github.com/google/go-querystring/query"
)

type Config struct {
	Name        string `yaml:"Name" json:"name,optional"`
	BaseApi     string `yaml:"BaseApi" json:"baseApi,optional"`
	AccessToken string `yaml:"AccessToken" json:"accessToken"`
	StartDate   string `yaml:"StartDate" json:"startDate"`
	EndDate     string `yaml:"EndDate" json:"endDate"`
	PageSize    int    `yaml:"PageSize" json:"pageSize,optional"`
}

type Market struct {
	Config Config
}

func New(config Config) *Market {
	if config.BaseApi == "" {
		config.BaseApi = "https://graph.facebook.com/v17.0"
	}
	return &Market{
		Config: config,
	}
}

func (m *Market) handleRequest(req *RequestData) *RequestData {
	req.TimeRange = "{'since':'" + req.DateStart + "','until':'" + req.DateStop + "'}"
	req.DateStart = ""
	req.DateStop = ""

	if req.Timezone == "" {
		//req.Timezone = "Asia/Shanghai"
	}
	if req.TimeIncrement <= 0 { //默认以天为单位
		req.TimeIncrement = 1 //以天为单位返回
	}

	return req
}

// 某个商户下所有账号信息 账号余额，状态等等
func (m *Market) GetAccountsByBusinessId(businessId string, pageSize int) (*AccountResponse, error) {
	if pageSize == 0 {
		pageSize = 10000
	}
	api := m.Config.BaseApi + ApiAccount
	api = fmt.Sprintf(api, businessId)

	req := &RequestData{
		AccessToken: m.Config.AccessToken,
		Limit:       pageSize,
		Fields:      accountFields,
	}

	req = m.handleRequest(req)
	params, _ := query.Values(req)
	request := &gohttp.Request{
		Url:          api,
		ParamsValues: params,
	}
	gh := gohttp.GoHttp[AccountResponse]{Request: request}
	result, err := gh.HttpGet(context.Background())
	if err != nil {
		return nil, err
	}
	return result, nil
}

// all data -------------------------------
func (m *Market) GetAllDataByAccountId(req *RequestData, accountId string) (*AllDataResponse, error) {
	api := m.Config.BaseApi + ApiAccountAdsets
	api = fmt.Sprintf(api, accountId)
	req = m.handleRequest(req)
	params, _ := query.Values(req)

	request := &gohttp.Request{
		Url:          api,
		ParamsValues: params,
	}
	gh := gohttp.GoHttp[AllDataResponse]{Request: request}
	result, err := gh.HttpGet(context.Background())

	if err != nil {
		return nil, err
	}
	return result, nil
}

// 根据数据类型获取某个详情，如：广告组详情 广告详情  -------------------------------
func (m *Market) GetDetailByDataId(req *RequestData, dataId string) (*DataDetailResponse, error) {
	api := m.Config.BaseApi + ApiDataDetails
	api = fmt.Sprintf(api, dataId)

	req = m.handleRequest(req)
	params, _ := query.Values(req)
	request := &gohttp.Request{
		Url:          api,
		ParamsValues: params,
	}
	gh := gohttp.GoHttp[DataDetailResponse]{
		Request: request,
	}
	result, err := gh.HttpGet(context.Background())
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 刷新token接口
func (m *Market) RefreshToken(clientId string, clientSecret string) (*TokenResponse, error) {
	req := &ApiRefreshTokenRequest{
		ClientId:        clientId,
		ClientSecret:    clientSecret,
		GrantType:       "fb_exchange_token",
		FbExchangeToken: m.Config.AccessToken,
	}
	api := m.Config.BaseApi + ApiRefreshToken
	params, _ := query.Values(req)

	request := &gohttp.Request{
		Url:          api,
		ParamsValues: params,
	}
	gh := gohttp.GoHttp[TokenResponse]{
		Request: request,
	}
	result, err := gh.HttpGet(context.Background())
	if err != nil {
		return nil, err
	}
	return result, nil
}

//
//func (m *Market) AccessKeys(accountId string) {
//	api := m.Config.BaseApi + "/" + accountId + "/access_keys"
//	request := &gohttp.Request{
//		Url:         api,
//		QueryParams: map[string]string{"access_token": m.Config.AccessToken},
//	}
//
//	gh := gohttp.GoHttp[TokenResponse]{
//		Request: request,
//	}
//
//	gh.HttpGet(context.Background())
//	//if err != nil {
//	//	return nil, err
//	//}
//	//return result, nil
//}

//---------------------------------------------------------------- 使用例子方法----------------------------------------------------------------

// 某个计划或者广告组所有详情数据以 国家小时为纬度的数据
// res.Paging.Cursors.After 通过这个参数重新请求下一页数据
func (m *Market) GetDetailsDataForCountry(outlineItem *AllDataItem, startDate, endDate string, pageSize int) (*DataDetailResponse, error) {
	req := &RequestData{
		Fields:      adFields,
		AccessToken: m.Config.AccessToken,
		DateStart:   startDate,
		DateStop:    endDate,
		Limit:       pageSize,
		Breakdowns:  "['country']", //默认以国家纬度数据请求
	}

	res, err := m.GetDetailByDataId(req, outlineItem.Id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// 概要数据加载，下一页数据 res.Paging.Cursors.After
func (m *Market) GetAccountAdSetsOutline(accountId string, startDate, endDate string, pageSize int) (*AllDataResponse, error) {
	req := &RequestData{
		Fields:      allDataFields,
		AccessToken: m.Config.AccessToken,
		DateStart:   startDate,
		DateStop:    endDate,
		Limit:       pageSize,
	}

	res, err := m.GetAllDataByAccountId(req, accountId)
	if err != nil {
		return nil, err
	}
	return res, nil
}
