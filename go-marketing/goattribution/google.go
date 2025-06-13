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
		Channel: h.Channel(),
	}, nil
}
