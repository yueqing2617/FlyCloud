package Db

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// 数据库时间类型
type XTime struct {
	time.Time
}

// 为XTime 重写 MarshalJSON() 方法，转换成时间戳
func (t XTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%d\"", t.Unix())
	return []byte(stamp), nil
}

// 为XTime 实现 Value() 方法，转换成时间戳
func (t XTime) Value() (driver.Value, error) {
	return t.Unix(), nil
}

// 为XTime 实现 Scan() 方法，转换成时间戳
func (t *XTime) Scan(value interface{}) error {
	var err error
	switch value.(type) {
	case time.Time:
		t.Time = value.(time.Time)
	case int64:
		t.Time = time.Unix(value.(int64), 0)
	case string:
		t.Time, err = time.Parse("2006-01-02 15:04:05", value.(string))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("can not convert %v to timestamp", value)
	}
	return nil
}
