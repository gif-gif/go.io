package goattribution

import "net/url"

func CreateBaseAttributeInfo(queryParams url.Values, channel string) *AttributeInfo {
	return &AttributeInfo{
		UtmSource:  queryParams.Get("utm_source"),
		UtmMedium:  queryParams.Get("utm_medium"),
		UtmContent: queryParams.Get("utm_content"),
		ClickId:    queryParams.Get("click_id"),
		Channel:    channel,
	}
}
