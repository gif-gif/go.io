# 常用工具包
#### 加密相关
- 生成 AES 密钥
- 生成 AES 密钥和 IV
- 计算文件md5(支持超大文件)
- 计算MD5大写、计算Md5小写
- SHA1、SHA256、HMacMd5、HMacSha1、SHAWithRSA
- Base64Encode、Base64Decode
- SHAWithRSA

#### email
- HideEmail
- IsEmail

#### 大数据计算
- bigint 计算

#### Goroutine 
- NewErrorGroup(context.TODO(), maxWorkers) 并发可取消的协程组
- AsyncFunc
- AsyncFuncPanic
- AsyncFuncGroup
- AsyncFuncGroupPanic
- MeasureExecutionTime
- IsContextDone

#### Utils

- func CheckSign(secret string, linkSignTimeout int64, ts int64, sign string) bool // 常用签名验证, sign md5 小写
- func CheckSignSha1(secret, nonce string, linkSignTimeout int64, ts int64, sign string) bool {
- func IsInArray[T any](arr []T, target T) bool // 元素都转换成字符串比较
- func IsInArrayX[T any](arr []T, exists func(target T) bool) bool 
- func IsInArrayXX[T any](arr []*T, exists func(target *T) bool) bool 
- func IfNot[T any](isTrue bool, a, b T) T  // 通用三目运算
- func IfString(isTrue bool, a, b string) string 
- func IfInt(isTrue bool, a, b int) int 
- func IfFloat32(isTrue bool, a, b float32) float32
- func IfFloat64(isTrue bool, a, b float64) float64 
- func ReverseArray(arr []*interface{}) 
- func PadStart(str, pad string, length int) string
- func MinInt64(a, b int64) int64
- func MinInt(a, b int) int
- func MaxInt64(a, b int64) int64 
- func MaxInt(a, b int) int
- func GenValidateCode(width int) string {//随机数
- func SplitStringArray(arr []string, size int) (list [][]string)
- func SplitIntArray(arr []int, size int) (list [][]int)
- func SplitInt64Array(arr []int64, size int) (list [][]int64)
- func SplitArray(arr []interface{}, size int) (list [][]interface{}) 
- func GetFieldValue(config interface{}, fieldName string) (interface{}, error) {// 通过反射获取结构体字段的值
- func CopyProperties[T any](target interface{}) // 复制对象
- func GetRuntimeStack() string // 获取运行堆栈
- func GenericSort  // 通用排序
- func InsertionSort // 通用插入排序
- func FillMissingNumbers 填充缺失数字
- func ConvertAmPmHourTo24HourFormat Am Pm 时间转换比

#### password
- func BcryptHash(password string) string // BcryptHash 使用 bcrypt 对密码进行加密
- func BcryptCheck(password, hash string) bool  // BcryptCheck 对比明文密码和数据库的哈希值
- func ValidPassword(str string) (msg string, matched bool) { //至少一位数字、大小字母,且长度6-20位
- func ValidPasswordV2(str string) (msg string, matched bool) { //至少一位数字、大小字母和特殊字符,且长度6-20位

#### time 时间
- time.go 
- timex.go

#### xml
- xml.go
-https://github.com/samber/lo
