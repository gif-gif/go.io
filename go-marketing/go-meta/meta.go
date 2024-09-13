package gometa

import (
	"context"
	gohttp "github.com/gif-gif/go.io/go-http"
	"github.com/google/go-querystring/query"
)

type Market struct {
	BaseApi     string
	AccessToken string
	StartDate   string
	EndDate     string
	PageSize    int
}

// 某个商户下所有账号信息 账号余额，状态等等
func (m *Market) GetAccountsByBusinessId(businessId string) (*AccountResponse, error) {
	req := &RequestData{
		AccessToken: m.AccessToken,
		Limit:       m.PageSize,
		Fields:      accountFields,
	}
	return m.getAccountsByBusinessId(req, businessId)
}

// 概要数据加载，下一页数据 res.Paging.Cursors.After
func (m *Market) GetAccountAdsetesOutline(accountId string) (*AllDataResponse, error) {
	req := &RequestData{
		Fields:      allDataFields,
		AccessToken: m.AccessToken,
		DateStart:   m.StartDate,
		DateStop:    m.EndDate,
		Limit:       m.PageSize,
	}

	res, err := m.getAllDataByAccountId(req, accountId, ApiAccountAdsets)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 某个计划或者广告组所有详情数据以 国家小时为纬度的数据
// res.Paging.Cursors.After 通过这个参数重新请求下一页数据
func (m *Market) GetDetailsData(outlineItem *AllDataItem) (*DataDetailResponse, error) {
	req := &DetailsDataRequest{
		RequestData: RequestData{
			Fields:      adFields,
			AccessToken: m.AccessToken,
			DateStart:   m.StartDate,
			DateStop:    m.EndDate,
			Limit:       m.PageSize,
		},
		Breakdowns: "['country']", //默认以国家纬度数据请求
	}

	res, err := m.getDetailByDataId(req, outlineItem.Id, ApiDataDetails)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// 刷新token接口
func (m *Market) RefreshToken(clientId string, clientSecret string) (*TokenResponse, error) {
	req := &ApiRefreshTokenRequest{
		ClientId:        clientId,
		ClientSecret:    clientSecret,
		GrantType:       "fb_exchange_token",
		FbExchangeToken: m.AccessToken,
	}
	api := m.BaseApi + ApiRefreshToken
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
