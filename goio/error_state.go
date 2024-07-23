package goio

const SUCCESS_REQUEST = 0

// 全局错误码

// client
const ERROR_REQUEST int64 = 400
const ERROR_TOKEN_EXPIRE int64 = 401
const ERROR_FORBIDDEN int64 = 403
const ERROR_INVALID int64 = -1

// server
const ERROR_SERVER int64 = 500
const ERROR_DATABASE int64 = 555
const ERROR_REDIS int64 = 556
const ERROR_PARAMS int64 = 700
