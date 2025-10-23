package goattribution

import (
	"net/url"
	"strings"

	"github.com/samber/lo"
)

type YandexAttributeHandler struct {
	_Channel    string
	_SubChannel string
}

func (h *YandexAttributeHandler) SubChannel() string {
	return h._SubChannel
}

func (h *YandexAttributeHandler) Channel() string {
	return CHANNEL_YANDEX
}

func (h *YandexAttributeHandler) Match(queryParams url.Values) bool {
	utm_medium := strings.TrimSpace(queryParams.Get("utm_medium"))
	utm_source := strings.TrimSpace(queryParams.Get("utm_source"))
	return utm_medium == CHANNEL_YANDEX || utm_source == CHANNEL_YANDEX
}

func (h *YandexAttributeHandler) Handle(queryParams url.Values) (*AttributeInfo, error) {
	h._Channel = CHANNEL_YANDEX
	utm_medium := queryParams.Get("utm_medium")
	//utm_source := queryParams.Get("utm_source")
	h._SubChannel = lo.If(strings.TrimSpace(utm_medium) != "", utm_medium).Else(CHANNEL_YANDEX)
	return CreateBaseAttributeInfo(queryParams, h.Channel(), h.SubChannel()), nil
}
