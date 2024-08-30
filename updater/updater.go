package updater

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/inconshreveable/go-update"
)

// Updater 处理自动更新过程
type Updater struct {
	config Config
	source *HTTPSource
}

// NewUpdater 创建一个新的 Updater 实例
func NewUpdater(config Config) (*Updater, error) {
	if config.BaseURL == "" {
		return nil, errors.New("更新Url不能为空")
	}
	if config.CurrentVersion == "" {
		return nil, errors.New("当前版本不能为空")
	}
	if config.Logger == nil {
		config.Logger = &NoOpLogger{} // 使用一个空操作的日志记录器作为默认值
	}
	source := NewHTTPSource(config.BaseURL)
	return &Updater{config: config, source: source}, nil
}

// uncompressAndUpdate 函数用于解压并更新二进制文件
func uncompressAndUpdate(src io.Reader, assetURL, cmdPath string) error {
	// 从完整路径中提取命令名称
	_, cmd := filepath.Split(cmdPath)

	// 解压命令文件
	asset, err := UncompressCommand(src, assetURL, cmd)
	if err != nil {
		return err // 如果解压失败，返回错误
	}

	// 应用更新
	return update.Apply(asset, update.Options{
		TargetPath: cmdPath, // 指定更新后的文件保存路径
	})
}

// compareVersions 比较两个版本字符串
func compareVersions(v1, v2 string) (int, error) {
	ver1, err := version.NewVersion(v1)
	if err != nil {
		return 0, fmt.Errorf("解析版本 %s 失败: %w", v1, err)
	}

	ver2, err := version.NewVersion(v2)
	if err != nil {
		return 0, fmt.Errorf("解析版本 %s 失败: %w", v2, err)
	}

	return ver1.Compare(ver2), nil
}

// CheckForUpdates 检查是否有可用的更新
func (u *Updater) CheckForUpdates(ctx context.Context) (bool, string, error) {
	u.config.Logger.Info("开始检查更新")
	latestVersion, err := u.source.GetLatestVersion(ctx)
	if err != nil {
		u.config.Logger.Error("获取最新版本失败: %v", err)
		return false, "", fmt.Errorf("获取最新版本失败: %w", err)
	}

	u.config.Logger.Info("当前版本: %s, 最新版本: %s", u.config.CurrentVersion, latestVersion)

	comparison, err := compareVersions(u.config.CurrentVersion, latestVersion)
	if err != nil {
		u.config.Logger.Error("比较版本失败: %v", err)
		return false, "", fmt.Errorf("比较版本失败: %w", err)
	}

	hasUpdate := comparison < 0
	if hasUpdate {
		u.config.Logger.Info("发现新版本")
	} else {
		u.config.Logger.Info("当前已是最新版本")
	}
	return hasUpdate, latestVersion, nil
}

// UpdateCommand 更新指定路径的命令
func (u *Updater) UpdateTo(ctx context.Context, cmdPath string) error {
	if runtime.GOOS == "windows" && !strings.HasSuffix(cmdPath, ".exe") {
		cmdPath = cmdPath + ".exe"
	}

	stat, err := os.Lstat(cmdPath)
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %w", err)
	}
	if stat.Mode()&os.ModeSymlink != 0 {
		p, err := filepath.EvalSymlinks(cmdPath)
		if err != nil {
			return fmt.Errorf("解析符号链接失败: %w", err)
		}
		cmdPath = p
	}

	hasUpdate, latestVersion, err := u.CheckForUpdates(ctx)
	if err != nil {
		u.config.Logger.Error("检查更新失败: %v", err)
		return err
	}
	if !hasUpdate {
		u.config.Logger.Info("无需更新")
		return nil
	}

	u.config.Logger.Info("准备下载版本 %s", latestVersion)
	updateInfo, err := u.source.GetUpdateInfo(ctx, latestVersion)
	if err != nil {
		u.config.Logger.Error("获取更新信息失败: %v", err)
		return fmt.Errorf("获取更新信息失败: %w", err)
	}

	u.config.Logger.Info("开始下载更新文件")
	file, err := u.source.DownloadFile(ctx, latestVersion)
	if err != nil {
		u.config.Logger.Error("下载更新文件失败: %v", err)
		return fmt.Errorf("下载更新文件失败: %w", err)
	}
	defer file.Close()

	u.config.Logger.Info("开始替换可执行文件")
	return uncompressAndUpdate(file, updateInfo.DownloadURL, cmdPath)
}

// UpdateSelf 更新当前可执行文件
func (u *Updater) UpdateSelf(ctx context.Context) error {
	cmdPath, err := os.Executable()
	if err != nil {
		return err
	}
	return u.UpdateTo(ctx, cmdPath)
}

// GetChangelog 获取指定版本的更新日志
func (u *Updater) GetChangelog(ctx context.Context, version string) (string, error) {
	return u.source.GetChangelog(ctx, version)
}

// GetLatestVersionInfo 获取最新版本的详细信息
func (u *Updater) GetLatestVersionInfo(ctx context.Context) (*LatestVersionInfo, error) {
	return u.source.GetLatestVersionInfo(ctx)
}
