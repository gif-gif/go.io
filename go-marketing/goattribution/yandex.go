package goattribution

import (
	"net/url"
	"strings"
)

type YandexAttributeHandler struct {
}

func (h *YandexAttributeHandler) Channel() string {
	return CHANNEL_YANDEX
}

func (h *YandexAttributeHandler) Match(queryParams url.Values) bool {
	utm_medium := strings.TrimSpace(queryParams.Get("utm_medium"))
	utm_source := strings.TrimSpace(queryParams.Get("utm_source"))
	return utm_medium == h.Channel() || utm_source == h.Channel()
}

func (h *YandexAttributeHandler) Handle(queryParams url.Values) (*AttributeInfo, error) {
	return CreateBaseAttributeInfo(queryParams, h.Channel()), nil
}
