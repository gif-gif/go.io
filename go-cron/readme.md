# 定时任务

```
gocron.Minute(func() {
    log.Println("Minute - 1")
}, func() {
    log.Println("Minute - 2")
}).MinuteX(2, func() {
    log.Println("MinuteX - 1")
}, func() {
    log.Println("MinuteX - 2")
}).Start()
```