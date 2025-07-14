package goerror

// 常用错误码
// 成功返回
const OK uint32 = 0

/**(前3位代表业务,后三位代表具体功能)**/

// 全局错误码
// client
const REQUEST_PARAM_ERROR uint32 = 1400
const TOKEN_EXPIRE_ERROR uint32 = 1401
const FORBIDDEN_ERROR uint32 = 1403
const NOT_FOUND_ERROR uint32 = 1404
const CAPTCHA_ERROR uint32 = 1700
const USER_NOT_EXISTS_ERROR uint32 = 1701
const USER_LOGIN_ERROR uint32 = 1703
const USER_EXISTS_ERROR uint32 = 1704

// server
const SERVER_COMMON_ERROR uint32 = 1500
const DB_ERROR uint32 = 1555
const REDIS_ERROR uint32 = 1666
const ETCD_ERROR uint32 = 1777

// 错误处理枚举
// showType?: number; // error display type： 0 silent; 1 message.warn; 2 message.error; 4 notification; 9 page
const ShowTypeSilent uint32 = 0
const ShowTypeMessageWarn uint32 = 1
const ShowTypeMessageError uint32 = 2
const ShowTypeNotification uint32 = 4
const ShowTypePage uint32 = 9 // 页面跳转
