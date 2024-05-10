# 日志级别

- Fatal 致命的
- Panic 严重的
- Error 错误的
- Warn 告警的
- Info 信息的
- Debug 调试的

# 文件说明

- `log.go` 对外开放的方法，默认console适配器，属于第1层级
- `logger.go` log对象，像 `SetAdapter` `WithHook` 是项目及全局方法，属于第2层级
- `entry.go` 实体类，每次产生一条log时，都要 `new` 实体类，主要处理消息内容、标签、附加数据等，属于第3层级
- `message.go` 消息内容对象，包装每条消息包含的字段信息
- `adapter.go` 适配器
- `console.go` 控制台适配器
- `file.go` 文件适配器
- `file_options.go` 文件适配器 选项

# 文件适配器

- `filepath` 日志保存目录
- `filename` 日志文件名，使用 `yyyymmdd.log`
- `maxSize` 文件大小最大值，默认512M，超过后，会切割文件
- `FilePathOption()` 定义文件路径
- `FileMaxSizeOption()` 定义文件大小最大值

# 输出对象

- console
- file
- kafka
- es