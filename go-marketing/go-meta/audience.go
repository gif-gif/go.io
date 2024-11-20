package gometa

import (
	"context"
	"fmt"
	gohttp "github.com/gif-gif/go.io/go-http"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/gogf/gf/util/gconv"
	"github.com/google/go-querystring/query"
	"strings"
	"time"
)

type ResponseItem struct {
	AppVersionName  string  `json:"app_version_name"`
	Platform        string  `json:"platform"`
	Format          string  `json:"format"`
	Date            string  `json:"date"`
	Hour            int64   `json:"hour"`
	AdUnit          string  `json:"ad_unit"`
	Country         string  `json:"country"`
	AdRequest       int64   `json:"ad_requests"`
	Clicks          int64   `json:"clicks"`
	Earnings        int64   `json:"earnings"` // 美分*10000. 这里如果 返回美分会损失精度
	Impressions     int64   `json:"impressions"`
	ImpressionCtr   float64 `json:"impression_ctr"`
	ImpressionRpm   float64 `json:"Impression_rpm"` //美元
	MatchedRequests int64   `json:"matched_requests"`
	MatchRate       float64 `json:"match_rate"`
	ShowRate        float64 `json:"show_rate"`
}

type ReportResponse struct {
	Items         []*ResponseItem
	NextPageToken string
}

// <ID> 是您的 Meta 企业编号、资产编号或应用编号
func (m *GoMeta) GetMetaReport(req *AudienceDataRequest, ID string) (*AudienceDataResponse, error) {
	if req.Limit > limitMax {
		return nil, fmt.Errorf("最大查询不能超过：" + gconv.String(limitMax) + "条")
	}
	api := m.Config.currentVersionBaseApi + ApiAdNetworkAnalytics
	api = fmt.Sprintf(api, ID)
	params, _ := query.Values(req)
	request := &gohttp.Request{
		Url:          api,
		ParamsValues: params,
	}
	gh := gohttp.GoHttp[AudienceDataResponse]{Request: request}
	result, err := gh.HttpGet(context.Background())
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *GoMeta) GetReport(req *AudienceDataRequest, ID string) (*ReportResponse, error) {
	rst, err := m.GetMetaReport(req, ID)
	if err != nil {
		return nil, err
	}
	res := ReportResponse{}
	rowsData := make(map[string]map[string]string) //一行数据, 列转行
	for _, item := range rst.Data {
		for _, result := range item.Results {
			key := result.Time
			for _, breakdown := range result.Breakdowns {
				key += "-" + breakdown.Key + ":" + breakdown.Value + ""
			}

			rowData := make(map[string]string)
			if rowsData[key] == nil {
				rowsData[key] = rowData
			} else {
				rowData = rowsData[key]
			}

			rowData[result.Metric] = result.Value
			rowsData[key] = rowData
		}
	}

	for rowsKey, value := range rowsData {
		vo := ResponseItem{}
		keys := strings.Split(rowsKey, "-")
		for i, key := range keys {
			if i == 0 {
				tt, err := goutils.ConvertToGMTTime(key)
				if err != nil {
					fmt.Println(err)
					return nil, err
				}
				// 将时间对象转换为 UTC 时区
				utcTime := tt.UTC()
				gmtTimeStr := utcTime.Format(time.DateOnly)
				hour := utcTime.Hour()
				vo.Date = gmtTimeStr
				vo.Hour = int64(hour)
			} else {
				kk := strings.Split(key, ":")
				switch kk[0] {
				case BREAKDOWN_COUNTRY:
					vo.Country = kk[1]
					break
				case BREAKDOWN_PLACEMENT:
					vo.AdUnit = kk[1]
					break
				case BREAKDOWN_PLATFORM:
					vo.Platform = kk[1]
					break
				}
			}
		}

		for metric, v := range value {
			switch metric {
			case Metrics_AD_NETWORK_REVENUE:
				vo.Earnings = gconv.Int64(gconv.Float64(v) * 1000000) //1美元=1000000微分, market api 统一返回格式
				break
			case Metrics_AD_NETWORK_CPM:
				vo.ImpressionRpm = gconv.Float64(v)
				break
			case Metrics_AD_NETWORK_IMP:
				vo.Impressions = gconv.Int64(v)
				break
			case Metrics_AD_NETWORK_REQUEST:
				vo.AdRequest = gconv.Int64(v)
				break
			case Metrics_AD_NETWORK_SHOW_RATE:
				vo.ShowRate = gconv.Float64(v)
				break
			case Metrics_AD_NETWORK_CTR:
				vo.ImpressionCtr = gconv.Float64(v)
				break
			case Metrics_AD_NETWORK_CLICK:
				vo.Clicks = gconv.Int64(v)
				break
			case Metrics_AD_NETWORK_FILL_RATE:
				vo.MatchRate = gconv.Float64(v)
				break
			case Metrics_AD_NETWORK_FILLED_REQUEST:
				vo.MatchedRequests = gconv.Int64(v)
				break
			}
		}
		res.Items = append(res.Items, &vo)
	}

	res.NextPageToken = rst.Paging.Cursors.After

	return &res, nil
}
