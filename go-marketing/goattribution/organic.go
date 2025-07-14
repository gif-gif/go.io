package goattribution

import (
	"net/url"
	"strings"
)

type OrganicHandler struct {
}

func (h *OrganicHandler) Channel() string {
	return CHANNEL_ORGANIC
}

func (h *OrganicHandler) Match(queryParams url.Values) bool {
	utm_medium := strings.TrimSpace(queryParams.Get("utm_medium"))
	utm_source := strings.TrimSpace(queryParams.Get("utm_source"))
	return utm_medium == h.Channel() || utm_source == h.Channel() ||
		(utm_source == "(not set)" && utm_medium == "(not set)") ||
		(utm_source == "" && utm_medium == "")
}

func (h *OrganicHandler) Handle(queryParams url.Values) (*AttributeInfo, error) {
	return CreateBaseAttributeInfo(queryParams, h.Channel()), nil
}
