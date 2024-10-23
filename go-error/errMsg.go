package goerror

var (
	errorsx = map[uint32]string{}

	message = map[uint32]string{
		OK:                    "ok",
		SERVER_COMMON_ERROR:   "server error",
		NOT_FOUND_ERROR:       "not found error",
		REQUEST_PARAM_ERROR:   "param error",
		TOKEN_EXPIRE_ERROR:    "token expired,please login again",
		USER_EXISTS_ERROR:     "user exists",
		USER_NOT_EXISTS_ERROR: "user not exists",
		USER_LOGIN_ERROR:      "username or password is invalid",
		FORBIDDEN_ERROR:       "forbidden",
		DB_ERROR:              "db error",
		CAPTCHA_ERROR:         "captcha error",
	}
)

var (
	ErrNotFound          = NewErrCodeMsg(NOT_FOUND_ERROR, "not found")
	ErrUserNoExists      = NewErrCodeMsg(USER_NOT_EXISTS_ERROR, "user not found")
	ErrUserExists        = NewErrCodeMsg(USER_EXISTS_ERROR, "user exists")
	ErrUserForbidden     = NewErrCode(FORBIDDEN_ERROR)
	ErrUserLogin         = NewErrCode(USER_LOGIN_ERROR)
	ErrUnauthorized      = NewErrCodeMsg(TOKEN_EXPIRE_ERROR, "unauthorized")
	ErrCaptcha           = NewErrCodeMsg(CAPTCHA_ERROR, "captcha error")
	ErrRequestParamError = NewErrCode(REQUEST_PARAM_ERROR)
	ErrDBError           = NewErrCode(DB_ERROR)
	ErrServerError       = NewErrCode(SERVER_COMMON_ERROR)
	ErrRedisError        = NewErrCode(REDIS_ERROR)
	ErrEtcdError         = NewErrCode(ETCD_ERROR)
)

// 扩展的错误类型这里初始化
func Init(errs map[uint32]string) {
	errorsx = errs
}

func MapErrMsg(errcode uint32) string {
	if msg, ok := message[errcode]; ok {
		return msg
	} else {
		return "server error"
	}
}

func IsCodeErr(errcode uint32) bool {
	if _, ok := message[errcode]; ok {
		return true
	} else {
		return false
	}
}
