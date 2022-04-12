package cache

import (
	"FlyCloud/serves/config"
	"FlyCloud/serves/logging"
	"github.com/allegro/bigcache"
	"log"
	"time"
)

// 声明一个全局的缓存对象
var Cache *bigcache.BigCache

// 初始化缓存
func InitCache(cfg *config.CacheConfig) {
	log.Println("------------------初始化缓存------------------")
	// 构建config
	config := bigcache.Config{
		Shards:             cfg.Shards,
		LifeWindow:         time.Duration(cfg.LifeWindow) * time.Minute,
		MaxEntriesInWindow: cfg.MaxEntriesWindow,
		MaxEntrySize:       cfg.MaxEntrySize,
		Verbose:            cfg.Verbose,
		HardMaxCacheSize:   cfg.HardMaxCacheSize,
	}
	// 初始化缓存
	cache, err := bigcache.NewBigCache(config)
	if err != nil {
		logging.Error("初始化缓存失败", err)
	}
	// 赋值给全局变量
	Cache = cache
	log.Println("------------------缓存初始化完成------------------")
}

// 获取缓存
func GetCache(key string) ([]byte, error) {
	return Cache.Get(key)
}

// 写入缓存
func SetCache(key string, value []byte) error {
	return Cache.Set(key, value)
}

// 删除缓存
func DeleteCache(key string) error {
	return Cache.Delete(key)
}

// 清空缓存
func ClearCache() error {
	return Cache.Reset()
}

// 获取缓存对象
func GetCacheObj() *bigcache.BigCache {
	return Cache
}
