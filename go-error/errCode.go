package goerror

// 成功返回
const OK uint32 = 0

/**(前3位代表业务,后三位代表具体功能)**/

// 全局错误码
// client
const REQUEST_PARAM_ERROR uint32 = 400
const TOKEN_EXPIRE_ERROR uint32 = 401
const FORBIDDEN_ERROR uint32 = 403
const NOT_FOUND_ERROR uint32 = 404
const CAPTCHA_ERROR uint32 = 700
const USER_NOT_EXISTS_ERROR uint32 = 701
const USER_LOGIN_ERROR uint32 = 703
const USER_EXISTS_ERROR uint32 = 704

// server
const SERVER_COMMON_ERROR uint32 = 500
const DB_ERROR uint32 = 555
const REDIS_ERROR uint32 = 666
const ETCD_ERROR uint32 = 777

//用户模块
