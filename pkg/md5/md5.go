package md5

import (
	"crypto/md5"
	"fmt"
)

// MD5 密码加密
func Encry(str string) string {
	has := md5.Sum([]byte(str))
	md5 := fmt.Sprintf("%x", has)
	return md5
}

// 密码比较
func Compare(str, str2 string) bool {
	return Encry(str) == str2
}
