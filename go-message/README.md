# 发送通知
飞书
```
普通群消息
gomessage.FeiShu(hookUrl, "这是普通的群消息")
```

钉钉
``` 
gomessage.InitDing("token","secret")

普通群消息
err := gomessage.DingDing("这是普通的群消息")

@特定人的消息
@对象必须为绑定钉钉的手机号
err := gomessage.DingDing("Lucy, Harvey, 你们的程序挂了", "18578924567", "+13414567890")

@所有人的消息
err := gomessage.DingDing("这是@所有人的消息", "*")

```

