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
			&YandexAttributeHandler{},
			&GoogleAttributeHandler{},
			&FacebookAttributeHandler{
				DecryptKey: config.DecryptKeys[CHANNEL_META],
			},
			&AppsFlyerAttributeHandler{},
			&BigoAttributeHandler{},
			&OrganicHandler{},
			&CommonAttributeHandler{},
		},
		Config: config,
	}
}

// 注册一个属性处理器
// 注意：注册的顺序很重要，因为会按照注册的顺序进行匹配
// 例如：如果先注册了 FacebookAttributeHandler，那么就会先匹配 FacebookAttributeHandler，然后再匹配 AppsFlyerAttributeHandler
// 这个函数会把handler插入到AttributeHandlers的头部
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
			v, err := handler.Handle(queryParams)
			if err != nil {
				continue
			}
			return v, nil
		}
	}

	return &AttributeInfo{
		Channel: "organic_unknown",
	}, nil
}
