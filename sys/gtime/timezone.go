package gtime

import "time"

var (
	TimeZoneAsiaShanghai, _ = time.LoadLocation("Asia/Shanghai")
)

func GetLocalTimezone() (int, error) {
	return 0, nil
}

func SetLocalTimezone(timezone int) error {
	return nil
}

func ParseTimezoneCode(tz string) (offset int, err error) {
	return 0, nil //timezone.GetOffset()
}
