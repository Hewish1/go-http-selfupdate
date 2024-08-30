package v1_test

import (
	"context"
	"testing"

	v1 "github.com/iiwish/go-http-selfupdate/v1"
)

func TestSelfUpdate(t *testing.T) {
	// 使用实际的 URL
	baseURL := "https://wishte-public.oss-cn-beijing.aliyuncs.com/"
	logger, err := v1.NewFileLogger("D:/go-selfupdate-test/logfile.log")
	if err != nil {
		// 处理错误
		t.Fatalf("创建日志记录器失败: %v", err)
	}
	defer logger.Close()
	// 创建更新器配置
	config := v1.Config{
		CurrentVersion: "1.0.0",
		BaseURL:        baseURL,
		Logger:         logger,
	}
	t.Logf("更新器配置: %v", config)
	// 创建更新器
	u, err := v1.NewUpdater(config)
	if err != nil {
		t.Fatalf("创建更新器失败: %v", err)
	}

	// 检查更新
	hasUpdate, latestVersion, err := v1.CheckForUpdates(context.Background(), u)
	if err != nil {
		t.Fatalf("检查更新失败: %v", err)
	}

	t.Logf("是否有更新: %v, 最新版本: %s", hasUpdate, latestVersion)

	if !hasUpdate {
		t.Logf("没有检测到更新，当前版本可能已经是最新的")
		return // 如果没有更新，测试会在这里结束
	}

	t.Logf("检测到新版本: %s", latestVersion)

	// 获取更新日志
	changelog, err := v1.GetChangelog(context.Background(), u, latestVersion)
	if err != nil {
		t.Logf("获取更新日志失败: %v", err)
	} else {
		t.Logf("更新日志: %s", changelog)
	}

	// 获取最新版本信息
	info, err := v1.GetLatestVersionInfo(context.Background(), u)
	if err != nil {
		t.Fatalf("获取最新版本信息失败: %v", err)
	}

	t.Logf("最新版本信息: 版本 %s, 发布日期 %s", info.Version, info.ReleaseDate)

	if info.Version != latestVersion {
		t.Errorf("版本不匹配: CheckForUpdates 返回 %s, 但 GetLatestVersionInfo 返回 %s", latestVersion, info.Version)
	}

	// 更新
	err = v1.UpdateSelf(context.Background(), u)
	if err != nil {
		t.Fatalf("更新失败: %v", err)
	}
}
