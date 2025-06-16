package goattribution

import (
	"fmt"
	"slices"
)

type (
	AttributeManager struct {
		AttributeHandlers []AttributeHandler
		Config            Config
	}
)

func New(config Config) *AttributeManager {
	if config.DecryptKeys == nil {
		config.DecryptKeys = make(map[string]string)
	}
	return &AttributeManager{
		AttributeHandlers: []AttributeHandler{
			&FacebookAttributeHandler{
				DecryptKey: config.DecryptKeys[CHANNEL_META],
			},
			&AppsFlyerAttributeHandler{},
			&BigoAttributeHandler{},
			&OrganicHandler{},
			&GoogleAttributeHandler{},
			//新的渠道处理器需要放在这里
			&CommonAttributeHandler{},
		},
		Config: config,
	}
}

// 注册一个属性处理器
func (m *AttributeManager) AddAttributeHandler(handler AttributeHandler) {
	m.AttributeHandlers = slices.Insert(m.AttributeHandlers, 0, handler)
}

func (m *AttributeManager) DecryptAttribute(referer string) (*AttributeInfo, error) {
	queryParams, err := ParseQuery(referer)
	if err != nil {
		return nil, fmt.Errorf("parse query error: %w", err)
	}
	for _, handler := range m.AttributeHandlers {
		if handler.Match(queryParams) {
			return handler.Handle(queryParams)
		}
	}

	return &AttributeInfo{
		Channel: "organic_unknown",
	}, nil
}
