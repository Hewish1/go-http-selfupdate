package updater

import (
	"context"
	"io"
)

// Config 定义了更新器的配置选项
type Config struct {
	// CurrentVersion 是当前应用程序的版本
	CurrentVersion string

	// Source 是用于获取更新的来源
	Source Source

	// ValidateChecksum 指定是否验证更新文件的md5
	ValidateChecksum bool

	// Logger 用于记录更新过程中的日志
	Logger Logger
}

// Source 定义了获取更新信息和下载更新文件的接口
type Source interface {
	// GetLatestVersion 获取最新版本信息
	GetLatestVersion(ctx context.Context) (string, error)

	// GetUpdateInfo 获取指定版本的更新信息
	GetUpdateInfo(ctx context.Context, version string) (*UpdateInfo, error)

	// DownloadFile 下载指定版本的更新文件
	DownloadFile(ctx context.Context, version string) (io.ReadCloser, error)
}
