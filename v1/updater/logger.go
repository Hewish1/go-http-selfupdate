package updater

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Logger 接口定义了日志记录的方法
type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
}

// FileLogger 实现了基于文件的日志记录
type FileLogger struct {
	logger *log.Logger
	file   *os.File
}

// NewFileLogger 创建一个新的文件日志记录器
func NewFileLogger(logPath string) (*FileLogger, error) {
	err := os.MkdirAll(filepath.Dir(logPath), 0755)
	if err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %w", err)
	}

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("打开日志文件失败: %w", err)
	}

	logger := log.New(file, "", log.LstdFlags)
	return &FileLogger{logger: logger, file: file}, nil
}

// Info 记录信息日志
func (l *FileLogger) Info(format string, v ...interface{}) {
	l.logger.Printf("[INFO] "+format, v...)
}

// Error 记录错误日志
func (l *FileLogger) Error(format string, v ...interface{}) {
	l.logger.Printf("[ERROR] "+format, v...)
}

// Close 关闭日志文件
func (l *FileLogger) Close() error {
	return l.file.Close()
}
