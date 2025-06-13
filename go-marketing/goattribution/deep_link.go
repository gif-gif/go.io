package goattribution

import (
	"fmt"
	"net/url"
	"strconv"
)

// referer示例:
//1. gclid=EAIaIQobChMIwLX_jcuziwMV1olQBh1cJASjEAEYASAAEgKDQPD_BwE&gbraid=0AAAAA-VTNA1RZmc_uqUMRcothWf8n9V1w&gad_source=5
//2. utm_source=(not%20set)&utm_medium=(not%20set)
//3. 客户端写死默认值 utm_source=google-play&utm_medium=organic

func ParseQuery(referer string) (url.Values, error) {
	// 对Unicode字符解码
	decodedStr, err := strconv.Unquote(`"` + referer + `"`)
	if err != nil {
		return nil, fmt.Errorf("decoding string error: %v", err)
	}

	// 解析URL参数
	return url.ParseQuery(decodedStr)
}
