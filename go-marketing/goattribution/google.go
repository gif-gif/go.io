package goattribution

import (
	"net/url"

	"github.com/samber/lo"
)

type GoogleAttributeHandler struct {
	_Channel    string
	_SubChannel string
}

func (h *GoogleAttributeHandler) SubChannel() string {
	return h._SubChannel
}

func (h *GoogleAttributeHandler) Channel() string {
	return CHANNEL_GOOGLE
}

// gclid=123456789&utm_medium=referral&utm_source=apps.facebook.com&utm_campaign=fb4a&utm_content=bytedanceglobal_E.C.P.C&facebook_app_id=
func (h *GoogleAttributeHandler) Match(queryParams url.Values) bool {
	return len(queryParams.Get("gclid")) > 0 //ads
}

func (h *GoogleAttributeHandler) Handle(queryParams url.Values) (*AttributeInfo, error) {
	utm_medium := queryParams.Get("utm_medium")
	//utm_source := queryParams.Get("utm_source")
	h._Channel = CHANNEL_GOOGLE
	h._SubChannel = lo.If(utm_medium != "", utm_medium).Else(CHANNEL_GOOGLE)

	return CreateBaseAttributeInfo(queryParams, h.Channel(), h.SubChannel()), nil
}
