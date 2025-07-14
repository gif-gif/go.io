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
	utm_medium := queryParams.Get("utm_medium")
	utm_source := queryParams.Get("utm_source")
	if utm_source != "" {
		h._Channel = utm_source
	} else if utm_medium != "" {
		h._Channel = utm_medium
	} else {
		h._Channel = CHANNEL_UNKWON
	}

	return CreateBaseAttributeInfo(queryParams, h.Channel()), nil
}
