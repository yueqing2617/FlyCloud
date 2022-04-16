package system

import (
	"math/rand"
	"os"
	"strconv"
	"time"
)

// 判断文件或文件夹是否存在
func IsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, err
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return false, err
}

// 创建文件夹
func MkDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// string转uint
func StrToUint(str string) uint {
	i, _ := strconv.ParseUint(str, 10, 64)
	return uint(i)
}

// 生成随机英文字符串
func RandString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

// 从time.date中获取当天日期
func GetDate() string {
	return time.Now().Format("2006-01-02")
}

// string转int64
func StrToInt64(str string) int64 {
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}

// InArray 判断字符串是否在数组中
func InArray(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

// GetRand 获取随机数
func GetRand(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

// UintToStr uint转string
func UintToStr(i uint) string {
	return strconv.FormatUint(uint64(i), 10)
}
