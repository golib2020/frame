# frame

[![Go Doc](https://godoc.org/github.com/golib2020/frame?status.svg)](https://godoc.org/github.com/golib2020/frame)
[![Build Status](https://travis-ci.org/golib2020/frame.svg?branch=master)](https://travis-ci.org/golib2020/frame)
[![Go Report](https://goreportcard.com/badge/github.com/golib2020/frame)](https://goreportcard.com/report/github.com/golib2020/frame)
[![Code Coverage](https://codecov.io/gh/golib2020/frame/branch/master/graph/badge.svg)](https://codecov.io/gh/golib2020/frame/branch/master)
[![Internal ready](https://img.shields.io/badge/internal-ready-success.svg)](https://github.com/golib2020/frame)
[![License](https://img.shields.io/github/license/golib2020/frame.svg?style=flat)](https://github.com/golib2020/frame)

基础开发框架

### 依赖

``` 
go get -u github.com/golib2020/frame
```

### 简单演示
```go
package main

import (
    fmt

    github.com/golib2020/frame/f
)

func main() { 
    c := f.Config()
    fmt.Println(c.GetString("db.default.dsn"))
}
```

### 配置

配置环境变量 `APP_ENV` 配置文件为 `config.<app_env>.json`

> 默认 config.json

## 模版
```json
{
  "db": {
    "default": {
      "driver": "<mysql | xorm所支持的数据库类型>", 
      "dsn": "<user>:<password>@(<host>:<port>)/<dbname>?charset=utf8mb4",
      "max": {
        "open": 10,
        "idle": 2,
        "life": "60s"
      }
    }
  },
  "redis": {
    "default": {
      "addr": "<host>:6379",
      "pass": "",
      "max": {
        "idle": 2,
        "active": 10
      }
    }
  },
  "storage": {
    "default": {
      "driver": "local",
      "root": "/storage/",
      "host": "/"
    },
    "aliyun": {
      "driver": "oss",
      "root": "< / | 根路径>",
      "host": "<访问域名的 如:https://domain/ >",
      "secret_id": "<secret_id>",
      "secret_key": "<secret_key>",
      "endpoint": "<endpoint>",
      "bucket_name": "<bucket_name>"
    },
    "tencent": {
      "driver": "cos",
      "root": "< / | 根路径>",
      "host": "<访问域名的>",
      "secret_id": "<secret_id>",
      "secret_key": "<secret_key>",
      "region": "<region>",
      "bucket_name": "<bucket_name>"
    }
  },
  "email": {
    "default": {
      "addr": "<host>:465",
      "user": "<user email>",
      "pass": "<password>",
      "name": "<发送人的名字>"
    }
  },
  "sms": {
    "default": {
      "driver": "wise",
      "api": "<api url>",
      "user": "<user>",
      "pass": "<password>"
    }
  },
  "cache": {
    "default": {
      "driver": "redis",
      "prefix": "<prefix>"
    },
    "local": {
      "driver": "local",
      "prefix": "<prefix>",
      "root": "./"
    }
  }
}
```