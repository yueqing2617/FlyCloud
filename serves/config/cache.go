package config

// 声明一个缓存配置
type CacheConfig struct {
	// 存储的条目数量，值必须是2的幂
	Shards int `mapstructure:"shards"`
	// 超时后条目将被删除
	LifeWindow int `mapstructure:"life_window"`
	// 在 LifeWindow 中的最大数量
	MaxEntriesWindow int `mapstructure:"max_entries_window"`
	// 条目的最大大小，单位是字节
	MaxEntrySize int `mapstructure:"max_entry_size"`
	// 设置缓存的最大值，以MB为单位，超过了不再分配内存，默认是0，表示不限制
	HardMaxCacheSize int `mapstructure:"hard_max_cache_size"`
	// Verbose
	Verbose bool `mapstructure:"verbose"`
}
