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
	utm_medium := queryParams.Get("utm_medium")
	utm_source := queryParams.Get("utm_source")
	return utm_source == h.Channel() || utm_medium == h.Channel()
}

func (h *BigoAttributeHandler) Handle(queryParams url.Values) (*AttributeInfo, error) {
	return CreateBaseAttributeInfo(queryParams, h.Channel()), nil
}
