# Google admob data Api

## Admob API 授权
官方调试平台（授权 获取Token 访问）
```json 
https://developers.google.com/oauthplayground/
```

### 1 把 redirect_uri 改为自己的服务器地址

```json
https://accounts.google.com/o/oauth2/v2/auth?redirect_uri=https%3A%2F%2Fbangbox.jidianle.cc&prompt=consent&response_type=code&client_id=273488495628-a56cdd6vrnkm5i5ors5vl2bmrj3rh622.apps.googleusercontent.com&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fadmob.readonly+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fadmob.report&access_type=offline
```
#### 跳转后选择权限 需要授权账号

### 2 授权后会返回 code 
### 3 获取AccessToken 和 RefreshToken


### 4 需求
- 按产品平台（安卓/IOS）、国家（全部国家、美国、德国、荷兰、瑞典、英国、澳大利亚、加拿大、巴西、缅甸、伊朗）、天、天小时数据 广告类型和广告单元
