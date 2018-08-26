package utils

// IsLeapYear // TODO 是否是闰年
func IsLeapYear(year int) bool {
	if year%4 == 0 && year%100 != 0 || year%400 == 0 {
		return true
	}
	return false
}

// GetDayInMonthOf 获取这个月有多少天
// param int month 月份
func GetDayInMonthOf(year, month int) int {
	var day int
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		day = 31
	case 4, 6, 9, 11:
		day = 30
	default:
		if IsLeapYear(year) {
			day = 29
		} else {
			day = 28
		}
	}
	return day
}
