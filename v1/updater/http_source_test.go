package updater

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPSource_GetLatestVersion(t *testing.T) {
	// 创建一个模拟的HTTP服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/latest.json" {
			t.Errorf("预期路径为 /latest.json, 实际为 %v", r.URL.Path)
		}
		json.NewEncoder(w).Encode(LatestVersionInfo{Version: "1.2.3"})
	}))
	defer server.Close()

	// 创建HTTPSource实例
	source := NewHTTPSource(server.URL)

	// 测试GetLatestVersion方法
	version, err := source.GetLatestVersion(context.Background())
	if err != nil {
		t.Fatalf("GetLatestVersion失败: %v", err)
	}
	if version != "1.2.3" {
		t.Errorf("预期版本为1.2.3, 实际为 %s", version)
	}
}

func TestHTTPSource_GetLatestVersionInfo(t *testing.T) {
	// 创建一个模拟的HTTP服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/latest.json" {
			t.Errorf("预期路径为 /latest.json, 实际为 %v", r.URL.Path)
		}
		json.NewEncoder(w).Encode(LatestVersionInfo{
			Version:     "1.2.3",
			ReleaseDate: "2023-07-01",
			Description: "测试版本",
			Downloads: map[string]Download{
				"windows-amd64": {URL: "http://example.com/app.zip", MD5: "abc123"},
			},
			ChangelogURL: "http://example.com/changelog.md",
		})
	}))
	defer server.Close()

	// 创建HTTPSource实例
	source := NewHTTPSource(server.URL)

	// 测试GetLatestVersionInfo方法
	info, err := source.GetLatestVersionInfo(context.Background())
	if err != nil {
		t.Fatalf("GetLatestVersionInfo失败: %v", err)
	}
	if info.Version != "1.2.3" {
		t.Errorf("预期版本为1.2.3, 实际为 %s", info.Version)
	}
	if info.ReleaseDate != "2023-07-01" {
		t.Errorf("预期发布日期为2023-07-01, 实际为 %s", info.ReleaseDate)
	}
}

func TestHTTPSource_GetUpdateInfo(t *testing.T) {
	// 创建一个模拟的HTTP服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(LatestVersionInfo{
			Version:     "1.2.3",
			ReleaseDate: "2023-07-01",
			Downloads: map[string]Download{
				"windows-amd64": {URL: "http://example.com/app.zip", MD5: "abc123"},
			},
		})
	}))
	defer server.Close()

	// 创建HTTPSource实例
	source := NewHTTPSource(server.URL)

	// 测试GetUpdateInfo方法
	info, err := source.GetUpdateInfo(context.Background(), "1.2.3")
	if err != nil {
		t.Fatalf("GetUpdateInfo失败: %v", err)
	}
	if info.Version != "1.2.3" {
		t.Errorf("预期版本为1.2.3, 实际为 %s", info.Version)
	}
}

func TestHTTPSource_DownloadFile(t *testing.T) {
	// 创建一个模拟的HTTP服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/latest.json" {
			json.NewEncoder(w).Encode(LatestVersionInfo{
				Version: "1.2.3",
				Downloads: map[string]Download{
					"windows-amd64": {URL: server.URL + "/download", MD5: "abc123"},
				},
			})
		} else if r.URL.Path == "/download" {
			w.Write([]byte("测试文件内容"))
		}
	}))
	defer server.Close()

	// 创建HTTPSource实例
	source := NewHTTPSource(server.URL)

	// 测试DownloadFile方法
	reader, err := source.DownloadFile(context.Background(), "1.2.3")
	if err != nil {
		t.Fatalf("DownloadFile失败: %v", err)
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("读取下载内容失败: %v", err)
	}
	if string(content) != "测试文件内容" {
		t.Errorf("下载内容不匹配,预期为'测试文件内容',实际为'%s'", string(content))
	}
}

func TestHTTPSource_GetChangelog(t *testing.T) {
	// 创建一个模拟的HTTP服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/latest.json" {
			json.NewEncoder(w).Encode(LatestVersionInfo{
				Version:      "1.2.3",
				ChangelogURL: server.URL + "/changelog",
			})
		} else if r.URL.Path == "/changelog" {
			w.Write([]byte("# 更新日志\n\n- 修复了一些bug\n- 添加了新功能"))
		}
	}))
	defer server.Close()

	// 创建HTTPSource实例
	source := NewHTTPSource(server.URL)

	// 测试GetChangelog方法
	changelog, err := source.GetChangelog(context.Background(), "1.2.3")
	if err != nil {
		t.Fatalf("GetChangelog失败: %v", err)
	}
	expectedChangelog := "# 更新日志\n\n- 修复了一些bug\n- 添加了新功能"
	if changelog != expectedChangelog {
		t.Errorf("更新日志内容不匹配,预期为'%s',实际为'%s'", expectedChangelog, changelog)
	}
}
