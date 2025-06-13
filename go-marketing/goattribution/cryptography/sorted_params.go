package cryptography

import (
	"encoding/json"
	"sort"
	"strings"
)

func GetSortedString(info any, ignoreKeys map[string]struct{}, ignoreEmpty bool) (string, error) {
	dataMap, ok := info.(map[string]string)
	if !ok {
		var newDataMap map[string]string
		dataJson, err := json.Marshal(info)
		if err != nil {
			return "", err
		}
		err = json.Unmarshal(dataJson, &newDataMap)
		if err != nil {
			return "", err
		}
		dataMap = newDataMap
	}

	keys := make([]string, 0, len(dataMap)-1)
	for key, val := range dataMap {
		if _, ok := ignoreKeys[key]; ok {
			continue
		}
		if ignoreEmpty {
			if val == "" || val == "0" || val == "0.0" {
				continue
			}
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	builder := strings.Builder{}
	for _, key := range keys {
		builder.WriteString(key)
		builder.WriteString("=")
		builder.WriteString(dataMap[key])
		builder.WriteString("&")
	}

	queryString := builder.String()
	queryString = queryString[:len(queryString)-1]
	return queryString, nil
}
