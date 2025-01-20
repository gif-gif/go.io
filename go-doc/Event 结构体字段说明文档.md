# Event 结构体字段说明文档

## 基础信息字段
| 字段名 | 类型 | 说明 |
|--------|------|------|
| AppCode | string | 产品 code |
| UserId | string | 用户 id |
| RegisterDate | string | 用户注册日期 |
| RegisterTime | time.Time | 用户注册时间戳(秒) |

## 请求相关字段
| 字段名 | 类型 | 说明 |
|--------|------|------|
| RequestTime | time.Time | 请求时间戳，使用客户端 time 或 LogTime（两者差值小于 30 分钟时使用客户端 time，否则使用 LogTime）|
| RequestDate | string | 请求日期（2025-01-11） |
| Hour | uint8 | 当前小时 |

## 设备和系统信息
| 字段名 | 类型 | 说明 |
|--------|------|------|
| OsType | string | 设备类型（mobile-ios, mobile-android, desktop-windows, desktop-macos）|
| SystemVersion | string | 系统版本号 |
| Brand | string | 手机型号(参考：文档最后 设备列表)|
| DeviceId | string | 设备 ID |
| GoogleId | string | Google ID |
| AndroidId | string | Android ID |

## 渠道和广告系列信息
| 字段名 | 类型 | 说明 |
|--------|------|------|
| Channel | string | 渠道 |
| CampaignId | string | 系列 Id |
| CampaignName | string | 系列名 |
| AdCostMode | string | 广告计费模型 |
| CampaignChannel | string | 系列渠道（如 fb）|
| CampaignPartner | string | 系列合作人 |

## 地理位置和网络信息
| 字段名 | 类型 | 说明 |
|--------|------|------|


| Ip | string | 广告的 IP（IsVpn=1 时为 漂移国家ip，否则为 用户国家 ip）|
| VpnIp | string | VPN 节点 IP |
| VpnIsp | string | VPN IP 对应的运营商（使用客户端上报的 VpnIdc）|
| GeoIp | string | 地理位置 IP |
| GeoCountry | string | 用户国家（通过 API 获取）|
| GeoIsp | string | 地理位置运营商 |
| City | string | IP 对应城市（使用客户端上报的 VpnRegion代替）|
| VpnIdc | string | VPN 节点所在的 IDC 机房或运营商 |
| VpnRegion | string | VPN 节点所在的地区 |

## 事件相关字段
| 字段名 | 类型 | 说明 |
|--------|------|------|
| EventMode | string | 事件类型（client, server）|
| Event | string | 事件枚举值（ad_show, ad_click, ad_error, ad_req_success）|
| EventMsg | string | 事件信息（如 ad_error 的详细错误信息）|
| ErrorType | string | 错误类型 |

## 版本和状态标识
| 字段名 | 类型 | 说明 |
|--------|------|------|
| Vc | uint16 | APP 版本号（如：1）|
| Vn | string | APP 版本名（如：1.0.0）|
| IsVpn | uint8 | 是否使用 VPN 请求 |
| IsVpnAd | uint8 | 是否为 VPN 广告 |

## 广告相关字段
| 字段名 | 类型 | 说明 |
|--------|------|------|
| Ecpm | float64 | 每次展示的广告价值（没有时用 0）|
| Platform | string | 广告平台（如 AdMob,Meta, Applovin, Pangle 等）|
| AdType | string | 广告位类型（app_open/rewarded_video/interstitial/banner/native）|
| AdCountryCode | string | 广告的国家（漂移国家-小写国家code,如：美国-》us） |
| PlatformAccountCode | string | 广告平台账号 code |
| PlatformAccountName | string | 广告平台账号名称 |
| PlatformAccountAppId | string | 平台账号下的 AppId |
| AdPlacementCode | string | 广告位 code |
| AdPlacementName | string | 广告位名称 |
| AdSegmentCode | string | 广告分组 code |
| AdSegmentName | string | 广告分组名称 |
| MediaPlatform | string | 第三方聚合平台（如 TopOn, Max 等）|
| MediaSlotId | string | 第三方聚合平台广告位 ID(TopOnID) |
| SlotId | string | 实际的广告位 ID（如 AdMob 的广告位 ID）|
| LoadTime | int64 | 请求加载时间（单位：秒）|
| PagePos | string | 页面位置 |
| CacheAd | uint8 | 缓存广告标识 |

## 其他字段
| 字段名 | 类型 | 说明 |
|--------|------|------|
| Days | int64 | 留存天数 (如：Days=1 是新用户)|
| LastNodeId | string | 节点 ID |
| CurrencyCode | string | 货币代码 |
| LogTime | time.Time | 日志服务器收到日志的时间 |
| IpSource | string | IP 来源 |
| LogId | string | 日志 ID |
| MergeVersion | uint32 | 日志合并版本号（值为：MaxUint32-RequestTime，目的：合并时使用 RequestTime 最小值）|





- 设备列表
```
var AppleDevicesNames = []string{
	"Simulator",
	"i386",
	"iPad",
	"iPad 10",
	"iPad 2",
	"iPad 3",
	"iPad 3G",
	"iPad 4",
	"iPad 5",
	"iPad 6",
	"iPad 7",
	"iPad 8",
	"iPad 9",
	"iPad Air",
	"iPad Air 2",
	"iPad Air 3",
	"iPad Air 4",
	"iPad Air 5",
	"iPad Air 6th Gen",
	"iPad Mini",
	"iPad Mini 2",
	"iPad Mini 3",
	"iPad Mini 4",
	"iPad Mini 5",
	"iPad Mini 6",
	"iPad Pro 10.5 inch",
	"iPad Pro 11-inch",
	"iPad Pro 11-inch 2nd gen",
	"iPad Pro 11-inch 3nd gen",
	"iPad Pro 11-inch 4th gen",
	"iPad Pro 11-inch 5th Gen",
	"iPad Pro 12.9",
	"iPad Pro 12.9 inch 2nd gen",
	"iPad Pro 12.9-inch 3rd gen",
	"iPad Pro 12.9-inch 4th gen",
	"iPad Pro 12.9-inch 5th gen",
	"iPad Pro 12.9-inch 6th gen",
	"iPad Pro 12.9-inch 7th Gen",
	"iPad Pro 9.7",
	"iPhone 11",
	"iPhone 11 Pro",
	"iPhone 11 Pro Max",
	"iPhone 12",
	"iPhone 12 Pro",
	"iPhone 12 Pro Max",
	"iPhone 12 mini",
	"iPhone 13",
	"iPhone 13 Pro",
	"iPhone 13 Pro Max",
	"iPhone 13 mini",
	"iPhone 14",
	"iPhone 14 Plus",
	"iPhone 14 Pro",
	"iPhone 14 Pro Max",
	"iPhone 15",
	"iPhone 15 Plus",
	"iPhone 15 Pro",
	"iPhone 15 Pro Max",
	"iPhone 16",
	"iPhone 16 Plus",
	"iPhone 16 Pro",
	"iPhone 16 Pro Max",
	"iPhone 4",
	"iPhone 4S",
	"iPhone 5",
	"iPhone 5c",
	"iPhone 5s",
	"iPhone 6",
	"iPhone 6 Plus",
	"iPhone 6s",
	"iPhone 6s Plus",
	"iPhone 7",
	"iPhone 7 Plus",
	"iPhone 8",
	"iPhone 8 Plus",
	"iPhone SE",
	"iPhone SE 2",
	"iPhone SE 3",
	"iPhone X",
	"iPhone XR",
	"iPhone XS",
	"iPhone XS Max",
	"iPod Touch 1G",
	"iPod Touch 2G",
	"iPod Touch 3G",
	"iPod Touch 4G",
	"iPod Touch 5G",
	"iPod Touch 6G",
	"iPod Touch 7G",
	"x86_64",
}
```