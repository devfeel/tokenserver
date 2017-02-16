package times

import (
	"time"
)

const (
	DefaultDateLayout     = "2006-01-02"
	DefaultFullTimeLayout = "2006-01-02 15:04:05.999999"
	DefaultTimeLayout     = "2006-01-02 15:04:05"
)

//将time转义成"2006-01-02 15:04:05" 形式的字符串
func ConvertDefaultTimeString(time time.Time) string {
	return time.Format(DefaultTimeLayout)
}

//将time转义成"2006-01-02 15:04:05.999999" 形式的字符串
func ConvertFullTimeString(time time.Time) string {
	return time.Format(DefaultFullTimeLayout)
}

//将time转义成"2006-01-02" 形式的字符串
func ConvertDateString(time time.Time) string {
	return time.Format(DefaultDateLayout)
}
