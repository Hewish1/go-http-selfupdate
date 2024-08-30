package updater

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"runtime"
	"strings"
)

// matchExecutableName 检查给定的命令名是否匹配目标可执行文件名
func matchExecutableName(cmd, target string) bool {
	if cmd == target {
		return true
	}

	o, a := runtime.GOOS, runtime.GOARCH

	// 当包含的可执行文件名是完整名称时（例如 foo_darwin_amd64），
	// 也被视为目标可执行文件。
	for _, d := range []rune{'_', '-'} {
		c := fmt.Sprintf("%s%c%s%c%s", cmd, d, o, d, a)
		if o == "windows" {
			c += ".exe"
		}
		if c == target {
			return true
		}
	}

	return false
}

// unarchiveTar 从tar归档中提取指定的命令文件
func unarchiveTar(src io.Reader, url, cmd string) (io.Reader, error) {
	t := tar.NewReader(src)
	for {
		h, err := t.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("解压.tar文件失败: %s", err)
		}
		_, name := filepath.Split(h.Name)
		if matchExecutableName(cmd, name) {
			log.Println("在tar归档中找到可执行文件", h.Name)
			return t, nil
		}
	}

	return nil, fmt.Errorf("在%s中未找到命令'%s'的文件", url, cmd)
}

// UncompressCommand 解压给定的源。压缩格式会根据'url'参数自动检测，
// 该参数代表资源的URL。此函数返回由'cmd'指定的解压后命令的读取器。
// 支持'.zip'、'.tar.gz'、'.tar.xz'、'.tgz'、'.gz'和'.xz'格式。
func UncompressCommand(src io.Reader, url, cmd string) (io.Reader, error) {
	if strings.HasSuffix(url, ".zip") {
		log.Println("正在解压zip文件", url)

		// Zip格式需要文件大小来解压。
		// 所以我们需要先将HTTP响应读入缓冲区。
		buf, err := io.ReadAll(src)
		if err != nil {
			return nil, fmt.Errorf("为zip文件创建缓冲区失败: %s", err)
		}
		// log.Printf("读取到的文件大小: %d 字节\n", len(buf))
		if len(buf) < 22 { // ZIP文件至少需要22字节的结束目录记录
			return nil, fmt.Errorf("ZIP文件太小，可能不完整")
		}

		r := bytes.NewReader(buf)
		// log.Println("r", r.Size())
		z, err := zip.NewReader(r, r.Size())
		if err != nil {
			return nil, fmt.Errorf("解压zip文件失败: %s", err)
		}
		for _, file := range z.File {
			_, name := filepath.Split(file.Name)
			if !file.FileInfo().IsDir() && matchExecutableName(cmd, name) {
				log.Println("在zip归档中找到可执行文件", file.Name)
				return file.Open()
			}
		}

		return nil, fmt.Errorf("在%s中未找到命令'%s'的文件", url, cmd)
	} else if strings.HasSuffix(url, ".tar.gz") || strings.HasSuffix(url, ".tgz") {
		log.Println("正在解压tar.gz文件", url)

		gz, err := gzip.NewReader(src)
		if err != nil {
			return nil, fmt.Errorf("解压.tar.gz文件失败: %s", err)
		}

		return unarchiveTar(gz, url, cmd)
	} else if strings.HasSuffix(url, ".gzip") || strings.HasSuffix(url, ".gz") {
		log.Println("正在解压gzip文件", url)

		r, err := gzip.NewReader(src)
		if err != nil {
			return nil, fmt.Errorf("解压从%s下载的gzip文件失败: %s", url, err)
		}

		name := r.Header.Name
		if !matchExecutableName(cmd, name) {
			return nil, fmt.Errorf("在%s中找到的文件名'%s'与命令'%s'不匹配", url, name, cmd)
		}
		log.Println("在gzip文件中找到可执行文件", name)
		return r, nil
	}

	log.Println("无需解压", url)
	return src, nil
}
