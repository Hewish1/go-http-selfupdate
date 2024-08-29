Go中阿里云oss、腾讯云cos和服务器托管应用程序的自更新库
==============================================================

这个库提供了一个简单的方法来实现基于云存储的软件自动更新功能,支持阿里云OSS、腾讯云COS以及服务器等可以公开读取的存储服务。

项目正在开发中，请不要使用！！！
项目正在开发中，请不要使用！！！
项目正在开发中，请不要使用！！！


## 云存储文件结构

建议的云存储文件结构如下:
```
/
├── latest.json
├── versions/
│   ├── v1.0.0/
│   │   ├── app-windows-amd64.zip
│   │   ├── app-linux-amd64.tar.gz
│   │   └── app-darwin-amd64.tar.gz
│   ├── v1.1.0/
│   │   ├── app-windows-amd64.zip
│   │   ├── app-linux-amd64.tar.gz
│   │   └── app-darwin-amd64.tar.gz
│   └── v1.2.0/
│       ├── app-windows-amd64.zip
│       ├── app-linux-amd64.tar.gz
│       └── app-darwin-amd64.tar.gz
└── changelogs/
    ├── v1.0.0.md
    ├── v1.1.0.md
    └── v1.2.0.md
```

### 版本信息文件

`latest.json` 文件包含最新版本的信息,格式如下:
```json
{
  "version": "1.2.0",
  "releaseDate": "2023-06-15",
  "description": "这个版本修复了一些bug并添加了新功能X",
  "downloads": {
    "windows-amd64": {
      "url":"https://your-bucket.com/versions/v1.2.0/app-windows-amd64.zip",
      "md5":"1234123123123asdasdasd123123123123123"
    },
    "linux-amd64": {
      "url":"https://your-bucket.com/versions/v1.2.0/app-linux-amd64.tar.gz",
      "md5":"1234123123123asdasdasd123123123123123"
    },
    "darwin-amd64": {
      "url":"https://your-bucket.com/versions/v1.2.0/app-darwin-amd64.tar.gz",
      "md5":"1234123123123asdasdasd123123123123123"
    }
  },
  "changelogUrl": "https://your-bucket.com/changelogs/v1.2.0.md"
}
```

### 更新日志文件

更新日志文件以markdown格式编写,位于 `changelogs/` 目录下,文件名格式为 `vX.Y.Z.md`。


