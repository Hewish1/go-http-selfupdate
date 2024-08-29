package updater

import (
	"context"
	"crypto/md5"
	"encoding/hex"
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
	Source Source
}

// NewUpdater 创建一个新的 Updater 实例
func NewUpdater(config Config) (*Updater, error) {
	if config.Source == nil {
		return nil, errors.New("更新源不能为空")
	}
	if config.CurrentVersion == "" {
		return nil, errors.New("当前版本不能为空")
	}
	if config.Logger == nil {
		config.Logger = &NoOpLogger{} // 使用一个空操作的日志记录器作为默认值
	}
	return &Updater{config: config}, nil
}

// uncompressAndUpdate 函数用于解压并更新二进制文件
func uncompressAndUpdate(src io.Reader, assetURL, cmdPath string) error {
	// 从完整路径中提取命令名称
	_, cmd := filepath.Split(cmdPath)

	// 解压命令文件
	// src: 压缩文件的数据源
	// assetURL: 资源的URL，可能用于日志或验证
	// cmd: 要解压的命令名称
	asset, err := UncompressCommand(src, assetURL, cmd)
	if err != nil {
		return err // 如果解压失败，返回错误
	}

	// 应用更新
	// asset: 解压后的文件数据
	// update.Options: 更新选项，指定目标路径
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
	latestVersion, err := u.config.Source.GetLatestVersion(ctx)
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
	updateInfo, err := u.config.Source.GetUpdateInfo(ctx, latestVersion)
	if err != nil {
		u.config.Logger.Error("获取更新信息失败: %v", err)
		return fmt.Errorf("获取更新信息失败: %w", err)
	}

	u.config.Logger.Info("开始下载更新文件")
	file, err := u.config.Source.DownloadFile(ctx, latestVersion)
	if err != nil {
		u.config.Logger.Error("下载更新文件失败: %v", err)
		return fmt.Errorf("下载更新文件失败: %w", err)
	}
	defer file.Close()

	if u.config.ValidateChecksum {
		u.config.Logger.Info("验证文件md5")
		if err := u.validateChecksum(file, updateInfo.Checksum); err != nil {
			return err
		}
	}

	u.config.Logger.Info("开始替换可执行文件")
	return uncompressAndUpdate(file, updateInfo.DownloadURL, cmdPath)
}

func (u *Updater) validateChecksum(file io.Reader, expectedChecksum string) error {
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return fmt.Errorf("计算MD5失败: %w", err)
	}
	computedChecksum := hex.EncodeToString(hash.Sum(nil))
	if computedChecksum != expectedChecksum {
		return fmt.Errorf("文件校验失败, 计算的MD5: %s, 预期的MD5: %s", computedChecksum, expectedChecksum)
	}
	return nil
}

func (u *Updater) UpdateSelf(ctx context.Context) error {
	cmdPath, err := os.Executable()
	if err != nil {
		return err
	}
	return u.UpdateTo(ctx, cmdPath)
}

// NoOpLogger 是一个不执行任何操作的日志记录器
type NoOpLogger struct{}

func (l *NoOpLogger) Info(format string, v ...interface{})  {}
func (l *NoOpLogger) Error(format string, v ...interface{}) {}
