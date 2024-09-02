Go中阿里云oss、腾讯云cos和服务器托管应用程序的自更新库
==============================================================

这个库提供了一个简单的方法来实现基于云存储的软件自动更新功能,支持阿里云OSS、腾讯云COS以及服务器等可以公开读取的存储服务。该项目是对go-update库的一个简单封装，便于直接用固定的格式进行go项目的自更新。

该项目目前仅在windows中测试了，mac和linux中如果需要使用请自行测试。

该项目参考了库https://github.com/rhysd/go-github-selfupdate。

# 使用说明

## 安装
```
go get -u github.com/iiwish/go-http-selfupdate
```

## 使用
最简单的示例
``` go
// 导入
import v1 "github.com/iiwish/go-http-selfupdate/v1"

// 创建日志文件
logger, err := v1.NewFileLogger(floder + "./update.log")
defer logger.Close()

// 创建更新器配置
config := v1.Config{
  // 获取当前版本号
  CurrentVersion: "v1.0.0",
  BaseURL:        "https://example.com/releases/", // 更新服务的根目录url
  Logger:         logger,
}

// 创建更新器
updater, err := v1.NewUpdater(config)

// 检查更新
hasUpdate, latestVersion, err := updater.CheckForUpdates()

// 获取最新版本信息,系统信息根据runtime自动获取，读取的是runtime.GOOS和runtime.GOARCH，对应latest.json中的配置
info, err := updater.GetLatestVersionInfo()

// 获取更新说明markdown文件
changeLog, err := updater.GetChangelog(latestVersion)

// 执行自更新,更新完成后软件重启才生效
err = updater.UpdateSelf()

```

## wails2中使用示例

<details>
<summary>点击展开查看</summary>

update.go文件代码如下
``` go
package update

import (
	v1 "github.com/iiwish/go-http-selfupdate/v1"
)

const (
	Version   = "v1.0.0" // 当前版本号，发版时记得修改
	UpdateURL = "https://example.com/releases/" // 更新服务的根目录url
)

type UpdaterService struct {
	updater *v1.Updater
}

func NewUpdaterService() (*UpdaterService, error) {
	// 创建日志
	logger, err := logger()
	if err != nil {
		return nil, err
	}

	// 创建更新器配置
	config := v1.Config{
		// 获取当前版本号
		CurrentVersion: Version,
		BaseURL:        UpdateURL,
		Logger:         logger,
	}

	// 创建更新器
	updater, err := v1.NewUpdater(config)
	if err != nil {
		return nil, err
	}

	return &UpdaterService{updater: updater}, nil
}

func (s *UpdaterService) CheckForUpdates() (bool, string, error) {
	// 检查更新
	hasUpdate, latestVersion, err := v1.CheckForUpdates(s.updater)
	if err != nil {
		return false, "", err
	}
	return hasUpdate, latestVersion, nil
}

func (s *UpdaterService) GetLatestVersionInfo() (*v1.LatestVersionInfo, error) {
	// 获取最新版本信息
	info, err := v1.GetLatestVersionInfo(s.updater)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (s *UpdaterService) UpdateSelf() error {
	// 执行自我更新
	err := v1.UpdateSelf(s.updater)
	if err != nil {
		return err
	}
	return nil
}

func (s *UpdaterService) Version() string {
	return Version
}

func logger() (*v1.FileLogger, error) {
	// 获取当前目录
	floder := "./"
	// 创建日志文件
	logger, err := v1.NewFileLogger(floder + "update.log")
	if err != nil {
		// 处理错误
		return nil, err
	}
	defer logger.Close()

	return logger, nil
}

```
</details>


# 服务器配置

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
  "version": "v1.2.0",
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

