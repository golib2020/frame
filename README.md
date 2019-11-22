# frame

基础开发框架

# 基础使用

``` 
go get -u github.com/golib2020/frame
```

### 配置
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

# 配置文件模版 config.json

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