package utils

// GetDayInMonthOf 获取这个月有多少天
// param int month 月份
func GetDayInMonthOf(month int) int {
	var day int
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		day = 31
	case 4, 6, 9, 11:
		day = 30
	default:
		day = 28 // TODO 没计划闰年
	}
	return day
}
