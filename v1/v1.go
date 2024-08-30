package v1

import (
	"context"

	"github.com/hewish1/go-selfupdate/v1/updater"
)

// Version 是当前应用程序的版本
var Version = "1.0.0"

// 直接导出 Updater 和 Config
type Updater = updater.Updater
type Config = updater.Config

// 导出其他常用类型
type (
	Logger            = updater.Logger
	HTTPSource        = updater.HTTPSource
	LatestVersionInfo = updater.LatestVersionInfo
)

// FileLogger 包装了 updater.FileLogger，并实现了 Logger 接口
type FileLogger struct {
	*updater.FileLogger
}

// NewFileLogger 创建一个新的文件日志记录器
func NewFileLogger(logPath string) (*FileLogger, error) {
	fl, err := updater.NewFileLogger(logPath)
	if err != nil {
		return nil, err
	}
	return &FileLogger{fl}, nil
}

// 导出常用函数
var (
	NewHTTPSource = updater.NewHTTPSource
)

// NewUpdater 创建一个新的 Updater 实例
func NewUpdater(config updater.Config) (*updater.Updater, error) {
	return updater.NewUpdater(config)
}

// CheckForUpdates 检查是否有可用的更新
func CheckForUpdates(ctx context.Context, u *updater.Updater) (bool, string, error) {
	return u.CheckForUpdates(ctx)
}

// UpdateSelf 执行更新操作
func UpdateSelf(ctx context.Context, u *updater.Updater) error {
	return u.UpdateSelf(ctx)
}

// GetChangelog 获取指定版本的更新日志
func GetChangelog(ctx context.Context, u *Updater, version string) (string, error) {
	return u.GetChangelog(ctx, version)
}

// GetLatestVersionInfo 获取最新版本的详细信息
func GetLatestVersionInfo(ctx context.Context, u *Updater) (*LatestVersionInfo, error) {
	return u.GetLatestVersionInfo(ctx)
}
