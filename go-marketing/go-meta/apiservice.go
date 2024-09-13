package gometa

import (
	"context"
	"fmt"
	gohttp "github.com/gif-gif/go.io/go-http"
	"github.com/google/go-querystring/query"
)

func (m *Market) getAccountsByBusinessId(req *RequestData, businessId string) (*AccountResponse, error) {
	api := m.BaseApi + ApiAccount
	api = fmt.Sprintf(api, businessId)
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
func (m *Market) getAllDataByAccountId(req *RequestData, accountId string, dataTypeUri string) (*AllDataResponse, error) {
	api := m.BaseApi + dataTypeUri
	api = fmt.Sprintf(api, accountId)
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
func (m *Market) getDetailByDataId(req *DetailsDataRequest, dataId string, dataTypeUri string) (*DataDetailResponse, error) {
	api := m.BaseApi + dataTypeUri
	api = fmt.Sprintf(api, dataId)
	req.TimeRange = "{'since':'" + req.DateStart + "','until':'" + req.DateStop + "'}"

	req.DateStart = ""
	req.DateStop = ""

	if req.Timezone == "" {
		//req.Timezone = "Asia/Shanghai"
	}

	if req.TimeIncrement <= 0 { //默认以天为单位
		req.TimeIncrement = 1 //以天为单位返回
	}

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
