package v1

import (
	"context"
	"fmt"

	"github.com/hewish1/go-selfupdate/v1/updater"
)

// Version 是当前应用程序的版本
var Version = "1.0.0"

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
func GetChangelog(ctx context.Context, u *updater.Updater, version string) (string, error) {
	if httpSource, ok := u.Source.(*updater.HTTPSource); ok {
		return httpSource.GetChangelog(ctx, version)
	}
	return "", fmt.Errorf("不支持的更新源类型")
}

// GetLatestVersionInfo 获取最新版本的详细信息
func GetLatestVersionInfo(ctx context.Context, u *updater.Updater) (*updater.LatestVersionInfo, error) {
	if httpSource, ok := u.Source.(*updater.HTTPSource); ok {
		return httpSource.GetLatestVersionInfo(ctx)
	}
	return nil, fmt.Errorf("不支持的更新源类型")
}
