package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
)

// HTTPSource 实现了基于 HTTP/HTTPS 的更新源
type HTTPSource struct {
	BaseURL string
	Client  *http.Client
}

// NewHTTPSource 创建一个新的 HTTPSource 实例
func NewHTTPSource(baseURL string) *HTTPSource {
	return &HTTPSource{
		BaseURL: baseURL,
		Client:  http.DefaultClient,
	}
}

// LatestVersionInfo 包含最新版本的详细信息
type LatestVersionInfo struct {
	Version      string              `json:"version"`
	ReleaseDate  string              `json:"releaseDate"`
	Description  string              `json:"description"`
	Downloads    map[string]Download `json:"downloads"`
	ChangelogURL string              `json:"changelogUrl"`
}

// Download 包含下载文件的信息
type Download struct {
	URL string `json:"url"`
	MD5 string `json:"md5"`
}

// UpdateInfo 包含更新文件的信息
type UpdateInfo struct {
	Version      string
	ReleaseDate  string
	DownloadURL  string
	Checksum     string
	Description  string
	ChangelogURL string
}

// GetLatestVersion 获取最新版本信息
func (s *HTTPSource) GetLatestVersion(ctx context.Context) (string, error) {
	info, err := s.GetLatestVersionInfo(ctx)
	if err != nil {
		return "", err
	}
	return info.Version, nil
}

// GetLatestVersionInfo 获取完整的最新版本信息
func (s *HTTPSource) GetLatestVersionInfo(ctx context.Context) (*LatestVersionInfo, error) {
	url := fmt.Sprintf("%s/latest.json", s.BaseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取最新版本失败,状态码: %d", resp.StatusCode)
	}

	var info LatestVersionInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &info, nil
}

// GetUpdateInfo 获取指定版本的更新信息
func (s *HTTPSource) GetUpdateInfo(ctx context.Context, version string) (*UpdateInfo, error) {
	info, err := s.GetLatestVersionInfo(ctx)
	if err != nil {
		return nil, err
	}

	if info.Version != version {
		return nil, fmt.Errorf("请求的版本 %s 与最新版本 %s 不匹配", version, info.Version)
	}

	platform := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
	download, ok := info.Downloads[platform]
	if !ok {
		return nil, fmt.Errorf("未找到平台 %s 的下载信息", platform)
	}

	return &UpdateInfo{
		Version:     info.Version,
		ReleaseDate: info.ReleaseDate,
		DownloadURL: download.URL,
		Checksum:    download.MD5,
	}, nil
}

// DownloadFile 下载指定版本的更新文件
func (s *HTTPSource) DownloadFile(ctx context.Context, version string) (io.ReadCloser, error) {
	info, err := s.GetUpdateInfo(ctx, version)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, info.DownloadURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("下载文件失败,状态码: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

// GetChangelog 获取指定版本的更新日志
func (s *HTTPSource) GetChangelog(ctx context.Context, version string) (string, error) {
	info, err := s.GetLatestVersionInfo(ctx)
	if err != nil {
		return "", err
	}

	if info.Version != version {
		return "", fmt.Errorf("请求的版本 %s 与最新版本 %s 不匹配", version, info.Version)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, info.ChangelogURL, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("获取更新日志失败,状态码: %d", resp.StatusCode)
	}

	changelog, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取更新日志失败: %w", err)
	}

	return string(changelog), nil
}
