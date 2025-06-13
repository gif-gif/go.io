package goattribution

import (
	"net/url"
)

type OrganicHandler struct {
}

func (h *OrganicHandler) Channel() string {
	return CHANNEL_ORGANIC
}

func (h *OrganicHandler) Match(queryParams url.Values) bool {
	return queryParams.Get("utm_medium") == h.Channel() || (queryParams.Get("utm_source") == "(not set)" && queryParams.Get("utm_medium") == "(not set)")
}

func (h *OrganicHandler) Handle(queryParams url.Values) (*AttributeInfo, error) {
	return &AttributeInfo{
		Channel: h.Channel(),
	}, nil
}
