package updater

// Config 定义了更新器的配置选项
type Config struct {
	// CurrentVersion 是当前应用程序的版本
	CurrentVersion string

	// BaseURL 是更新源的基础URL
	BaseURL string

	// ValidateChecksum 指定是否验证更新文件的md5
	// ValidateChecksum bool

	// Logger 用于记录更新过程中的日志
	Logger Logger
}
