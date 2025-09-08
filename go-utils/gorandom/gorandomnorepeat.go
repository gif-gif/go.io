package gorandom

import "math/rand/v2"

type WeightedItem struct {
	Value  interface{} // 元素值
	Weight float64     // 权重
}

// WeightedRandomSelector 加权随机选择器
type WeightedRandomSelector struct {
	items       []WeightedItem
	totalWeight float64
}

// NewWeightedRandomSelector 创建新的加权随机选择器
func NewWeightedRandomSelector(items []WeightedItem) *WeightedRandomSelector {
	selector := &WeightedRandomSelector{
		items: make([]WeightedItem, len(items)),
	}

	// 复制元素并计算总权重
	copy(selector.items, items)
	for _, item := range selector.items {
		selector.totalWeight += item.Weight
	}

	return selector
}

// SelectMultiple 选择多个不重复的元素
func (w *WeightedRandomSelector) SelectMultiple(count int) []WeightedItem {
	if count <= 0 || count > len(w.items) {
		return nil
	}

	// 创建可用元素的副本
	availableItems := make([]WeightedItem, len(w.items))
	copy(availableItems, w.items)
	totalWeight := w.totalWeight

	result := make([]WeightedItem, 0, count)

	for i := 0; i < count; i++ {
		if len(availableItems) == 0 {
			break
		}

		// 生成随机数
		randomWeight := rand.Float64() * totalWeight

		// 找到对应的元素
		currentWeight := 0.0
		selectedIndex := -1

		for j, item := range availableItems {
			currentWeight += item.Weight
			if randomWeight <= currentWeight {
				selectedIndex = j
				break
			}
		}

		// 如果没有找到（理论上不应该发生），选择最后一个
		if selectedIndex == -1 {
			selectedIndex = len(availableItems) - 1
		}

		// 添加选中的元素到结果中
		selectedItem := availableItems[selectedIndex]
		result = append(result, selectedItem)

		// 从可用元素中移除已选择的元素
		totalWeight -= selectedItem.Weight
		availableItems = append(availableItems[:selectedIndex], availableItems[selectedIndex+1:]...)
	}

	return result
}

// Select 选择单个元素
func (w *WeightedRandomSelector) Select() interface{} {
	if len(w.items) == 0 {
		return nil
	}

	randomWeight := rand.Float64() * w.totalWeight
	currentWeight := 0.0

	for _, item := range w.items {
		currentWeight += item.Weight
		if randomWeight <= currentWeight {
			return item.Value
		}
	}

	// 理论上不应该到达这里，但为了安全起见返回最后一个元素
	return w.items[len(w.items)-1].Value
}
