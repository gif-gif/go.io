package goattribution

import (
	"net/url"
)

type GoogleAttributeHandler struct {
}

func (h *GoogleAttributeHandler) Channel() string {
	return CHANNEL_GOOGLE
}

func (h *GoogleAttributeHandler) Match(queryParams url.Values) bool {
	return len(queryParams.Get("gclid")) > 0
}

func (h *GoogleAttributeHandler) Handle(queryParams url.Values) (*AttributeInfo, error) {
	return &AttributeInfo{
		UtmSource:  queryParams.Get("utm_source"),
		UtmMedium:  queryParams.Get("utm_medium"),
		UtmContent: queryParams.Get("utm_content"),
		Channel:    h.Channel(),
	}, nil
}
