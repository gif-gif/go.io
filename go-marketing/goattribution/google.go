package goattribution

import (
	"net/url"
)

type GoogleAttributeHandler struct {
}

func (h *GoogleAttributeHandler) Channel() string {
	return CHANNEL_GOOGLE
}

// gclid=123456789&utm_medium=referral&utm_source=apps.facebook.com&utm_campaign=fb4a&utm_content=bytedanceglobal_E.C.P.C&facebook_app_id=
func (h *GoogleAttributeHandler) Match(queryParams url.Values) bool {
	return len(queryParams.Get("gclid")) > 0 //ads
}

func (h *GoogleAttributeHandler) Handle(queryParams url.Values) (*AttributeInfo, error) {
	return CreateBaseAttributeInfo(queryParams, h.Channel()), nil
}
