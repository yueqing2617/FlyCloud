package models

import "github.com/jinzhu/gorm"

// 定义系统设置模型
type Settings struct {
	Key string `gorm:"column:key;type:varchar(255);primary_key" json:"key"`
	Val string `gorm:"column:value;type:text" json:"value"`
}

// TableName 设置表名
func (Settings) TableName() string {
	return "settings"
}

// InitSettingsTable 初始化settings
func InitSettingsTable(Db *gorm.DB) {
	// 判断表是否存在
	if !Db.HasTable(&Settings{}) {
		// 创建表
		Db.CreateTable(&Settings{})
		// 添加默认设置
		Db.Create(&Settings{
			Key: "site_name",
			Val: "FlyCloud",
		})
		Db.Create(&Settings{
			Key: "site_description",
			Val: "FlyCloud is a file storage service.",
		})
		Db.Create(&Settings{
			Key: "site_keywords",
			Val: "FlyCloud,file storage,storage service",
		})
		Db.Create(&Settings{
			Key: "site_url",
			Val: "http://flycloud.inzj.cn",
		})
		Db.Create(&Settings{
			Key: "site_email",
			Val: "empty@inzj.cn",
		})
		Db.Create(&Settings{
			Key: "site_icp",
			Val: "",
		})
		Db.Create(&Settings{
			Key: "site_copyright",
			Val: "Copyright © 2019 inzj.cn",
		})
		Db.Create(&Settings{
			Key: "site_tongji",
			Val: "",
		})
		Db.Create(&Settings{
			Key: "site_status",
			Val: "1",
		})
		Db.Create(&Settings{
			Key: "site_theme",
			Val: "default",
		})
		Db.Create(&Settings{
			Key: "site_upload_file_size",
			Val: "15728640",
		})
		Db.Create(&Settings{
			Key: "site_upload_ext",
			Val: "jpg,jpeg,png,gif,bmp,zip,rar,7z,doc,docx,xls,xlsx,ppt,pptx,pdf,txt,mp4,avi,mp3,wma,wmv,flv,swf,mkv,rm,rmvb,mov,asf,asx,vob,dat,ts,m4v,m3u8,3gp,3g2,m4a,aac,ape,ogg,wav,flac,ape,wma,mpc,mp+",
		})
		Db.Create(&Settings{
			Key: "site_upload_image_size",
			Val: "2097152",
		})
		Db.Create(&Settings{
			Key: "site_upload_image_ext",
			Val: "jpg,jpeg,png,gif,bmp",
		})
	}
}

// 过滤空值，并生成查询条件
func (settings *Settings) Filter() map[string]interface{} {
	var where map[string]interface{}
	if settings.Key != "" {
		where = make(map[string]interface{})
		where["key"] = settings.Key
	}
	return where
}

// 根据传入的key获取设置
func GetSettingsByKey(DB *gorm.DB, key string) (Settings, error) {
	var settings Settings
	err := DB.Where("key = ?", key).First(&settings).Error
	return settings, err
}

// 根据keys获取设置值
func GetSettingsByKeys(DB *gorm.DB, keys []string) (map[string]string, error) {
	var settingsList []Settings
	err := DB.Where("key in (?)", keys).Find(&settingsList).Error
	if err != nil {
		return nil, err
	}
	settings := make(map[string]string)
	for _, v := range settingsList {
		settings[v.Key] = v.Val
	}
	return settings, nil
}

// 根据key获取设置
func (settings *Settings) GetByKey(DB *gorm.DB) error {
	return DB.Where(settings.Filter()).First(settings).Error
}

// 更新设置
func (settings *Settings) Update(DB *gorm.DB) error {
	return DB.Model(settings).Where(settings.Filter()).Updates(settings).Error
}

// 获取所有设置
func (settings *Settings) GetAll(DB *gorm.DB) ([]Settings, error) {
	var settingsList []Settings
	err := DB.Find(&settingsList).Error
	return settingsList, err
}
