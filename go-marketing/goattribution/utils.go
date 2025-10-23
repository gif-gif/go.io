package goattribution

import "net/url"

func CreateBaseAttributeInfo(queryParams url.Values, channel string, subChannel string) *AttributeInfo {
	clickId := queryParams.Get("click_id")
	if clickId == "" {
		clickId = queryParams.Get("clickid")
	}
	if clickId == "" {
		clickId = queryParams.Get("sid")
	}
	return &AttributeInfo{
		UtmSource:  queryParams.Get("utm_source"),
		UtmMedium:  queryParams.Get("utm_medium"),
		UtmContent: queryParams.Get("utm_content"),
		ClickId:    clickId,
		Channel:    channel,
		SubChannel: subChannel,
	}
}
