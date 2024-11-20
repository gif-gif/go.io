package gometa

import (
	"context"
	"fmt"
	gohttp "github.com/gif-gif/go.io/go-http"
	"github.com/gogf/gf/util/gconv"
	"github.com/google/go-querystring/query"
	"time"
)

func (m *GoMeta) handleRequest(req *RequestData) *RequestData {
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

func (m *GoMeta) DecryptEcpms(appId string, encryptedEcpms []string) (*EncryptedEcpmRes, error) {
	api := m.Config.baseApi + "/" + appId + "/aggregate_revenue"
	req := EncryptedEcpmReq{
		AccessToken: m.Config.AccessToken,
		Ecpms:       encryptedEcpms,
		RequestId:   gconv.String(time.Now().UnixNano()),
		SyncApi:     true,
	}

	request := &gohttp.Request{
		Url:  api,
		Body: req,
	}

	gh := gohttp.GoHttp[EncryptedEcpmRes]{Request: request}
	result, err := gh.HttpPostJson(context.Background())

	if err != nil {
		return nil, err
	}
	return result, nil
}

// 某个商户下所有账号信息 账号余额，状态等等
func (m *GoMeta) GetMarketAccountsByBusinessId(businessId string, pageSize int) (*AccountResponse, error) {
	if pageSize == 0 {
		pageSize = 10000
	}
	api := m.Config.currentVersionBaseApi + ApiAccount
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
func (m *GoMeta) GetMarketAllDataByAccountId(req *RequestData, accountId string) (*AllDataResponse, error) {
	api := m.Config.currentVersionBaseApi + ApiAccountAdsets
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
func (m *GoMeta) GetMarketDetailByDataId(req *RequestData, dataId string) (*DataDetailResponse, error) {
	api := m.Config.currentVersionBaseApi + ApiDataDetails
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

//---------------------------------------------------------------- 使用例子方法----------------------------------------------------------------

// 某个计划或者广告组所有详情数据以 国家小时为纬度的数据
// res.Paging.Cursors.After 通过这个参数重新请求下一页数据
func (m *GoMeta) GetMarketDetailsDataForCountry(outlineItem *AllDataItem, startDate, endDate string, pageSize int) (*DataDetailResponse, error) {
	req := &RequestData{
		Fields:      adFields,
		AccessToken: m.Config.AccessToken,
		DateStart:   startDate,
		DateStop:    endDate,
		Limit:       pageSize,
		Breakdowns:  "['country']", //默认以国家纬度数据请求
	}

	res, err := m.GetMarketDetailByDataId(req, outlineItem.Id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// 概要数据加载，下一页数据 res.Paging.Cursors.After
func (m *GoMeta) GetMarketAccountAdSetsOutline(accountId string, startDate, endDate string, pageSize int) (*AllDataResponse, error) {
	req := &RequestData{
		Fields:      allDataFields,
		AccessToken: m.Config.AccessToken,
		DateStart:   startDate,
		DateStop:    endDate,
		Limit:       pageSize,
	}

	res, err := m.GetMarketAllDataByAccountId(req, accountId)
	if err != nil {
		return nil, err
	}
	return res, nil
}
