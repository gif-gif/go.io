# 有状态信息的错误封装

```go
package main

import (
	goerror "github.com/gif-gif/go.io/go-error"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/pkg/errors"
)

func main() {
	err := errors.Wrapf(goerror.NewErrCode(goerror.DB_ERROR), "find customer db err, in:%v , err:%v", "args", "err")
	if goerror.IsCodeErr(goerror.DB_ERROR) {
		golog.WithTag("goerror1").Error(err.Error(), " db error")
	}

	if goerror.IsErrCode(err, goerror.DB_ERROR) {
		golog.WithTag("goerror2").Error(err.Error(), " db error")
	}
	errCode, errMsg := goerror.GetErrCodeMsg(err)
	golog.WithTag("goerror3").Error(errCode, errMsg)
}

```
