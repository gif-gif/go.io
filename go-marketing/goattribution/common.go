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
	return CreateBaseAttributeInfo(queryParams, h.Channel()), nil
}
