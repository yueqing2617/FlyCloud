package config

// 声明一个MySQL数据库配置
type DatabaseConfig struct {
	// 连接类型
	Type string `mapstructure:"type"`
	// 数据库地址
	Host string `mapstructure:"host"`
	// 数据库端口
	Port int `mapstructure:"port"`
	// 数据库名称
	Database string `mapstructure:"database"`
	// 数据库用户名
	Username string `mapstructure:"username"`
	// 数据库密码
	Password string `mapstructure:"password"`
	// 数据库连接池最大连接数
	MaxOpenConns int `mapstructure:"max_open_conns"`
	// 数据库连接池最大空闲连接数
	MaxIdleConns int `mapstructure:"max_idle_conns"`
	// 数据库连接池连接超时时间
	ConnMaxLifetime int `mapstructure:"conn_max_lifetime"`
	// Suffix
	Suffix  string `mapstructure:"suffix"`
	LogMode bool   `mapstructure:"log_mode"`
}
