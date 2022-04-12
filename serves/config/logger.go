package config

// 声明一个日志配置
type LoggerConfig struct {
	// 日志文件路径
	Path string `mapstructure:"path"`
}
