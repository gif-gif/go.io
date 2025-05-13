package gowechat

import (
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/gogf/gf/util/gconv"
	"math"
	"sort"
	"strings"
	"time"
)

// 常用签名验证, sign Sha1 小写
func CheckSignSha1(secret, nonce string, linkSignTimeout int64, ts int64, sign string) bool {
	if linkSignTimeout == 0 {
		linkSignTimeout = 20
	}
	tsStep := time.Now().Unix() - ts
	if math.Abs(gconv.Float64(tsStep)) > gconv.Float64(linkSignTimeout) { //连接超时
		return false
	}

	args := []string{secret, gconv.String(ts), nonce}
	sort.Strings(args)
	// 小写
	serverSign := goutils.SHA1([]byte(strings.Join(args, "")))
	return serverSign == sign
}
