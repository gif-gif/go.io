package goerror

var (
	message = map[uint32]string{
		OK:                            "ok",
		SERVER_COMMON_ERROR:           "server error",
		REQUEST_PARAM_ERROR:           "param error",
		TOKEN_EXPIRE_ERROR:            "token expired,please login again",
		TOKEN_GENERATE_ERROR:          "token create error",
		FORBIDDEN_ERROR:               "forbidden",
		DB_ERROR:                      "db error",
		DB_UPDATE_AFFECTED_ZERO_ERROR: "no changes",
		CAPTCHA_ERROR:                 "captcha error",
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
)

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
