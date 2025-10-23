package goattribution

import (
	"net/url"
)

type CommonAttributeHandler struct {
	_Channel    string
	_SubChannel string
}

func (h *CommonAttributeHandler) SubChannel() string {
	return h._SubChannel
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
	h._Channel = CHANNEL_UNKWON
	h._SubChannel = CHANNEL_UNKWON
	if utm_source != "" {
		h._Channel = utm_source
	}

	if utm_medium != "" {
		h._SubChannel = utm_medium
	}

	return CreateBaseAttributeInfo(queryParams, h.Channel(), h.SubChannel()), nil
}
