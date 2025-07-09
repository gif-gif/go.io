package goattribution

import (
	"net/url"
)

type BigoAttributeHandler struct {
}

func (h *BigoAttributeHandler) Channel() string {
	return CHANNEL_BIGO
}

func (h *BigoAttributeHandler) Match(queryParams url.Values) bool {
	val := queryParams.Get("utm_medium")
	return val == h.Channel()
}

func (h *BigoAttributeHandler) Handle(queryParams url.Values) (*AttributeInfo, error) {
	return CreateBaseAttributeInfo(queryParams, h.Channel()), nil
}
