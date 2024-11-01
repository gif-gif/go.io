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
  POST https://admob.googleapis.com/v1/accounts/pub-4328354313035484/networkReport:generate

 https://admob.googleapis.com/v1/accounts/pub-4328354313035484/networkReport:generate
 
```
curl -X POST \
      https://admob.googleapis.com/v1/accounts/pub-XXXXXXXX/networkReport:generate \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer token" \
      -H "X-Goog-AuthUser": "0" \
      --data @- << EOF
{
 "report_spec": {
   "date_range": {
     "start_date": {"year": 2024, "month": 4, "day": 1},
     "end_date": {"year": 2024, "month": 4, "day": 2}
   },
   "dimensions": ["DATE"],
   "metrics": ["CLICKS", "AD_REQUESTS", "IMPRESSIONS", "ESTIMATED_EARNINGS"],
   "dimension_filters": [{"dimension": "COUNTRY", "matches_any": {"values": ["US"]}}],
   "sort_conditions": [{"metric":"CLICKS", order: "DESCENDING"}],
   "localization_settings": {"currency_code": "USD", "language_code": "en-US"}
 }
}
EOF
```

```json
[ {
  "header" : {
    "dateRange" : {
      "startDate" : {
        "year" : 2024,
        "month" : 8,
        "day" : 20
      },
      "endDate" : {
        "year" : 2024,
        "month" : 8,
        "day" : 21
      }
    },
    "localizationSettings" : {
      "currencyCode" : "USD"
    }
  }
}, {
  "row" : {
    "dimensionValues" : {
      "DATE" : {
        "value" : "20240820"
      },
      "APP" : {
        "value" : "ca-app-pub-4328354313035484~1283901700",
        "displayLabel" : "JumpJumpVPN- Fast & Secure VPN"
      },
      "COUNTRY" : {
        "value" : "AD"
      }
    },
    "metricValues" : {
      "CLICKS" : {
        "integerValue" : "0"
      },
      "ESTIMATED_EARNINGS" : {
        "microsValue" : "238"
      }
    }
  }
}, {
  "row" : {
    "dimensionValues" : {
      "DATE" : {
        "value" : "20240820"
      },
      "APP" : {
        "value" : "ca-app-pub-4328354313035484~1283901700",
        "displayLabel" : "JumpJumpVPN- Fast & Secure VPN"
      },
      "COUNTRY" : {
        "value" : "AE"
      }
    },
    "metricValues" : {
      "CLICKS" : {
        "integerValue" : "22"
      },
      "ESTIMATED_EARNINGS" : {
        "microsValue" : "1427375"
      }
    }
  }
}, {
  "row" : {
    "dimensionValues" : {
      "DATE" : {
        "value" : "20240820"
      },
      "APP" : {
        "value" : "ca-app-pub-4328354313035484~2152486667",
        "displayLabel" : "VPN - biubiuVPN Fast & Secure"
      },
      "COUNTRY" : {
        "value" : "AE"
      }
    },
    "metricValues" : {
      "CLICKS" : {
        "integerValue" : "223"
      },
      "ESTIMATED_EARNINGS" : {
        "microsValue" : "7850960"
      }
    }
  }
}, {
  "row" : {
    "dimensionValues" : {
      "DATE" : {
        "value" : "20240820"
      },
      "APP" : {
        "value" : "ca-app-pub-4328354313035484~7713423190",
        "displayLabel" : "biubiuVPN : VPN"
      },
      "COUNTRY" : {
        "value" : "AE"
      }
    },
    "metricValues" : {
      "CLICKS" : {
        "integerValue" : "71"
      },
      "ESTIMATED_EARNINGS" : {
        "microsValue" : "4050666"
      }
    }
  }
}, {
  "row" : {
    "dimensionValues" : {
      "DATE" : {
        "value" : "20240820"
      },
      "APP" : {
        "value" : "ca-app-pub-4328354313035484~8600583746",
        "displayLabel" : "JumpJumpVPN- Fast & Secure VPN"
      },
      "COUNTRY" : {
        "value" : "AE"
      }
    },
    "metricValues" : {
      "CLICKS" : {
        "integerValue" : "259"
      },
      "ESTIMATED_EARNINGS" : {
        "microsValue" : "6430629"
      }
    }
  }
}, {
  "row" : {
    "dimensionValues" : {
      "DATE" : {
        "value" : "20240820"
      },
      "APP" : {
        "value" : "ca-app-pub-4328354313035484~1283901700",
        "displayLabel" : "JumpJumpVPN- Fast & Secure VPN"
      },
      "COUNTRY" : {
        "value" : "AF"
      }
    },
    "metricValues" : {
      "CLICKS" : {
        "integerValue" : "0"
      },
      "ESTIMATED_EARNINGS" : {
        "microsValue" : "3914"
      }
    }
  }
}, {
  "row" : {
    "dimensionValues" : {
      "DATE" : {
        "value" : "20240820"
      },
      "APP" : {
        "value" : "ca-app-pub-4328354313035484~2152486667",
        "displayLabel" : "VPN - biubiuVPN Fast & Secure"
      },
      "COUNTRY" : {
        "value" : "AF"
      }
    },
    "metricValues" : {
      "CLICKS" : {
        "integerValue" : "13"
      },
      "ESTIMATED_EARNINGS" : {
        "microsValue" : "127095"
      }
    }
  }
}, {
  "row" : {
    "dimensionValues" : {
      "DATE" : {
        "value" : "20240820"
      },
      "APP" : {
        "value" : "ca-app-pub-4328354313035484~7713423190",
        "displayLabel" : "biubiuVPN : VPN"
      },
      "COUNTRY" : {
        "value" : "AF"
      }
    },
    "metricValues" : {
      "CLICKS" : {
        "integerValue" : "2"
      },
      "ESTIMATED_EARNINGS" : {
        "microsValue" : "60361"
      }
    }
  }
}, {
  "row" : {
    "dimensionValues" : {
      "DATE" : {
        "value" : "20240820"
      },
      "APP" : {
        "value" : "ca-app-pub-4328354313035484~8600583746",
        "displayLabel" : "JumpJumpVPN- Fast & Secure VPN"
      },
      "COUNTRY" : {
        "value" : "AF"
      }
    },
    "metricValues" : {
      "CLICKS" : {
        "integerValue" : "8"
      },
      "ESTIMATED_EARNINGS" : {
        "microsValue" : "67432"
      }
    }
  }
}, {
  "row" : {
    "dimensionValues" : {
      "DATE" : {
        "value" : "20240820"
      },
      "APP" : {
        "value" : "ca-app-pub-4328354313035484~2152486667",
        "displayLabel" : "VPN - biubiuVPN Fast & Secure"
      },
      "COUNTRY" : {
        "value" : "AL"
      }
    },
    "metricValues" : {
      "CLICKS" : {
        "integerValue" : "0"
      },
      "ESTIMATED_EARNINGS" : {
        "microsValue" : "341"
      }
    }
  }
}, {
  "footer" : {
    "matchingRowCount" : "846"
  }
} ]
```
