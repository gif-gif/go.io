package goattribution

import (
	"net/url"

	"github.com/samber/lo"
)

type BigoAttributeHandler struct {
	_Channel    string
	_SubChannel string
}

func (h *BigoAttributeHandler) SubChannel() string {
	return h._SubChannel
}

func (h *BigoAttributeHandler) Channel() string {
	return CHANNEL_BIGO
}

func (h *BigoAttributeHandler) Match(queryParams url.Values) bool {
	utm_medium := queryParams.Get("utm_medium")
	utm_source := queryParams.Get("utm_source")
	return utm_source == CHANNEL_BIGO || utm_medium == CHANNEL_BIGO
}

func (h *BigoAttributeHandler) Handle(queryParams url.Values) (*AttributeInfo, error) {
	h._Channel = CHANNEL_BIGO
	utm_medium := queryParams.Get("utm_medium")
	h._SubChannel = lo.If(utm_medium != "", utm_medium).Else(CHANNEL_BIGO)
	return CreateBaseAttributeInfo(queryParams, h.Channel(), h.SubChannel()), nil
}
