package goattribution

import (
	"net/url"
)

type CommonAttributeHandler struct {
	_Channel string
}

func (h *CommonAttributeHandler) Channel() string {
	return h._Channel
}

func (h *CommonAttributeHandler) Match(queryParams url.Values) bool {
	return true
}

func (h *CommonAttributeHandler) Handle(queryParams url.Values) (*AttributeInfo, error) {
	h._Channel = queryParams.Get("utm_medium")
	return &AttributeInfo{
		UtmSource:  queryParams.Get("utm_source"),
		UtmMedium:  queryParams.Get("utm_medium"),
		UtmContent: queryParams.Get("utm_content"),
		Channel:    h.Channel(),
	}, nil
}
